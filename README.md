# Limo

> A CLI for managing starred Git repositories

[![Build Status](https://travis-ci.org/hoop33/limo.svg?branch=master)](https://travis-ci.org/hoop33/limo)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](http://opensource.org/licenses/MIT)
[![Issues](https://img.shields.io/github/issues/hoop33/limo.svg)](https://github.com/hoop33/limo/issues)
[![Coverage Status](https://coveralls.io/repos/github/hoop33/limo/badge.svg?branch=master)](https://coveralls.io/github/hoop33/limo?branch=master)
[![codebeat badge](https://codebeat.co/badges/9ab79648-de9b-4585-918c-85c043bf7971)](https://codebeat.co/projects/github-com-hoop33-limo)
[![Go Report Card](https://goreportcard.com/badge/hoop33/limo)](https://goreportcard.com/report/hoop33/limo)
[![Join the chat at https://gitter.im/hoop33/limo](https://badges.gitter.im/hoop33/limo.svg)](https://gitter.im/hoop33/limo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Table of Contents

* [Introduction](#introduction)
* [Installation](#installation)
* [Usage](#usage)
* [FAQ](#faq)
* [Contributing](#contributing)
* [Credits](#credits)
* [License](#license)

## Introduction

Both [GitHub](https://github.com) and [GitLab](https://gitlab.com) allow you to "star" repositories, and [Bitbucket](https://bitbucket.org) lets you "watch" them. "Starring" or "watching" lets you keep track of repositories you find interesting, but none of the services provide ways to search or tag your repositories so you can easily find them.

Limo lets you manage your starred repositories from the command line. You can do things like tag them, search them, or list them by language. Think of Limo as the CLI version of [Astral](https://app.astralapp.com/) (also worth looking into).

## Installation

If you have a working Go installation, type:

```sh
$ go get -u github.com/hoop33/limo
```

You can also download a pre-built binary for your platform and copy it to a location on your path.

* [macOS](https://github.com/hoop33/limo/raw/master/dist/macos/limo)
* [Linux](https://github.com/hoop33/limo/raw/master/dist/linux/limo)
* [Windows](https://github.com/hoop33/limo/raw/master/dist/windows/limo.exe)

## Usage

Limo supports both GitHub and GitLab, and Bitbucket support is coming.

You can read the full usage documentation at <https://www.gitbook.com/book/hoop33/limo/details>.

Here's how to get started:

### Log In to GitHub and GitLab

First, create API keys for your GitHub and GitLab accounts on their respective sites, and then type:

```sh
$ limo login
Enter your GitHub API token:

$ limo login --service gitlab
Enter your GitLab API token:
```

### Update Your Local Database

```sh
$ limo update
Updating . . . /
Created: 10; Updated: 46; Errors: 0

$ limo update --service gitlab
Updating . . . /
Created: 5; Updated: 23; Errors: 0
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
jaxbot/github-issues.vim ★ :347 VimL https://github.com/jaxbot/github-issues.vim.git
jaxbot/semantic-highlight.vim ★ :209 VimL https://github.com/jaxbot/semantic-highlight.vim.git
...
```

### List Stars for a Specific Language

```sh
$ limo list stars -l viml
...
jaxbot/github-issues.vim ★ :347 VimL https://github.com/jaxbot/github-issues.vim.git
jaxbot/semantic-highlight.vim ★ :209 VimL https://github.com/jaxbot/semantic-highlight.vim.git
...
```

### Tag a Star

```sh
$ limo tag jaxbot vim github
Star 'jaxbot' ambiguous:
jaxbot/github-issues.vim ★ :347 VimL https://github.com/jaxbot/github-issues.vim.git
jaxbot/semantic-highlight.vim ★ :209 VimL https://github.com/jaxbot/semantic-highlight.vim.git
Narrow your search
```

```sh
$ limo tag github-issues vim github
jaxbot/github-issues.vim ★ :347 VimL https://github.com/jaxbot/github-issues.vim.git
Added tag 'vim'
Added tag 'github'
```

### Show Details of a Star

```sh
$ limo show github-issues
jaxbot/github-issues.vim ★ :347 VimL https://github.com/jaxbot/github-issues.vim.git
github, vim
Github issue lookup in Vim
Home page: http://jaxbot.me/articles/github-issues-vim-plugin-5-7-2014
Starred on Fri Feb 21 16:02:49 UTC 2014
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
vim/vim ★ :5267 C https://github.com/vim/vim.git
tybenz/vimdeck ★ :953 Ruby https://github.com/tybenz/vimdeck.git
jaxbot/github-issues.vim ★ :347 VimL https://github.com/jaxbot/github-issues.vim.git
vicoapp/vico ★ :666 Objective-C https://github.com/vicoapp/vico.git
```

### Perform a Full-text Search on Your Stars

```sh
$ limo search text editor
(0.708069) limetext/lime ★ :12588 https://github.com/limetext/lime.git
(0.617632) zyedidia/micro ★ :2030 Go https://github.com/zyedidia/micro.git
(0.614035) driusan/de ★ :116 Go https://github.com/driusan/de.git
(0.611212) martanne/vis ★ :2171 C https://github.com/martanne/vis.git
(0.608427) Cocoanetics/DTRichTextEditor ★ :259 Objective-C https://github.com/Cocoanetics/DTRichTextEditor.git
(0.608427) atom/atom ★ :29489 CoffeeScript https://github.com/atom/atom.git
(0.608427) skx/kilua ★ :80 C++ https://github.com/skx/kilua.git
(0.600297) antirez/kilo ★ :2672 C https://github.com/antirez/kilo.git
(0.597658) xmementoit/vim-ide ★ :150 VimL https://github.com/xmementoit/vim-ide.git
(0.592483) edwardloveall/atom-replacement-icon ★ :133 Shell https://github.com/edwardloveall/atom-replacement-icon.git
```

You can read the full usage documentation at <https://www.gitbook.com/book/hoop33/limo/details>.

## FAQ

* Why the name "limo"?
	* If you know anything about Hollywood, you know that limos carry . . . stars.
* Where is this information stored?
  * The configuration is stored in `~/.config/limo`. Inside that directory, you'll find:
    * `limo.yaml`: Configuration information
    * `limo.db`: The SQLite database that stores all your stars and tags
    * `limo.idx`: The Bleve search index
* How do I change the "updating" spinner?
  * Limo uses <https://github.com/briandowns/spinner> for its "updating" spinner. You can override which spinner is used, what color to make it, and the spin interval in your `limo.yaml` file, like this:
    ```yaml
    outputs:
      color:
        spinnerIndex: 10
        spinnerInterval: 300
        spinnerColor: blue
    ```

## Contributing

Please note that this project is released with a [Contributor Code of Conduct](http://contributor-covenant.org/). By participating in this project you agree to abide by its terms. See [CODE_OF_CONDUCT](CODE_OF_CONDUCT.md) file.

Contributions are welcome! Please open pull requests with code that passes all the checks. See *Building* for more information.

### Building

You must have a working Go development environment to contribute code. `limo` uses vendoring, so requires Go 1.6+ (or Go 1.5 with `GO15VENDOREXPERIMENT=1`, though I haven't tested that).

The included makefile performs various checks on the code. To get started, run:

```sh
$ make deps
```

This will install the necessary dependencies. You should have to do this only once.

Then, you can run:

```sh
$ make
```

To run the code checks and tests. To build and install, run:

```sh
$ make install
```

## Credits

Limo uses the following open source libraries -- thank you!

* [Bleve](https://github.com/blevesearch/bleve)
* [Cobra](https://github.com/spf13/cobra.git)
* [Color](https://github.com/fatih/color)
* [GORM](https://github.com/jinzhu/gorm)
* [go-github](https://github.com/google/go-github)
* [go-gitlab](https://github.com/xanzy/go-gitlab)
* [go-humanize](https://github.com/dustin/go-humanize)
* [spinner](https://github.com/briandowns/spinner)
* [Open](https://github.com/skratchdot/open-golang)
* [Testify](https://github.com/stretchr/testify)
* [xdgbasedir](https://github.com/cep21/xdgbasedir)
* [YAML](https://github.com/go-yaml/yaml/tree/v2)

Also, thanks to contributors!

* [Tyson Warner](https://github.com/nylon22) : List stars by tag and language

Apologies if I've inadvertently omitted any library or any contributor.

## License

Copyright &copy; 2016 Rob Warner

Licensed under the [MIT License](https://hoop33.mit-license.org/)
