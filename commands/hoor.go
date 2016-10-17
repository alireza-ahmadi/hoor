// Package commands define command-line interface for Hoor, current implementation
// is based on the Cobra package.
package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/source"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// HoorCmd represents the base command when called without any subcommands.
var HoorCmd = &cobra.Command{
	Use:   "hoor",
	Short: "Add Shamsi date to Hugo based websites with ease",
	Run: func(cmd *cobra.Command, args []string) {
		site, err := loadHugoSite()
		if err != nil {
			return
		}

		if input != "" {
			// do something
		}

		parseAndApply(site)
	},
}

// Available command line flags
var (
	cfgFile    string
	sourcePath string
	contentDir string
	input      string
)

// init initilizes required flags for Hoor.
func init() {
	// Persistent flags
	HoorCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")

	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	HoorCmd.Flags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)

	// Common flags
	HoorCmd.Flags().StringVarP(&sourcePath, "source", "s", "", "filesystem path to read files relative from")
	HoorCmd.Flags().StringVarP(&contentDir, "contentDir", "c", "", "filesystem path to content directory")
	HoorCmd.Flags().StringVarP(&input, "input", "i", "", "filesystem path of the input files")
}

// Setup adds all child commands to the root command and sets flags appropriately.
func Setup() {
	// Load configuration options into Viper
	if err := hugolib.LoadGlobalConfig(sourcePath, cfgFile); err != nil {
		jww.ERROR.Println("Cannot find configurations file")
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
		return
	}

	// TODO: improve the application to support multi-site installations
	site = hugoSites.Sites[0]

	if err = site.Initialise(); err != nil {
		jww.ERROR.Println("Invalid configuration file")
		return
	}

	return
}

func parseAndApply(site *hugolib.Site) {
	errs := make(chan error)

	files := site.Source.Files()

	if len(files) < 1 {
		close(errs)
		return
	}

	results := make(chan hugolib.HandledResult)
	filechan := make(chan *source.File)

	procs := getGoMaxProcs()
	wg := &sync.WaitGroup{}

	wg.Add(procs * 4)
	for i := 0; i < procs*4; i++ {
		fmt.Println(i)
		go sourceReader(site, filechan, results, wg)
	}

	go process(results)

	for _, file := range files {
		filechan <- file
	}

	close(filechan)
	wg.Wait()
	close(results)
}

func sourceReader(s *hugolib.Site, files <-chan *source.File, results chan<- hugolib.HandledResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		readSourceFile(s, file, results)
	}
}

func readSourceFile(s *hugolib.Site, file *source.File, results chan<- hugolib.HandledResult) {
	h := hugolib.NewMetaHandler(file.Extension())
	if h != nil {
		h.Read(file, s, results)
	} else {
		log.Println("Unsupported File Type", file.Path())
	}
}

func process(results <-chan hugolib.HandledResult) {
	for r := range results {
		p := r.Page()
		if p != nil {
			fmt.Println(filepath.Join(helpers.AbsPathify(viper.GetString("contentDir")+"/"), p.FullFilePath()))
		}
	}
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
