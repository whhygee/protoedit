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
func (e *Editor) reparse() error {
	handler := reporter.NewHandler(nil)
	file, err := parser.Parse("input.proto", strings.NewReader(e.content), handler)
	if err != nil {
		return fmt.Errorf("reparse proto: %w", err)
	}
	e.file = file
	return nil
}

// nodeDepth calculates the nesting depth of a node by counting ancestors.
func (e *Editor) nodeDepth(target ast.Node) int {
	depth := 0
	found := false

	var walk func(node ast.Node, d int)
	walk = func(node ast.Node, d int) {
		if found {
			return
		}
		if node == target {
			depth = d
			found = true
			return
		}
		switch n := node.(type) {
		case *ast.FileNode:
			for _, decl := range n.Decls {
				walk(decl, d)
			}
		case *ast.MessageNode:
			for _, decl := range n.Decls {
				walk(decl, d+1)
			}
		case *ast.ServiceNode:
			for _, decl := range n.Decls {
				walk(decl, d+1)
			}
		case *ast.EnumNode:
			for _, decl := range n.Decls {
				walk(decl, d+1)
			}
		}
	}

	walk(e.file, 0)
	return depth
}

// insertIntoNode inserts content into a node (before its closing brace).
func (e *Editor) insertIntoNode(node ast.Node, content string) error {
	var closeBrace ast.Node
	switch n := node.(type) {
	case *ast.MessageNode:
		closeBrace = n.CloseBrace
	case *ast.ServiceNode:
		closeBrace = n.CloseBrace
	case *ast.EnumNode:
		closeBrace = n.CloseBrace
	default:
		return fmt.Errorf("unsupported node type")
	}

	depth := e.nodeDepth(node) + 1 // Content is one level inside the node
	pos := e.file.NodeInfo(closeBrace).Start().Offset
	unit := detectIndent(e.content)
	content = applyIndent(content, depth, unit)

	// Check if inserting into empty inline block `{}`
	if pos > 0 && e.content[pos-1] == '{' {
		e.content = e.content[:pos] + "\n" + content + "\n" + strings.Repeat(unit, depth-1) + e.content[pos:]
	} else {
		e.content = e.content[:pos] + "\n" + content + "\n" + e.content[pos:]
	}
	return e.reparse()
}

// detectIndent returns the indent unit used in content (tab or spaces).
func detectIndent(content string) string {
	for _, line := range strings.Split(content, "\n") {
		if len(line) > 0 && (line[0] == '\t' || line[0] == ' ') {
			if line[0] == '\t' {
				return "\t"
			}
			return line[:len(line)-len(strings.TrimLeft(line, " "))]
		}
	}
	return "  "
}

// applyIndent ensures each line has the given indentation.
func applyIndent(content string, depth int, unit string) string {
	indent := strings.Repeat(unit, depth)
	lines := strings.Split(strings.TrimSpace(content), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = indent + strings.TrimLeft(line, " \t")
		}
	}
	return strings.Join(lines, "\n")
}
