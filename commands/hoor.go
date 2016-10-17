// Package commands define command-line interface for Hoor, current implementation
// is based on the Cobra package.
package commands

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
	"github.com/spf13/hugo/source"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/yaa110/go-persian-calendar/ptime"
)

// HoorCmd represents the base command when called without any subcommands.
var HoorCmd = &cobra.Command{
	Use:   "hoor",
	Short: "Add Shamsi date to Hugo based websites with ease",
	PreRun: func(cmd *cobra.Command, args []string) {
		startTime = time.Now()

		if err := defineWorkingDir(); err != nil {
			return
		}

		if debugMode {
			jww.SetStdoutThreshold(jww.LevelTrace)
		}

		// set default persian date format
		viper.SetDefault("shamsiDateFormat", "dd MM yyyy")
	},
	Run: func(cmd *cobra.Command, args []string) {
		site, err := loadHugoSite()
		if err != nil {
			return
		}

		if input != "" {
			applyToFile(input)
			return
		}

		applyToSite(site)
	},
}

// Available command line flags
var (
	cfgFile    string
	sourcePath string
	contentDir string
	input      string
	debugMode  bool
	startTime  time.Time
)

// init initilizes required flags for Hoor.
func init() {
	HoorCmd.Flags().BoolVar(&debugMode, "debug", false, "runs the application in debug mode")
	HoorCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")

	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	HoorCmd.Flags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)

	HoorCmd.Flags().StringVarP(&sourcePath, "source", "s", "", "filesystem path to read files relative from")
	HoorCmd.Flags().StringVarP(&contentDir, "contentDir", "c", "", "filesystem path to content directory")
	HoorCmd.Flags().StringVarP(&input, "input", "i", "", "filesystem path of the input files")
}

// defineWorkingDir sets working directory variable, it uses source flag value in
// case of existence. Otherwise, it gets working dir from standard library.
func defineWorkingDir() error {
	var (
		dir string
		err error
	)

	if sourcePath != "" {
		if dir, err = filepath.Abs(sourcePath); err != nil {
			jww.ERROR.Println("Unable to parse source flag")
			return err
		}
	} else {
		if dir, err = os.Getwd(); err != nil {
			jww.ERROR.Println("Unable to gather working directory")
			return err
		}
	}

	viper.Set("WorkingDir", dir)
	return nil
}

// Setup adds all child commands to the root command and sets flags appropriately.
func Setup() {
	// Load configuration options into Viper
	if err := hugolib.LoadGlobalConfig(sourcePath, cfgFile); err != nil {
		jww.ERROR.Println("Cannot find configurations file")
		jww.DEBUG.Println(err)
		return
	}

	HoorCmd.AddCommand(versionCmd)

	if err := HoorCmd.Execute(); err != nil {
		jww.ERROR.Println(err)
		os.Exit(-1)
	}
}

// loadHugoSite loads hugo website configuration and initializes hugo site instance
func loadHugoSite() (site *hugolib.Site, err error) {
	var hugoSites *hugolib.HugoSites

	if hugoSites, err = hugolib.NewHugoSitesFromConfiguration(); err != nil {
		jww.ERROR.Println("Invalid configuration file")
		jww.DEBUG.Println(err)
		return
	}

	// TODO: improve the application to support multi-site installations
	site = hugoSites.Sites[0]

	if err = site.Initialise(); err != nil {
		jww.ERROR.Println("Invalid configuration file")
		jww.DEBUG.Println(err)
		return
	}

	return
}

// applyToFile generates absolute path of input file and passes it to the file processor.
func applyToFile(input string) {
	inputFile := helpers.AbsPathify(input)
	process(inputFile)
}

// applyToSite walks through Hugo site and passes each file to the file processor.
func applyToSite(site *hugolib.Site) {
	jww.FEEDBACK.Println("Started applying Shamsi date")

	sourceFiles := site.Source.Files()
	if len(sourceFiles) < 1 {
		jww.ERROR.Println("No file found in this website")
		return
	}

	results := make(chan hugolib.HandledResult)
	filechan := make(chan *source.File)
	procs := getGoMaxProcs()

	// generate concurrent sourceReader and processor
	srouceReaderWg := &sync.WaitGroup{}
	srouceReaderWg.Add(procs * 4)
	for i := 0; i < procs*4; i++ {
		go sourceReader(site, filechan, results, srouceReaderWg)
	}

	resultHandlerWg := &sync.WaitGroup{}
	resultHandlerWg.Add(procs * 4)
	for i := 0; i < procs*4; i++ {
		go resultHandler(results, resultHandlerWg)
	}

	// pass source files to filechan channel
	for _, file := range sourceFiles {
		filechan <- file
	}

	close(filechan)
	srouceReaderWg.Wait()
	close(results)
	resultHandlerWg.Wait()

	jww.FEEDBACK.Printf("Completed in %v ms\n", int(1000*time.Since(startTime).Seconds()))
}

// sourceReader receives and processes items received from filechan in concurrent.
func sourceReader(site *hugolib.Site, files <-chan *source.File, results chan<- hugolib.HandledResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		readSourceFile(site, file, results)
	}
}

// readSourceFile reads and processes each file, determines it type and publish
// the result to results channel.
func readSourceFile(site *hugolib.Site, file *source.File, results chan<- hugolib.HandledResult) {
	if handler := hugolib.NewMetaHandler(file.Extension()); handler != nil {
		handler.Read(file, site, results)
		return
	}

	jww.WARN.Println("Unsupported file type", file.Path())
}

// resultHandler waits for the results from the sourceReader and passes them to
// handle result for validation and process.
func resultHandler(results <-chan hugolib.HandledResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for r := range results {
		handleResult(r)
	}
}

// handleResult checks passed result and send it to file processor for process
// if the result is a page.
func handleResult(r hugolib.HandledResult) {
	p := r.Page()
	if p != nil {
		filePath := filepath.Join(helpers.AbsPathify(viper.GetString("contentDir")+"/"), p.FullFilePath())
		process(filePath)
	}
}

// process receives a filepath, parses it and add ShamsiDate to it.
func process(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		jww.ERROR.Println("Unable to open file", filePath)
		jww.DEBUG.Println(err)
		return
	}

	parsedPage, err := parser.ReadFrom(file)
	if err != nil {
		jww.ERROR.Println("Unable to parse file", filePath)
		jww.DEBUG.Println(err)
		return
	}

	metadata, err := parsedPage.Metadata()
	if err != nil {
		jww.ERROR.Println("Unable to extract metadata from", filePath)
		jww.DEBUG.Println(err)
		return
	}

	metadataItems := cast.ToStringMap(metadata)
	frontMatter := parsedPage.FrontMatter()

	page, err := hugolib.NewPage("pagename")
	if err != nil {
		jww.ERROR.Println("Unable to initialize page", filePath)
		jww.DEBUG.Println(err)
		return
	}

	if date, found := metadataItems["date"]; found {
		var parsedDate time.Time

		if parsedDate, err = time.Parse(time.RFC3339, date.(string)); err != nil {
			if parsedDate, err = time.Parse("2006-01-02", date.(string)); err != nil {
				jww.ERROR.Println("Invalid input date format", filePath)
				jww.DEBUG.Println(err)
				return
			}
		}

		shamsiDateFormat := viper.GetString("shamsiDateFormat")
		shamsiDate := ptime.New(parsedDate)
		metadataItems["shamsiDateF"] = persianizeNumbers(shamsiDate.Format(shamsiDateFormat))

		page.SetSourceContent(parsedPage.Content())
		err = page.SetSourceMetaData(metadataItems, rune(frontMatter[0]))
		if err != nil {
			jww.ERROR.Println("Unable to insert metadata", filePath)
			jww.DEBUG.Println(err)
			return
		}

		err = page.SaveSourceAs(filePath)
		if err != nil {
			jww.ERROR.Println("Unable to save file", filePath)
			jww.DEBUG.Println(err)
		}
	}
}

// persianizeNumbers converts English number characters to Persian number characters
func persianizeNumbers(input string) string {
	r := strings.NewReplacer("0", "۰",
		"1", "۱",
		"2", "۲",
		"3", "۳",
		"4", "۴",
		"5", "۵",
		"6", "۶",
		"7", "۷",
		"8", "۸",
		"9", "۹")

	return r.Replace(input)
}

// getGoMaxProcs returns the number of threads that the application is allowed to
// use. Function body is borrowed from the exact same function in Hugo.
func getGoMaxProcs() int {
	if gmp := os.Getenv("GOMAXPROCS"); gmp != "" {
		if p, err := strconv.Atoi(gmp); err != nil {
			return p
		}
	}

	return 1
}
