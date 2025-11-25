//go:build tools
// +build tools

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	fmt.Println("Generating RBAC HTML documentation...")
	docuDir := "web/docu"
	if err := os.MkdirAll(docuDir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", docuDir, err)
		os.Exit(1)
	}
	if err := convertMarkdownToHTML("docs/rbac.md", docuDir+"/rbac.html", "RBAC Permissions"); err != nil {
		fmt.Printf("Error converting rbac.md to HTML: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("RBAC HTML documentation generated successfully!")
}

func convertMarkdownToHTML(inputFile, outputFile, title string) error {
	md, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read markdown file: %w", err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Tables
	p := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	htmlBytes := markdown.ToHTML(md, p, renderer)

	fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - kaf-mirror Documentation</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
    <style>
        body { padding-top: 5rem; }
        .container { max-width: 1200px; }
        .navbar { border-bottom: 1px solid #ddd; }
        h1 { margin-bottom: 30px; color: #333; }
        h2 { 
            margin-top: 40px; 
            margin-bottom: 20px; 
            color: #495057;
            border-bottom: 2px solid #007bff;
            padding-bottom: 8px;
        }
        table { 
            margin: 25px 0;
            border-collapse: collapse;
            width: 100%%;
            background: white;
            border: 1px solid #dee2e6;
            border-radius: 6px;
            overflow: hidden;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            table-layout: fixed;
        }
        th { 
            background-color: #f8f9fa;
            color: #495057;
            font-weight: 600;
            padding: 15px 20px;
            text-align: left;
            border-bottom: 2px solid #dee2e6;
            font-size: 14px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        th:first-child { 
            width: 60%%;
        }
        th:last-child { 
            width: 40%%;
            text-align: left;
        }
        td { 
            padding: 12px 20px;
            border-bottom: 1px solid #dee2e6;
            vertical-align: top;
        }
        td:first-child { 
            text-align: left;
        }
        td:last-child { 
            text-align: left;
        }
        tr:hover {
            background-color: #f8f9fa;
        }
        code {
            background-color: #f1f3f4;
            color: #d73a49;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 13px;
            font-weight: 500;
        }
        tbody tr:last-child td {
            border-bottom: none;
        }
    </style>
</head>
<body>
    <nav class="navbar navbar-expand-md navbar-light bg-light fixed-top">
        <div class="container">
            <a class="navbar-brand" href="/docu/index.html">kaf-mirror Documentation</a>
            <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarsExampleDefault" aria-controls="navbarsExampleDefault" aria-expanded="false" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarsExampleDefault">
                <ul class="navbar-nav ml-auto">
                    <li class="nav-item">
                        <a class="nav-link" href="/docu/index.html">Home</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="/docu/cli-commands.html">CLI Commands</a>
                    </li>
                    <li class="nav-item active">
                        <a class="nav-link" href="/docu/rbac.html">RBAC<span class="sr-only">(current)</span></a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="/docu/swagger/">API (Swagger)</a>
                    </li>
                </ul>
            </div>
        </div>
    </nav>
    <div class="container">
        %s
    </div>
    <script src="https://code.jquery.com/jquery-3.5.1.slim.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.5.4/dist/umd/popper.min.js"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/js/bootstrap.min.js"></script>
</body>
</html>`, title, string(htmlBytes))

	if err := ioutil.WriteFile(outputFile, []byte(fullHTML), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
	return nil
}
