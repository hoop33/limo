# Limo

> A CLI for managing starred Git repositories

## Status

Limo is currently under development, and some things may change. Not all of it is implemented yet. Here's some of what you can do:

### Log In to GitHub

First, create an API key for your GitHub account, and then type:

```sh
$ limo login
Enter your GitHub API token:
```

### Update Your Local Database from GitHub

```sh
$ limo update
..............................................................................................................................................................................................................................................................................................................................................................................................................................................................................
Created: 0; Updated: 462; Errors: 0
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
Awesome
cli
git
github
go
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

## Installation

If you have a working Go installation, you can:

```sh
$ go get -u github.com/hoop33/limo
```

Binaries for the various platforms not yet available.

## Usage

The Limo commands take the form:

```sh
$ limo [flags] <verb> <noun> [arguments...]
```

For example, to list your starred repositories for Github in JSON format, you type:

```sh
$ limo --service github --output json list repositories
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
`tag <star> <tag>...` | Tag starred repository `star` using tags `tag`
`untag <star> <tag>...` | Untag starred repository `star` using `tag` or leave blank for all tags
`update [--service <service>]` | Update `service` or leave blank for all services
`version` | Show program version

## FAQ

* Why the name "limo"?
	* If you know anything about Hollywood, you know that limos carry . . . stars.
* Where is this information stored?
	* The configuration is stored in `~/.config/limo`. Inside that directory, you'll find `limo.yaml`, which contains your GitHub API Key (so guard it!) and the path to your SQLite database, which defaults to the same directory.

## Contributing

Contributions are welcome! I've included a makefile that runs tests/lint/etc. to make sure code is clean, so please run `make check` before opening a pull request.

## License

Copyright &copy; 2016 Rob Warner

Licensed under the [MIT License](https://hoop33.mit-license.org/)
