package highlight

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"regexp"

	sitter "github.com/smacker/go-tree-sitter"
)

func populateSliceWithNodeData(node *sitter.Node, code []byte) []string {
	htmlParts := []string{}

	if node.Parent() == nil && node.StartByte() > 0 {
		htmlParts = append(htmlParts, html.EscapeString(string(code[0:node.StartByte()])))
	}

	matchUnderscoreLowerLettersRe, err := regexp.Compile(`^[a-z_]+$`)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return nil
	}

	isNamed := "false"
	if node.IsNamed() {
		isNamed = "true"
	}
	class := "syntax_node"
	if matchUnderscoreLowerLettersRe.MatchString(node.Type()) {
		class = node.Type()
	}

	content := node.Content(code)
	hash := sha256.Sum256([]byte(content))
	id := hex.EncodeToString(hash[:])[:12]

	// Start span with metadata and ID
	htmlParts = append(htmlParts, fmt.Sprintf(
		`<span id="%s" class="%s" type="%s" is_named="%s">`,
		id, class, node.Type(), isNamed))

	if node.ChildCount() != 0 && node.StartByte() < node.Child(0).StartByte() {
		htmlParts = append(htmlParts, html.EscapeString(string(code[node.StartByte():node.Child(0).StartByte()])))
	}

	if node.ChildCount() == 0 {
		htmlParts = append(htmlParts, html.EscapeString(content))
	}

	for i := uint32(0); i < node.ChildCount(); i++ {
		intI := int(i)
		if node.ChildCount() > 1 && i > 0 && node.Child(intI).StartByte() > node.Child(intI-1).EndByte() {
			htmlParts = append(htmlParts, html.EscapeString(string(code[node.Child(intI-1).EndByte():node.Child(intI).StartByte()])))
		}
		htmlParts = append(htmlParts, populateSliceWithNodeData(node.Child(intI), code)...)
	}

	if node.ChildCount() != 0 && node.EndByte() > node.Child(int(node.ChildCount()-1)).EndByte() {
		htmlParts = append(htmlParts, html.EscapeString(string(code[node.Child(int(node.ChildCount()-1)).EndByte():node.EndByte()])))
	}

	htmlParts = append(htmlParts, "</span>")
	return htmlParts
}
