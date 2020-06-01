# multilisten

![](https://github.com/TheCount/go-multilisten/workflows/CI/badge.svg)
[![Documentation](https://godoc.org/github.com/TheCount/go-multilisten/multilisten?status.svg)](https://godoc.org/github.com/TheCount/go-multilisten/multilisten)

multilisten is a Go package for bundling multiple net.Listeners into a single one.

This package is useful when a third-party package expects a single `net.Listener` to build some service, but you would actually like to listen on several endpoints (ports, specific interfaces, files).

## Install

```sh
go get github.com/TheCount/go-multilisten/multilisten
```

## Usage

For the detailed API, see the [Documentation](https://godoc.org/github.com/TheCount/go-multilisten/multilisten).

Essentially, all you have to do is bundle your listeners like this:

```golang
bundle, err := multilisten.Bundle(l1, l2, l3)
```

and then use `bundle` like a single listener. It will accept from `l1`, `l2` and `l3` simultaneously.
