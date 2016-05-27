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

`-l, --language <language>` : The computer language
`-o, --output <color|text|csv|json>` : The output format
`-s, --service <github|gitlab|bitbucket>` : The service
`-t, --tag <name>` : The tag name
`-u, --user <user>` : The user

### Verbs [Aliases]

`add`
`configure [config]`
`delete [rm]`
`list [ls]`
`login`
`open`
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

`repositories [repos]`
`repository [repo]`
`star`
`stars`
`tag`
`tags`
`trending`

### Commands

`add repository <URL> [tag]...` *Star repository at `URL` and optionally add tags*
`add tag <tag>` *Add tag `tag`*

`delete repository <URL|name>` *Unstar repository at `URL` or `name`*
`delete tag <tag>` *Delete tag `tag`*

`help [command]` *Show help for `command` or leave blank for all commands*

`list repositories [--tag <tag>]... [--language <language>]... [--service <service>]... [--user <user>] [--output <output>]` *List repositories*
`list tags [--output <output>]` *List tags*
`list trending [--language <language>]... [--service <service>]... [--output output]` *List trending repositories*

`login <service>` *Logs in to a service*

`rename <tag> <new>` *Rename tag `tag` to `new`*

`show <URL|name> [--output <output>]` *Show repository details*

`tag <repository> <tag>...` *Tag repository `repository` using tags `tag`*

`untag <repository> <tag>...` *Untag repository `repository` using `tag` or leave blank for all tags*

`update [--service <service>]...` *Update `service` or leave blank for all services*

`version` *Show program version*

## License

Copyright &copy; 2016 Rob Warner

Licensed under the [MIT License](http://hoop33.mit-license.org/)
