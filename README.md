# protoedit

A Go library for programmatically editing Protocol Buffer (`.proto`) files while preserving formatting and comments.

## Functions

- `AppendToService` – Add new RPC definitions to existing service blocks
- `AppendToMessage` – Add fields or nested types to message blocks
- `AppendToEnum` – Add new values to enum definitions
- `Append` – Add new top-level declarations (messages, enums, services)

> **Comment preservation** – All original comments and formatting are retained

## Installation

```bash
go get github.com/whhygee/protoedit
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/whhygee/protoedit"
)

const proto = `syntax = "proto3";

service MyService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {}
}
`

func main() {
    editor, err := protoedit.New(proto)
    if err != nil {
        panic(err)
    }

    // Add a new RPC to the service
    editor.AppendToService("MyService", `  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}`)

    // Add new message definitions
    editor.Append(`message CreateUserRequest {
  string name = 1;
}

message CreateUserResponse {
  string id = 1;
}`)

    fmt.Println(editor.String())
}
```

## How It Works

`protoedit` uses [bufbuild/protocompile](https://github.com/bufbuild/protocompile) to parse proto files and extract position information. Edits are performed by inserting content at the correct byte offsets, then re-parsing to update positions for subsequent modifications.
