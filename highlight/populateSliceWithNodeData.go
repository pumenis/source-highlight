package highlight

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"regexp"
	"strings"

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
		`<span id="h-%s" class="%s" type="%s" is_named="%s">`,
		id, class, strings.ReplaceAll(node.Type(), `"`, "&quot;"), isNamed))

	if node.ChildCount() != 0 && node.StartByte() < node.Child(0).StartByte() {
		htmlParts = append(htmlParts, html.EscapeString(string(code[node.StartByte():node.Child(0).StartByte()])))
	}

	if node.ChildCount() == 0 {
		if node.Type() == "raw_string_literal" {
			if strings.HasPrefix(content, "`-- sql") {
				htmlParts = append(htmlParts, GetSQLHighlighted(content))
			}
			if strings.HasPrefix(content, "`<") {
				htmlParts = append(htmlParts, GetHTMLHighlighted(content))
			}
		} else {
			htmlParts = append(htmlParts, html.EscapeString(content))
		}
	}

	if node.Type() == "string" && strings.HasPrefix(content, `"-- sql`) {
		htmlParts = append(htmlParts, `<span id="h-8a331fdde703" class="syntax_node" type="&quot;" is_named="false">"</span>`)
		htmlParts = append(htmlParts, GetSQLHighlighted(content[1:len(content)-1]))
		htmlParts = append(htmlParts, `<span id="h-8a331fdde703" class="syntax_node" type="&quot;" is_named="false">"</span>`)
	} else {
		for i := uint32(0); i < node.ChildCount(); i++ {
			intI := int(i)
			if node.ChildCount() > 1 && i > 0 && node.Child(intI).StartByte() > node.Child(intI-1).EndByte() {
				htmlParts = append(htmlParts, html.EscapeString(string(code[node.Child(intI-1).EndByte():node.Child(intI).StartByte()])))
			}
			htmlParts = append(htmlParts, populateSliceWithNodeData(node.Child(intI), code)...)
		}
	}

	if node.ChildCount() != 0 && node.EndByte() > node.Child(int(node.ChildCount()-1)).EndByte() {
		htmlParts = append(htmlParts, html.EscapeString(string(code[node.Child(int(node.ChildCount()-1)).EndByte():node.EndByte()])))
	}

	htmlParts = append(htmlParts, "</span>")
	return htmlParts
}
