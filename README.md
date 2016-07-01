# Limo

> A CLI for managing starred Git repositories

[![Build Status](https://travis-ci.org/hoop33/limo.svg?branch=master)](https://travis-ci.org/hoop33/limo)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](http://opensource.org/licenses/MIT)
[![Issues](https://img.shields.io/github/issues/hoop33/limo.svg)](https://github.com/hoop33/limo/issues)
[![Coverage Status](https://coveralls.io/repos/github/hoop33/limo/badge.svg?branch=master)](https://coveralls.io/github/hoop33/limo?branch=master)
[![codebeat badge](https://codebeat.co/badges/9ab79648-de9b-4585-918c-85c043bf7971)](https://codebeat.co/projects/github-com-hoop33-limo)
[![Go Report Card](https://goreportcard.com/badge/hoop33/limo)](https://goreportcard.com/report/hoop33/limo)
[![Join the chat at https://gitter.im/hoop33/limo](https://badges.gitter.im/hoop33/limo.svg)](https://gitter.im/hoop33/limo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Installation

If you have a working Go installation, you can:

```sh
$ go get -u github.com/hoop33/limo
```

Binaries for the various platforms not yet available.

## Usage

Limo is currently under development, and some things may change. Not all of it is implemented yet. Right now, GitHub is the only service supported. Here's how to get started:

### Log In to GitHub

First, create an API key for your GitHub account, and then type:

```sh
$ limo login
Enter your GitHub API token:
```

### Update Your Local Database from GitHub

```sh
$ limo update
........................................
Created: 10; Updated: 46; Errors: 0
```

### List the Languages You Have Stars In

```sh
$ limo list languages
...
Go
...
VimL
```

### List All Your Stars

```sh
$ limo list stars
...
jaxbot/github-issues.vim (VimL)
jaxbot/semantic-highlight.vim (VimL)
...
```

### List Stars for a Specific Language

```sh
$ limo list stars -l viml
...
jaxbot/github-issues.vim (VimL)
jaxbot/semantic-highlight.vim (VimL)
...
```

### Tag a Star

```sh
$ limo tag jaxbot vim github
Star 'jaxbot' ambiguous:
jaxbot/github-issues.vim (★ : 344) (VimL)
jaxbot/semantic-highlight.vim (★ : 204) (VimL)
Narrow your search
```

```sh
$ limo tag github-issues vim github
jaxbot/github-issues.vim (★ : 344) (VimL)
Added tag 'vim'
Added tag 'github'
```

### Show Details of a Star

```sh
$ limo show github-issues
jaxbot/github-issues.vim (★ : 344) (VimL)
vim, github
Github issue lookup in Vim
Home page: http://jaxbot.me/articles/github-issues-vim-plugin-5-7-2014
URL: https://github.com/jaxbot/github-issues.vim.git
Starred at Fri Feb 21 16:02:49 UTC 2014
```

### List All Your Tags

```sh
$ limo list tags
cli
git
vim
web
```

### List Stars for a Specific Tag

```sh
$ limo list stars -t vim
vim/vim (★ : 4979) (C)
tybenz/vimdeck (★ : 946) (Ruby)
jaxbot/github-issues.vim (★ : 344) (VimL)
```

### All Commands

The Limo commands take the form:

```sh
$ limo [flags] <verb> <noun> [arguments...]
```

For example, to list your starred repositories for GitHub in plaintext format, you type:

```sh
$ limo --service github --output text list stars
```

Some verbs don't require nouns, and flags can go before or after the verb-noun-arguments clause.

### Flags

Flag | Description
--- | ---
`-l, --language <language>` | The computer language
`-o, --output <color/text/csv/json>` | The output format
`-s, --service <github/gitlab/bitbucket>` | The service
`-t, --tag <name>` | The tag name
`-u, --user <user>` | The user
`-v, --verbose` | Turn on verbose output

### Verbs [Aliases]

`add`
`delete [rm]`
`help`
`list [ls]`
`login`
`open`
`prune`
`rename [mv]`
`show`
`search [find, q]`
`tag`
`untag`
`update`
`version`

### Nouns [Aliases]

`star`, `stars`
`tag`, `tags`
`trending`

### Commands

Command | Description
--- | ---
`add star <URL> [tag]...` | Star repository at `URL` and optionally add tags
`add tag <tag>` | Add tag `tag`
`delete star <URL/name>` | Unstar repository at `URL` or `name`
`delete tag <tag>` | Delete tag `tag`
`help [command]` | Show help for `command` or leave blank for all commands
`list languages` | List languages
`list stars [--tag <tag>]... [--language <language>]... [--service <service>]... [--user <user>] [--output <output>]` | List starred repositories
`list tags [--output <output>]` | List tags
`list trending [--language <language>]... [--service <service>]... [--output output]` | List trending repositories
`login [--service <service>]` | Logs in to a service (default: github)
`open <star>` | Opens the URL of a star in your default browser
`prune [--dryrun]` | Prunes unstarred items from your local database
`rename <tag> <new>` | Rename tag `tag` to `new`
`show <star> [--output <output>]` | Show repository details
`search <search string> [--output <output>]` | Search your stars
`tag <star> <tag>...` | Tag starred repository `star` using tags `tag`
`untag <star> <tag>...` | Untag starred repository `star` using `tag` or leave blank for all tags
`update [--service <service>] [--user <user>]` | Update `service` or leave blank for all services for user `user` or leave blank for you
`version` | Show program version

## Credits

Limo uses the following open source libraries -- thank you!

* [Bleve](https://github.com/blevesearch/bleve)
* [Cobra](https://github.com/spf13/cobra.git)
* [Color](https://github.com/fatih/color)
* [GORM](https://github.com/jinzhu/gorm)
* [go-github](https://github.com/google/go-github)
* [go-homedir](https://github.com/mitchellh/go-homedir)
* [Open](https://github.com/skratchdot/open-golang)
* [Testify](https://github.com/stretchr/testify)
* [YAML](https://github.com/go-yaml/yaml/tree/v2)

Apologies if I've inadvertently omitted any.

## FAQ

* Why the name "limo"?
	* If you know anything about Hollywood, you know that limos carry . . . stars.
* Where is this information stored?
	* The configuration is stored in `~/.config/limo`. Inside that directory, you'll find `limo.yaml`, which contains your GitHub API Key (so guard it!) and the path to your SQLite database, which defaults to the same directory.

## Contributing

Contributions are welcome! Please open pull requests with code that passes all the checks. See *Building* for more information.

### Building

You must have a working Go development environment to contribute code. I have tested so far only on Go 1.6.2 on OS X. `limo` uses a vendor folder, so requires Go 1.6+ or Go 1.5 with `GO15VENDOREXPERIMENT=1` (though I haven't tested that).

The included makefile performs various checks on the code. To get started, run:

```sh
$ make get-deps
```

This will install `golint` and `errcheck`. You should have to do this only once.

Then, you can run:

```sh
$ make
```

To run the code checks and tests. To build and install, run:

```sh
$ make install
```

## License

Copyright &copy; 2016 Rob Warner

Licensed under the [MIT License](https://hoop33.mit-license.org/)
