# [sauron](https://github.com/etcinit/sauron) [![GoDoc](https://godoc.org/github.com/etcinit/sauron?status.svg)](https://godoc.org/github.com/etcinit/sauron)

CLI tool for tail-ing a whole log directory recursively and dynamically

[![wercker status](https://app.wercker.com/status/ed8e8b86cb05d50c598dcff7ef070df2/m "wercker status")](https://app.wercker.com/project/bykey/ed8e8b86cb05d50c598dcff7ef070df2)

## What and why?

`tail` is really nice for monitoring logs. It also works recursively:
`tail -f **/*.log`, which means you can use it to monitor a directory of logs.
However, I couldn't find a simple command for doing this and automatically
following new files.

Sauron is a tool that allows you to do exactly that. Run it in a directory and
it will essentially output the result of calling tail -f on every single file
within the directory (even ones created after the command is ran).

Sauron is also a library, so you can extend the functionality to push logs to
an external service or database.

## Can I haz?

Yes,

```sh
# Tap
brew tap etcinit/homebrew-etcinit

# Install
brew install sauron

# Get usage info
sauron help
```

or if you have GOPATH setup:

```sh
# Get the package
go get github.com/etcinit/sauron

# Get usage info (make sure $GOPATH/bin is in your $PATH)
sauron help
```

## Future ideas:

- Filter by regex
- Integrated support for sinks (InfluxDB, ElasticSearch, MySQL, syslogd)
- Better error reporting
