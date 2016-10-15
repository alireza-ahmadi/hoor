# Hoor

> Hoor is the Persian word for sun.

Hoor is a command line tool for adding Shamsi date feature to [Hugo](http://gohugo.io)
based websites.

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

Found any issues? create an issue in the issues page.