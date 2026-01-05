// Package protoedit provides a simple API for programmatically editing .proto files.
// It wraps github.com/bufbuild/protocompile to parse proto files and uses position
// information to perform edits while preserving formatting and comments.
package protoedit

import (
	"fmt"
	"strings"

	"github.com/bufbuild/protocompile/ast"
	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
)

// Editor provides methods to edit proto file content.
type Editor struct {
	content string
	file    *ast.FileNode
}

// New creates a new Editor from proto file content.
// Returns an error if the content cannot be parsed.
func New(content string) (*Editor, error) {
	handler := reporter.NewHandler(nil)
	file, err := parser.Parse("input.proto", strings.NewReader(content), handler)
	if err != nil {
		return nil, fmt.Errorf("parse proto: %w", err)
	}
	return &Editor{
		content: content,
		file:    file,
	}, nil
}

// AppendToService appends content at the end of a service block's body
// (right before the closing brace). The content should be valid proto service elements
// (e.g., RPC definitions, options).
func (e *Editor) AppendToService(serviceName, content string) error {
	for _, decl := range e.file.Decls {
		if svc, ok := decl.(*ast.ServiceNode); ok {
			if svc.Name.AsIdentifier() == ast.Identifier(serviceName) {
				closeInfo := e.file.NodeInfo(svc.CloseBrace)
				insertPos := closeInfo.Start().Offset
				return e.insertBeforePos(insertPos, content)
			}
		}
	}
	return fmt.Errorf("service %q not found", serviceName)
}

// AppendToMessage appends content at the end of a message block's body
// (right before the closing brace). The content should be valid proto message elements
// (e.g., field definitions, nested messages, options).
func (e *Editor) AppendToMessage(messageName, content string) error {
	var found bool
	ast.Walk(e.file, &ast.SimpleVisitor{
		DoVisitMessageNode: func(msg *ast.MessageNode) error {
			if msg.Name.AsIdentifier() == ast.Identifier(messageName) {
				closeInfo := e.file.NodeInfo(msg.CloseBrace)
				insertPos := closeInfo.Start().Offset
				if err := e.insertBeforePos(insertPos, content); err != nil {
					return err
				}
				found = true
			}
			return nil
		},
	})
	if !found {
		return fmt.Errorf("message %q not found", messageName)
	}
	return nil
}

// AppendToEnum appends content at the end of an enum block's body
// (right before the closing brace). The content should be valid proto enum elements
// (e.g., enum values, options).
func (e *Editor) AppendToEnum(enumName, content string) error {
	var found bool
	ast.Walk(e.file, &ast.SimpleVisitor{
		DoVisitEnumNode: func(en *ast.EnumNode) error {
			if en.Name.AsIdentifier() == ast.Identifier(enumName) {
				closeInfo := e.file.NodeInfo(en.CloseBrace)
				insertPos := closeInfo.Start().Offset
				if err := e.insertBeforePos(insertPos, content); err != nil {
					return err
				}
				found = true
			}
			return nil
		},
	})
	if !found {
		return fmt.Errorf("enum %q not found", enumName)
	}
	return nil
}

// Append appends content at the end of the proto file.
// Useful for adding new top-level definitions (messages, enums, services).
func (e *Editor) Append(content string) {
	if !strings.HasSuffix(e.content, "\n") {
		e.content += "\n"
	}
	e.content += "\n" + content
}

// String returns the current proto content with all modifications applied.
func (e *Editor) String() string {
	return e.content
}

// insertBeforePos inserts content before the given position and reparses.
func (e *Editor) insertBeforePos(pos int, content string) error {
	e.content = e.content[:pos] + "\n" + content + "\n" + e.content[pos:]
	return e.reparse()
}

// reparse re-parses the content after modifications.
// This is needed to update position information for subsequent edits.
func (e *Editor) reparse() error {
	handler := reporter.NewHandler(nil)
	file, err := parser.Parse("input.proto", strings.NewReader(e.content), handler)
	if err != nil {
		return fmt.Errorf("reparse proto: %w", err)
	}
	e.file = file
	return nil
}
