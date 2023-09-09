package main

import (
	"bytes"
	"fmt"
	"github.com/danaketh/ga-confluence-gen/confluence"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/virtomize/confluence-go-api"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func CollectMarkdownFiles(roots []string) (map[string][]string, error) {
	mdFiles := make(map[string][]string)

	for _, root := range roots {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				dir := filepath.Dir(path)
				mdFiles[dir] = append(mdFiles[dir], info.Name())
			}
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("error walking the path %s: %v", root, err)
		}
	}
	return mdFiles, nil
}

func main() {
	if len(os.Args) < 6 {
		panic("Usage: app <confluence-domain> <confluence-user> <confluence-token> <confluence-space> <source-path>")
	}

	var confluenceDomain = os.Args[1]
	var confluenceUser = os.Args[2]
	var confluenceToken = os.Args[3]
	var confluenceSpace = os.Args[4]
	var sourcePath = strings.Split(os.Args[5], "\n")

	// Configure API client
	//goconfluence.SetDebug(true)
	api, err := goconfluence.NewAPI(fmt.Sprintf("https://%s.atlassian.net/wiki/rest/api", confluenceDomain), confluenceUser, confluenceToken)
	if err != nil {
		panic(err)
	}

	// Collect all Markdown files
	fmt.Println("Looking for Markdown files...")
	mdFiles, err := CollectMarkdownFiles(sourcePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a new Markdown parser with the frontmatter extension
	md := goldmark.New(
		goldmark.WithExtensions(&frontmatter.Extender{
			Mode: frontmatter.SetMetadata,
		}),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
		),
	)

	// Print the map of Markdown files
	fmt.Println("Processing files...")
	for dir, files := range mdFiles {
		fmt.Printf("Directory: %s\n", dir)
		for _, file := range files {
			fmt.Printf("- %s", file)

			data, err := os.ReadFile(filepath.Join(dir, file))
			if err != nil {
				fmt.Println("\nerror reading file:", err)
				continue
			}

			documentRoot := md.Parser().Parse(text.NewReader(data))
			document := documentRoot.OwnerDocument()
			meta := document.Meta()
			// Skip files that don't have the confluence metadata
			pageID, ok := meta["confluence"]
			if !ok {
				fmt.Printf(" [SKIP]\n")
				continue
			}

			var buf bytes.Buffer
			if err := md.Convert(data, &buf); err != nil {
				fmt.Println("\nerror converting file:", err)
				continue
			}

			// Modify HTML with Confluence components
			var confluenceMarkup = confluence.ConvertToConfluence(buf.String())

			// Fetch the page from Confluence
			page, err := api.GetContentByID(fmt.Sprintf("%v", pageID), goconfluence.ContentQuery{
				SpaceKey: confluenceSpace,
				Expand:   []string{"version"},
			})
			if err != nil {
				log.Println(err)
				continue
			}

			// Table of contents
			_, ok = meta["toc"]
			if ok {
				confluenceMarkup = confluence.PrependTableOfContents(confluenceMarkup)
			}

			// Prepend and append warning messages
			_, ok = meta["ag-warning-pre"]
			if ok {
				confluenceMarkup = confluence.PrependWarningMessage(confluenceMarkup)
			}

			_, ok = meta["ag-warning-post"]
			if ok {
				confluenceMarkup = confluence.AppendWarningMessage(confluenceMarkup)
			}

			// Everything is good, let's publish to Confluence
			pageContent := &goconfluence.Content{
				ID:    fmt.Sprintf("%v", pageID),
				Type:  "page",
				Title: meta["title"].(string),
				Body: goconfluence.Body{
					Storage: goconfluence.Storage{
						Value:          confluenceMarkup,
						Representation: "storage",
					},
				},
				Space: &goconfluence.Space{
					Key: confluenceSpace,
				},
				Version: &goconfluence.Version{
					Number: page.Version.Number + 1,
				},
			}

			// Update the page
			_, err = api.UpdateContent(pageContent)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf(" [OK]\n")
		}
	}
}
