# source-highlight a go module

Golang module for highlighting sources that outputs html

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/pumenis/source-highlight/highlight"
)

func main() {
	filePath := os.Args[1]

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return
	}
	sourceCode := string(content)

	result := highlight.GetGoHighlighted(sourceCode)
	fmt.Println(`<!DOCTYPE html>
	  <html><head>
	  <link rel="stylesheet" href="gosyntax.css" />
		</head><div><pre>` + result + `</pre></div></html>`)
}
```

```
go run ./examples/highlight_go_file/ highlight/getGoHighlighted.go >index.html

```
