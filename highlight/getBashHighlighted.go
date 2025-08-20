package highlight

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	bash "github.com/smacker/go-tree-sitter/bash"
)

func GetBashHighlighted(sourceCode string) string {
	code := []byte(sourceCode)

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(bash.GetLanguage())

	tree := parser.Parse(nil, code)
	defer tree.Close()

	htmlParts := populateSliceWithNodeData(tree.RootNode(), code)
	return strings.Join(htmlParts, "")
}
