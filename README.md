# protoedit

A Go library for programmatically editing Protocol Buffer (`.proto`) files while preserving formatting and comments.

## API

| Method | Appends to |
|--------|------------|
| `Append.ToService(name, content)` | Service body (RPCs, options) |
| `Append.ToMessage(name, content)` | Message body (fields, nested types) |
| `Append.ToEnum(name, content)` | Enum body (values) |
| `Append.ToFile(content)` | End of file (top-level declarations) |

All methods reparse the AST after modification, enabling chained operations.

> **Comment preservation** – All original comments and formatting are retained

## Installation

```bash
go get github.com/whhygee/protoedit
```

## Usage

```go
editor, err := protoedit.New(protoContent)
if err != nil {
    return err
}

// Add a new RPC to a service
editor.Append.ToService("UserService", rpcDefinition)

// Add a field to a message
editor.Append.ToMessage("CreateUserRequest", fieldDefinition)

// Add a value to an enum
editor.Append.ToEnum("Status", enumValueDefinition)

// Add new top-level definitions
editor.Append.ToFile(messageDefinitions)

result := editor.String()
```

## How It Works

`protoedit` uses [bufbuild/protocompile](https://github.com/bufbuild/protocompile) to parse proto files and extract position information. Edits are performed by inserting content at the correct byte offsets, then re-parsing to update positions for subsequent modifications.

Chained operations work correctly:

```go
editor.Append.ToFile(newMessage)
editor.Append.ToMessage("NewMessage", field)  // AST is fresh
```
