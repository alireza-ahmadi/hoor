# Hoor [![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/alireza-ahmadi/hoor)

> Hoor is the Persian word for sun.

Hoor is a command line tool for adding Shamsi date feature to [Hugo](http://gohugo.io)
based websites.

## Installation

You can download the latest version of **Hoor** from [Releases](https://github.com/alireza-ahmadi/hoor/releases)
page. Also, if you already installed Golang toolchain, you can install hoor by running
the following command:

    go get github/alireza-ahamdi/hoor

## Usage

### IMPORTANT NOTE

**Current version is the first implementation of this utility, if you want to use
it, please make sure that you have a backup of your website first.**

Adding Shamsi date to all contents(or changing date format):

```sh
cd path/to/your/website
hoor
```

Adding Shamsi date to a single content:

```sh
cd path/to/your/website
hoor -i content/post/foobar.md
``` 

After adding Shamsi date to your contents, you can display Shamsi date in your posts
by using `{{ .Params.shamsidate }}` in your templates. The default format of shamsi
date is `d M yyyy` but you can change it by adding `shamsiDateFormat` configuration
option to your website config file.

```toml
baseURL = "http://alireza.es"
title = "بلاگ علیرضا احمدی"
theme = "vivid"
...
shamsiDateFormat = "d MMMماه yyyy"
...
[indexes]
   topic = "topics"
```

These are the available formmating options:

```
yyyy, yyy, y     year (e.g. 1394)
yy               2-digits representation of year (e.g. 94)
MMM              the Persian name of month (e.g. فروردین)
MMI              the Dari name of month (e.g. حمل)
MM               2-digits representation of month (e.g. 01)
M                month (e.g. 1)
rw               remaining weeks of year
w                week of year
W                week of month
RD               remaining days of year
D                day of year
rd               remaining days of month
dd               2-digits representation of day (e.g. 01)
d                day (e.g. 1)
E                the Persian name of weekday (e.g. شنبه)
e                the Persian short name of weekday (e.g. ش)
A                the Persian name of 12-Hour marker (e.g. قبل از ظهر)
a                the Persian short name of 12-Hour marker (e.g. ق.ظ)
HH               2-digits representation of hour [00-23]
H                hour [0-23]
kk               2-digits representation of hour [01-24]
k                hour [1-24]
hh               2-digits representation of hour [01-12]
h                hour [1-12]
KK               2-digits representation of hour [00-11]
K                hour [0-11]
mm               2-digits representation of minute [00-59]
m                minute [0-59]
ss               2-digits representation of seconds [00-59]
s                seconds [0-59]
ns               nanoseconds
S                3-digits representation of milliseconds (e.g. 001)
z                the name of location
Z                zone offset (e.g. +03:30)
```

Learn more about formatting options on [ptime](https://github.com/yaa110/go-persian-calendar)
repository.

Also, if you want to get help about command line options, run the following command:

    hoor -h

## How it works?

Due to the lack of plugin system in Hugo, there is no way to change date format
in the build process. I approached a simple solution for fixing this problem. By
running `hoor` command in your hugo site source code, **Hoor** will automatically
finds all of your posts/pages, read the gregorian date and convert it to a Shamsi(Jalali) 
date string then add the result to the [front matter](https://gohugo.io/content/front-matter/)
list. Afterwards, you can use that string as a [template variable](https://gohugo.io/templates/variables/).

#### So why shouldn't I generate that string manually, using a simple code editor?

Of course you can write it manually, but applying any change in date format would
be a great pain in future, by using **Hoor** you can change date format of 50+ files
in a few seconds while changing them manually would take so much longer. It's all
about maintainability.

## Issues

Found any issues? [Create an issue](https://github.com/alireza-ahmadi/hoor/issues)
in the issues page.

## TODO

- Add more tests

## Credits

**Hoor** is built upon [Hugo](https://github.com/spf13/hugo) and [ptime](https://github.com/yaa110/go-persian-calendar).