// Package highlight provides syntax highlighting
// for Go, HTML, JS, Bash source codes using tree-sitter.
package highlight

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/sql"
)

func GetSQLHighlighted(sourceCode string) string {
	code := []byte(sourceCode)

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(sql.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, code)
	if err != nil {
		fmt.Println("error parsing the code")
	}
	defer tree.Close()

	htmlParts := populateSliceWithNodeData(tree.RootNode(), code)
	return strings.Join(htmlParts, "")
}
