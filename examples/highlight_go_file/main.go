package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/pumenis/source-highlight/highlight"
)

//go:embed gosyntax.css
var goFile embed.FS

func main() {
	filePath := os.Args[1]

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return
	}
	sourceCode := string(content)

	css, err := goFile.ReadFile("gosyntax.css")
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return
	}

	result := highlight.GetGoHighlighted(sourceCode)
	fmt.Println(`<!DOCTYPE html>
	  <html><head>
	 	<style>` + string(css) + `</style> 
		</head><div><pre>` + result + `</pre></div></html>`)
}
