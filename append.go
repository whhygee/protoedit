package protoedit

import (
	"fmt"
	"strings"

	"github.com/bufbuild/protocompile/ast"
)

// AppendOps provides methods for appending content to proto elements.
type AppendOps struct {
	editor *Editor
}

// ToService appends content at the end of a service block's body
// (right before the closing brace). The content should be valid proto service elements
// (e.g., RPC definitions, options). Reparses the AST for chained operations.
func (a *AppendOps) ToService(serviceName, content string) error {
	for _, decl := range a.editor.file.Decls {
		if svc, ok := decl.(*ast.ServiceNode); ok {
			if svc.Name.AsIdentifier() == ast.Identifier(serviceName) {
				closeInfo := a.editor.file.NodeInfo(svc.CloseBrace)
				insertPos := closeInfo.Start().Offset
				return a.editor.insertBeforePos(insertPos, content)
			}
		}
	}
	return fmt.Errorf("service %q not found", serviceName)
}

// ToMessage appends content at the end of a message block's body
// (right before the closing brace). The content should be valid proto message elements
// (e.g., field definitions, nested messages, options). Reparses the AST for chained operations.
func (a *AppendOps) ToMessage(messageName, content string) error {
	var found bool
	ast.Walk(a.editor.file, &ast.SimpleVisitor{
		DoVisitMessageNode: func(msg *ast.MessageNode) error {
			if msg.Name.AsIdentifier() == ast.Identifier(messageName) {
				closeInfo := a.editor.file.NodeInfo(msg.CloseBrace)
				insertPos := closeInfo.Start().Offset
				if err := a.editor.insertBeforePos(insertPos, content); err != nil {
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

// ToEnum appends content at the end of an enum block's body
// (right before the closing brace). The content should be valid proto enum elements
// (e.g., enum values, options). Reparses the AST for chained operations.
func (a *AppendOps) ToEnum(enumName, content string) error {
	var found bool
	ast.Walk(a.editor.file, &ast.SimpleVisitor{
		DoVisitEnumNode: func(en *ast.EnumNode) error {
			if en.Name.AsIdentifier() == ast.Identifier(enumName) {
				closeInfo := a.editor.file.NodeInfo(en.CloseBrace)
				insertPos := closeInfo.Start().Offset
				if err := a.editor.insertBeforePos(insertPos, content); err != nil {
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

// ToFile appends content at the end of the proto file.
// Useful for adding new top-level definitions (messages, enums, services).
// Reparses the AST for chained operations.
func (a *AppendOps) ToFile(content string) error {
	if !strings.HasSuffix(a.editor.content, "\n") {
		a.editor.content += "\n"
	}
	a.editor.content += "\n" + content
	return a.editor.reparse()
}
