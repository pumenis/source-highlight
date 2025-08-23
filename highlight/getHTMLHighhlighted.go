package highlight

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

	if node.Parent() == nil && node.StartByte() > 0 {
		htmlParts = append(htmlParts, html.EscapeString(string(code[0:node.StartByte()])))
	}

	class := "syntax_node"
	if matchUnderscoreLowerLetters.MatchString(node.Type()) {
		class = node.Type()
	}

	isNamed := "false"
	if node.IsNamed() {
		isNamed = "true"
	}

	text := node.Content(code)
	if !utf8.ValidString(text) {
		text = string(code[node.StartByte():node.EndByte()])
	}
	hash := sha256.Sum256([]byte(text))
	id := hex.EncodeToString(hash[:])[:12] // 12-char hash for brevity

	htmlParts = append(htmlParts, fmt.Sprintf(
		`<span id="h-%s" class="%s" type="%s" is_named="%s">`,
		id, class, strings.ReplaceAll(node.Type(), `"`, "&quot;"), isNamed))

	if node.ChildCount() > 0 && node.StartByte() < node.Child(0).StartByte() {
		htmlParts = append(htmlParts, html.EscapeString(string(code[node.StartByte():node.Child(0).StartByte()])))
	}

	if node.ChildCount() == 0 {
		if (node.Type() == "raw_text" && node.Parent() != nil && node.Parent().Type() == "script_element") ||
			(node.Type() == "attribute_value" &&
				node.Parent() != nil && node.Parent().Parent() != nil &&
				matchScriptAttributes.MatchString(node.Parent().Parent().Child(0).Content(code))) {
			htmlParts = append(htmlParts, GetJSHighlighted(text))
		} else if node.Type() == "raw_text" && node.Parent() != nil && node.Parent().Type() == "style_element" {
			htmlParts = append(htmlParts, GetCSSHighlighted(text))
		} else {
			htmlParts = append(htmlParts, html.EscapeString(text))
		}
	}

	for i := uint32(0); i < node.ChildCount(); i++ {
		intI := int(i)
		if i > 0 && node.Child(intI).StartByte() > node.Child(intI-1).EndByte() {
			htmlParts = append(htmlParts, html.EscapeString(string(code[node.Child(intI-1).EndByte():node.Child(intI).StartByte()])))
		}
		htmlParts = append(htmlParts, populateHTMLSliceWithNodeData(node.Child(intI), code)...)
	}

	if node.ChildCount() > 0 && node.EndByte() > node.Child(int(node.ChildCount()-1)).EndByte() {
		htmlParts = append(htmlParts, html.EscapeString(string(code[node.Child(int(node.ChildCount()-1)).EndByte():node.EndByte()])))
	}

	htmlParts = append(htmlParts, "</span>")
	return htmlParts
}

func GetHTMLHighlighted(sourceCode string) string {
	code := []byte(sourceCode)

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(htmlLang.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, code)
	if err != nil {
		fmt.Println("cannot parse code")
	}
	defer tree.Close()

	htmlParts := populateHTMLSliceWithNodeData(tree.RootNode(), code)
	return strings.Join(htmlParts, "")
}
