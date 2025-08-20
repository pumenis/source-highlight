package highlight

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	golang "github.com/smacker/go-tree-sitter/go"
)

func GetGoHighlighted(sourceCode string) string {
	code := []byte(sourceCode)

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(golang.GetLanguage())

	tree := parser.Parse(nil, code)
	defer tree.Close()

	htmlParts := populateSliceWithNodeData(tree.RootNode(), code)
	return strings.Join(htmlParts, "")
}

