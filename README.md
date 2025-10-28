# multilisten

![](https://github.com/TheCount/go-multilisten/workflows/CI/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/TheCount/go-multilisten.svg)](https://pkg.go.dev/github.com/TheCount/go-multilisten/pkg/multilisten)

multilisten is a Go package for bundling multiple net.Listeners into a single one.

This package is useful when a third-party package expects a single `net.Listener` to build some service, but you would actually like to listen on several endpoints (ports, specific interfaces, files).

## Install

In your `go.mod` directory, run:
```sh
go get github.com/TheCount/go-multilisten/pkg/multilisten
```
See also https://go.dev/doc/modules/managing-dependencies#workflow.

## Usage

For the detailed API, see the [Documentation](https://pkg.go.dev/github.com/TheCount/go-multilisten/pkg/multilisten).

Essentially, all you have to do is bundle your listeners like this:

```golang
bundle, err := multilisten.Bundle(l1, l2, l3)
```

and then use `bundle` like a single listener. It will accept from `l1`, `l2` and `l3` simultaneously.

For a more flexible API, you can also use:
```golang
set := multilisten.NewSet(mainAddr)
if err := set.Add(l1); err != nil {
  return err
}
if err := set.Add(l2); err != nil {
  return err
}
if err := set.InjectConn(conn); err != nil {
  return err
}
…
```
