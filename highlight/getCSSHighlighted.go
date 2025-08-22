// Package highlight provides syntax highlighting
// for Go, HTML, JS, Bash source codes using tree-sitter.
package highlight

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/css"
)

func hexToRGB(hex string) (r, g, b int64) {
	if hex[0] == '#' {
		hex = hex[1:]
	}
	r, _ = strconv.ParseInt(hex[0:2], 16, 64)
	g, _ = strconv.ParseInt(hex[2:4], 16, 64)
	b, _ = strconv.ParseInt(hex[4:6], 16, 64)
	return
}

func getContrastColor(hex string) string {
	r, g, b := hexToRGB(hex)
	// Calculate luminance
	luminance := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	if luminance > 128 {
		return "#000000" // dark text
	}
	return "#FFFFFF" // light text
}

func populateCSSSliceWithNodeData(node *sitter.Node, code []byte) []string {
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

	colorStyle := ""
	if node.Type() == "color_value" {
		colorStyle = ` style="background-color:` + content + `;color:` + getContrastColor(content) + `;"`
	}
	htmlParts = append(htmlParts, fmt.Sprintf(
		`<span id="h-%s" class="%s" type="%s" is_named="%s"%s>`,
		id, class, node.Type(), isNamed, colorStyle))

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
		htmlParts = append(htmlParts, populateCSSSliceWithNodeData(node.Child(intI), code)...)
	}

	if node.ChildCount() != 0 && node.EndByte() > node.Child(int(node.ChildCount()-1)).EndByte() {
		htmlParts = append(htmlParts, html.EscapeString(string(code[node.Child(int(node.ChildCount()-1)).EndByte():node.EndByte()])))
	}

	htmlParts = append(htmlParts, "</span>")
	return htmlParts
}

func GetCSSHighlighted(sourceCode string) string {
	code := []byte(sourceCode)

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(css.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, code)
	if err != nil {
		fmt.Println("error parsing the code")
	}
	defer tree.Close()

	htmlParts := populateCSSSliceWithNodeData(tree.RootNode(), code)
	return strings.Join(htmlParts, "")
}
