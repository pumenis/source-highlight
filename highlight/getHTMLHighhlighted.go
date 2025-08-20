package highlight

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode/utf8"

	sitter "github.com/smacker/go-tree-sitter"
	htmlLang "github.com/smacker/go-tree-sitter/html"
)

// Precompiled regex patterns for performance
var (
	matchUnderscoreLowerLetters = regexp.MustCompile(`^[a-z_]+$`)
	matchScriptAttributes       = regexp.MustCompile(`^(x-|on|@).*$`)
)

func populateHTMLSliceWithNodeData(node *sitter.Node, code []byte) []string {
	var htmlParts []string

	// Handle leading text before root node
	if node.Parent() == nil && node.StartByte() > 0 {
		htmlParts = append(htmlParts, string(code[0:node.StartByte()]))
	}

	// Determine class name
	class := "syntax_node"
	if matchUnderscoreLowerLetters.MatchString(node.Type()) {
		class = node.Type()
	}

	isNamed := "false"
	if node.IsNamed() {
		isNamed = "true"
	}

	htmlParts = append(htmlParts, fmt.Sprintf(
		`<span class="%s" type="%s" is_named="%s">`,
		class, node.Type(), isNamed))

	// Handle text between parent and first child
	if node.ChildCount() > 0 && node.StartByte() < node.Child(0).StartByte() {
		htmlParts = append(htmlParts, string(code[node.StartByte():node.Child(0).StartByte()]))
	}

	// Leaf node handling
	if node.ChildCount() == 0 {
		text := node.Content(code)
		if !utf8.ValidString(text) {
			text = string(code[node.StartByte():node.EndByte()])
		}

		if (node.Type() == "raw_text" && node.Parent() != nil && node.Parent().Type() == "script_element") ||
			(node.Type() == "attribute_value" &&
				node.Parent() != nil && node.Parent().Parent() != nil &&
				matchScriptAttributes.MatchString(node.Parent().Parent().Child(0).Content(code))) {
			htmlParts = append(htmlParts, GetJSHighlighted(text))
		} else {
			htmlParts = append(htmlParts, html.EscapeString(text))
		}
	}

	// Recursively process children
	for i := uint32(0); i < node.ChildCount(); i++ {
		if i > 0 && node.Child(i).StartByte() > node.Child(i-1).EndByte() {
			htmlParts = append(htmlParts, string(code[node.Child(i-1).EndByte():node.Child(i).StartByte()]))
		}
		htmlParts = append(htmlParts, populateHTMLSliceWithNodeData(node.Child(i), code)...)
	}

	// Handle trailing text after last child
	if node.ChildCount() > 0 && node.EndByte() > node.Child(node.ChildCount()-1).EndByte() {
		htmlParts = append(htmlParts, string(code[node.Child(node.ChildCount()-1).EndByte():node.EndByte()]))
	}

	htmlParts = append(htmlParts, "</span>")
	return htmlParts
}

func GetHTMLHighlighted(sourceCode string) string {
	code := []byte(sourceCode)

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(htmlLang.GetLanguage())

	tree := parser.Parse(nil, code)
	defer tree.Close()

	htmlParts := populateHTMLSliceWithNodeData(tree.RootNode(), code)
	return strings.Join(htmlParts, "")
}
