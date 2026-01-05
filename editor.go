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

	Append *AppendOps
}

// New creates a new Editor from proto file content.
// Returns an error if the content cannot be parsed.
func New(content string) (*Editor, error) {
	handler := reporter.NewHandler(nil)
	file, err := parser.Parse("input.proto", strings.NewReader(content), handler)
	if err != nil {
		return nil, fmt.Errorf("parse proto: %w", err)
	}
	e := &Editor{
		content: content,
		file:    file,
	}
	e.Append = &AppendOps{editor: e}
	return e, nil
}

// String returns the current proto content with all modifications applied.
func (e *Editor) String() string {
	return e.content
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

// insertBeforePos inserts content before the given position and reparses.
func (e *Editor) insertBeforePos(pos int, content string) error {
	e.content = e.content[:pos] + "\n" + content + "\n" + e.content[pos:]
	return e.reparse()
}
