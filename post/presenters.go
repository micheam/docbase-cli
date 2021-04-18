package post

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

	"github.com/micheam/docbase-cli/text"
	"github.com/micheam/go-docbase"
	"gopkg.in/yaml.v2"
)

func marshal(v interface{}) string {
	a, _ := yaml.Marshal(v)
	return string(a)
}

// define ResultHandlers
type (
	PostHandler     func(ctx context.Context, p docbase.Post) error
	PostListHandler func(ctx context.Context, ps []docbase.Post, m docbase.Meta) error
)

func WritePost(out io.Writer, n int) PostHandler {
	type M struct {
		docbase.Post
		Total int
		Lines []string
	}
	const tmplPostDetail = `---
ID:        {{.ID}}
Title:     {{.Title}}
Tags:      {{range .Tags}}{{- printf "#%s " .Name }}{{end}}
CreatedAt: {{.CreatedAt}}
UpdatedAt: {{.UpdatedAt}}
Draft:     {{.Draft}}
Archived:  {{.Archived}}
---

{{range .Lines}}
  {{- .}}
{{end}}

Showed {{len .Lines}} of {{.Total}}
`
	funcMap := template.FuncMap{}
	tmpl, err := template.New("get-post").Funcs(funcMap).
		Parse(tmplPostDetail)
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, post docbase.Post) error {
		// TODO(micheam): Win対応
		lines := strings.Split(text.Dos2Unix(post.Body), "\n")
		total := len(lines)
		if n > 0 {
			if n > len(lines) {
				n = len(lines)
			}
			lines = lines[:n]
		}
		err = tmpl.Execute(out, M{
			Post:  post,
			Total: total,
			Lines: lines,
		})
		if err != nil {
			return err
		}
		return nil
	}
}

func OpenBrowser(_ context.Context, post docbase.Post) error {
	openbrowser(post.URL)
	return nil
}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func BuildListResultHandler(withMeta bool) (PostListHandler, error) {
	const _tmplPostsList = `{{range .}}{{printf "%d\t%s" .ID (summary .)}}{{"\n"}}{{end}}`
	tmplPostsList, err := template.New("list-posts").Funcs(template.FuncMap{
		"summary": summary,
	}).Parse(_tmplPostsList)
	if err != nil {
		return nil, err
	}
	const _tmplMetaData = `---
Total: {{.Total}}
{{with .NextPageURL}}Next: {{.}}{{"\n"}}{{end -}} 
{{with .PreviousPageURL}}Prev: {{.}}{{end -}}`
	tmplMetaData, err := template.New("meta").Parse(_tmplMetaData)
	if err != nil {
		return nil, err
	}
	if withMeta {
		return func(ctx context.Context, posts []docbase.Post, meta docbase.Meta) error {
			err := tmplPostsList.Execute(os.Stdout, posts)
			if err != nil {
				return err
			}
			return tmplMetaData.Execute(os.Stdout, meta)
		}, nil
	}
	return func(ctx context.Context, posts []docbase.Post, _ docbase.Meta) error {
		return tmplPostsList.Execute(os.Stdout, posts)
	}, nil
}
