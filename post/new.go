package post

import (
	"context"
	"fmt"
	"io"

	"github.com/micheam/docbase-cli/pointer"
	"github.com/micheam/go-docbase"
)

type CreateRequest struct {
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

type CreationResultHandler func(ctx context.Context, created *docbase.Post) error

func Create(ctx context.Context, req CreateRequest, handler CreationResultHandler) error {
	opt := DefaultPostOption
	if req.Option != nil {
		opt = *req.Option
	}
	created, err := docbase.NewPost(ctx, req.Domain, req.Title, req.Body, opt)
	if err != nil {
		return fmt.Errorf("failed to create new post: %w", err)
	}
	return handler(ctx, created)
}
