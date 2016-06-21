# Limo

> A CLI for managing starred Git repositories

## Status

Limo is currently pre-alpha. Not a lot is working yet -- it's very much under development.

What you can currently do:

```sh
$ limo login
Enter your GitHub API token:
$ limo update
..............................................................................................................................................................................................................................................................................................................................................................................................................................................................................
Created: 0; Updated: 462; Errors: 0
$ limo list languages
AppleScript
C
C++
CSS
Clojure
CoffeeScript
Emacs Lisp
Go
HTML
Haskell
Java
JavaScript
Makefile
Objective-C
PHP
Perl
Python
Ruby
Rust
Scala
Shell
Swift
TypeScript
VimL
$ limo list stars -l viml
Shougo/dein.vim (VimL)
Shougo/neosnippet-snippets (VimL)
Yggdroot/indentLine (VimL)
chxuan/vimplus (VimL)
diepm/vim-rest-console (VimL)
flazz/vim-colorschemes (VimL)
itchyny/vim-cursorword (VimL)
jaxbot/github-issues.vim (VimL)
jaxbot/semantic-highlight.vim (VimL)
junegunn/vim-plug (VimL)
koron/minimap-vim (VimL)
mattn/emmet-vim (VimL)
mhinz/vim-galore (VimL)
morhetz/gruvbox (VimL)
neovim/neovim (VimL)
nicwest/QQ.vim (VimL)
rust-lang/rust.vim (VimL)
ryanoasis/vim-devicons (VimL)
samuelsimoes/vim-jsx-utils (VimL)
tpope/tpope (VimL)
tpope/vim-rails (VimL)
xmementoit/vim-ide (VimL)
xolox/vim-colorscheme-switcher (VimL)
```

## Installation

If you have a working Go installation, you can:

```sh
$ go get -u github.com/hoop33/limo
```

You can also download the proper binary for your system and put it in a directory in your path.

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
`show <URL/name> [--output <output>]` | Show repository details
`tag <star> <tag>...` | Tag starred repository `star` using tags `tag`
`untag <star> <tag>...` | Untag starred repository `star` using `tag` or leave blank for all tags
`update [--service <service>]` | Update `service` or leave blank for all services
`version` | Show program version

## FAQ

* Why the name "limo"?
	* If you know anything about Hollywood, you know that limos carry . . . stars.

## Contributing

Contributions are welcome! I've included a makefile that runs tests/lint/etc. to make sure code is clean, so please run `make check` before opening a pull request.

## License

Copyright &copy; 2016 Rob Warner

Licensed under the [MIT License](https://hoop33.mit-license.org/)
