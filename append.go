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

// ToService appends content at the end of a service block's body.
// Reparses the AST for chained operations.
func (a *AppendOps) ToService(serviceName, content string) error {
	for _, decl := range a.editor.file.Decls {
		if svc, ok := decl.(*ast.ServiceNode); ok {
			if svc.Name.AsIdentifier() == ast.Identifier(serviceName) {
				return a.editor.insertIntoNode(svc, content)
			}
		}
	}
	return fmt.Errorf("service %q not found", serviceName)
}

// ToMessage appends content at the end of a message block's body.
// Reparses the AST for chained operations.
func (a *AppendOps) ToMessage(messageName, content string) error {
	var target *ast.MessageNode
	ast.Walk(a.editor.file, &ast.SimpleVisitor{
		DoVisitMessageNode: func(msg *ast.MessageNode) error {
			if msg.Name.AsIdentifier() == ast.Identifier(messageName) {
				target = msg
			}
			return nil
		},
	})
	if target == nil {
		return fmt.Errorf("message %q not found", messageName)
	}
	return a.editor.insertIntoNode(target, content)
}

// ToEnum appends content at the end of an enum block's body.
// Reparses the AST for chained operations.
func (a *AppendOps) ToEnum(enumName, content string) error {
	var target *ast.EnumNode
	ast.Walk(a.editor.file, &ast.SimpleVisitor{
		DoVisitEnumNode: func(en *ast.EnumNode) error {
			if en.Name.AsIdentifier() == ast.Identifier(enumName) {
				target = en
			}
			return nil
		},
	})
	if target == nil {
		return fmt.Errorf("enum %q not found", enumName)
	}
	return a.editor.insertIntoNode(target, content)
}

// ToFile appends content at the end of the proto file.
// Reparses the AST for chained operations.
func (a *AppendOps) ToFile(content string) error {
	if !strings.HasSuffix(a.editor.content, "\n") {
		a.editor.content += "\n"
	}
	a.editor.content += "\n" + strings.TrimSpace(content) + "\n"
	return a.editor.reparse()
}
