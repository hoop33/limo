# Limo

> A CLI for managing starred repositories

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
`configure [config]`
`delete [rm]`
`list [ls]`
`login`
`open`
`prune`
`rename [mv]`
`search`
`show`
`star`
`tag`
`unstar`
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
`prune [--dryrun]` | Prunes unstarred items from your local database
`rename <tag> <new>` | Rename tag `tag` to `new`
`show <URL/name> [--output <output>]` | Show repository details
`tag <star> <tag>...` | Tag starred repository `star` using tags `tag`
`untag <star> <tag>...` | Untag starred repository `star` using `tag` or leave blank for all tags
`update [--service <service>]` | Update `service` or leave blank for all services
`version` | Show program version

## License

Copyright &copy; 2016 Rob Warner

Licensed under the [MIT License](https://hoop33.mit-license.org/)
