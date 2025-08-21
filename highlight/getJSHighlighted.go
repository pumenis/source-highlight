package highlight

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	javascript "github.com/smacker/go-tree-sitter/javascript"
)

func GetJSHighlighted(sourceCode string) string {
	code := []byte(sourceCode)

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(javascript.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, code)
	if err != nil {
		fmt.Println("cannot parse code")
	}
	defer tree.Close()

	htmlParts := populateSliceWithNodeData(tree.RootNode(), code)
	return strings.Join(htmlParts, "")
}
