package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/pumenis/source-highlight/highlight"
)

//go:embed sqlsyntax.css
var sqlFile embed.FS

func main() {
	filePath := os.Args[1]

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return
	}
	sourceCode := string(content)

	css, err := sqlFile.ReadFile("sqlsyntax.css")
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return
	}

	result := highlight.GetSQLHighlighted(sourceCode)
	fmt.Println(`<!DOCTYPE html>
	  <html><head>
	 	<style>` + string(css) + `</style> 
		</head><div><pre>` + result + `</pre></div></html>`)
}
