package docbasecli

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/micheam/docbase-cli/pointer"
	"github.com/micheam/docbase-cli/text"
	"github.com/micheam/go-docbase"
	"gopkg.in/yaml.v2"
)

/***************************************
 * Update Post
 ***************************************/

type UpdatePostRequest struct {
	Domain string
	ID     docbase.PostID
	Body   io.Reader
}

func UpatePost(ctx context.Context, req UpdatePostRequest, handle PostHandler) error {
	updated, err := docbase.UpdatePost(ctx, req.Domain, req.ID, req.Body, docbase.UpdateFields{})
	if err != nil {
		return fmt.Errorf("failed to create new post: %w", err)
	}
	return handle(ctx, *updated)
}

/***************************************
 * Get Post
 ***************************************/

type GetPostRequest struct {
	ID     docbase.PostID
	Domain string
}

func GetPost(ctx context.Context, req GetPostRequest, handle PostHandler) error {
	log.Printf("get post with req: %v", req)
	post, err := docbase.GetPost(ctx, req.Domain, req.ID)
	if err != nil {
		return err
	}
	return handle(ctx, *post)
}

type ListPostsRequest struct {
	Query   *string
	Page    *int
	PerPage *int
	Domain  string
}

func ListPosts(ctx context.Context, req ListPostsRequest, handle PostCollectionHandler) error {
	param := url.Values{}
	if req.Query != nil {
		param.Add("q", *req.Query)
	}
	if req.Page != nil {
		param.Add("page", fmt.Sprint(*req.Page))
	}
	if req.PerPage != nil {
		param.Add("per_page", fmt.Sprint(*req.PerPage))
	}

	log.Printf("list posts with req: %v", req)

	posts, meta, err := docbase.ListPosts(ctx, req.Domain, param)
	if err != nil {
		return err
	}
	return handle(ctx, posts, *meta)
}

type CreatePostRequest struct {
	Domain string
	Title  string
	Body   io.Reader

	// Option メモ作成時のオプション
	// 省略した場合は DefaultPostOption が適用される
	Option *docbase.PostOption
}

// DefaultPostOption メモ作成時のデフォルトオプション
var DefaultPostOption = docbase.PostOption{
	Draft:  pointer.BoolPtr(true),
	Tags:   []string{},
	Scope:  "private",
	Groups: []int{},
}

func CreatePost(ctx context.Context, req CreatePostRequest, handler PostHandler) error {
	opt := DefaultPostOption
	if req.Option != nil {
		opt = *req.Option
	}
	created, err := docbase.NewPost(ctx, req.Domain, req.Title, req.Body, opt)
	if err != nil {
		return fmt.Errorf("failed to create new post: %w", err)
	}
	return handler(ctx, *created)
}

func marshal(v interface{}) string {
	a, _ := yaml.Marshal(v)
	return string(a)
}

// define ResultHandlers
type (
	PostHandler           func(ctx context.Context, p docbase.Post) error
	PostCollectionHandler func(ctx context.Context, ps []docbase.Post, m docbase.Meta) error
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

func BuildPostCollectionHandler(withMeta bool) (PostCollectionHandler, error) {
	const _tmplPostsList = `{{range .}}{{printf "%d\t%s" .ID (summary .)}}{{"\n"}}{{end}}`
	tmplPostsList, err := template.New("list-posts").Funcs(template.FuncMap{
		"summary": summarizePost,
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

// メモを要約した文字列を生成する
func summarizePost(post docbase.Post) string {
	sb := new(strings.Builder)
	var prefixed bool
	if post.Archived {
		prefixed = true
		_, _ = sb.Write([]byte("[archived]"))
	}
	if post.Scope == docbase.ScopePrivate {
		prefixed = true
		_, _ = sb.Write([]byte("[" + post.Scope + "]"))
	}
	if prefixed {
		_, _ = sb.Write([]byte(" "))
	}
	_, _ = sb.Write([]byte(post.Title))
	for i := range post.Tags {
		tag := post.Tags[i]
		_, _ = sb.Write([]byte(" #" + tag.Name))
	}
	return sb.String()
}
