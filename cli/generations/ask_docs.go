package generations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/v2/cli/services/indexing"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
)

func AskDocs(_ AskDocsOptions) (err error) {
	var pages map[string]indexing.IndexedPage
	if pages, err = indexing.Index(indexing.IndexOptions{
		Address:     "https://razshare.github.io/frizzante-docs",
		Context:     context.Background(),
		StickToHost: true,
		Depth:       1,
	}); err != nil {
		return
	}
	directoryName := filepath.Join(".gen")
	if !files.IsDirectory(directoryName) {
		if err = os.Mkdir(directoryName, os.ModePerm); err != nil {
			return
		}
	}
	var builder strings.Builder
	for address, page := range pages {
		builder.WriteString(fmt.Sprintf("<!-- Beginning of page %s -->\n", page.Title))
		builder.WriteString("<html>\n")
		builder.WriteString("    <head>\n")
		builder.WriteString("        <title>\n")
		if page.Title == "" {
			builder.WriteString(address + "\n")
		} else {
			builder.WriteString(page.Title + "\n")
		}
		builder.WriteString("        </title>\n")
		builder.WriteString("    </head>\n")
		builder.WriteString("    <body>\n")
		builder.WriteString(page.Body + "\n")
		builder.WriteString("    </body>\n")
		builder.WriteString("</html>\n")
		builder.WriteString(fmt.Sprintf("<!-- Ending of page %s -->\n\n", page.Title))
	}
	fileName := filepath.Join(directoryName, "ask.docs.md")
	if err = os.WriteFile(fileName, []byte(builder.String()), os.ModePerm); err != nil {
		return
	}
	return
}
