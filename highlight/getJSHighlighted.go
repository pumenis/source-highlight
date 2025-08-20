package highlight

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	javascript "github.com/smacker/go-tree-sitter/javascript"
)

func GetJSHighlighted(sourceCode string) string {
	code := []byte(sourceCode)

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(javascript.GetLanguage())

	tree := parser.Parse(nil, code)
	defer tree.Close()

	htmlParts := populateSliceWithNodeData(tree.RootNode(), code)
	return strings.Join(htmlParts, "")
}
