# Tog Go Client

[![CircleCI](https://img.shields.io/circleci/build/github/escaletech/tog-go/master)](https://circleci.com/gh/escaletech/workflows/tog-go)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/escaletech/tog-go?sort=semver)](https://github.com/escaletech/tog-go/releases)
[![API reference](https://img.shields.io/badge/godoc-reference-5272B4)](https://pkg.go.dev/github.com/escaletech/tog-go?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/escaletech/tog-go)](https://goreportcard.com/report/github.com/escaletech/tog-go)

Go client library that implements the Tog specification for feature flags over Redis

## Usage

### For using sessions

See [`sessions` docs](https://pkg.go.dev/github.com/escaletech/tog-go/sessions?tab=doc).

```sh
$ go get github.com/escaletech/tog-node/sessions
```

```go
package main

import (
  "context"
  "fmt"

  "github.com/escaletech/tog-node/sessions"
)

func main() {
  client, err := sessions.NewClient(flags.ClientOptions{
    Addr: "redis://localhost:6379/2",
    Cluster: false,
    OnError: func(ctx context.Context, err error) {
      // errors will always be sent here
      fmt.Println("Error:", err.Error())
    },
  })
  if err != nil {
    panic(err)
  }

  defer fc.Close()
  ctx := context.Background()

  // wherever you whish to retrieve a session
  sess := client.Session(ctx, "my_app", "the-session-id", nil)

  buttonColor := "red"
  if sess.IsSet("blue-button") {
    buttonColor = "blue"
  }

  fmt.Println("the button is", buttonColor)
}
```

### For managing flags

See [`flags` docs](https://pkg.go.dev/github.com/escaletech/tog-go/flags?tab=doc).

```sh
$ go get github.com/escaletech/tog-node/flags
```

```go
package main

import (
  "context"
  "fmt"

  "github.com/escaletech/tog-node/flags"
)

func main() {
  client, err := flags.NewClient(flags.ClientOptions{
    Addr: "redis://localhost:6379/2",
    Cluster: false,
  })
  if err != nil {
    panic(err)
  }

  defer fc.Close()
  ctx := context.Background()

  // Save a flag
  flag, err := client.SaveFlag(ctx, flags.Flag{
    Namespace: "my_app",
    Name: "blue-button",
    Description: "Makes the call-to-action button blue",
    Rollout: []flags.Rollout{{ Percentage: 30, Value: true }},
  })
  fmt.Println(flag, err)

  // List flags
  allFlags, err := client.ListFlags(ctx, "my_app")
  fmt.Println(allFlags, err)
}
```
