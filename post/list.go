package post

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/micheam/go-docbase"
)

type ListRequest struct {
	Query   *string
	Page    *int
	PerPage *int
	Domain  string
}

func List(ctx context.Context, req ListRequest, handle PostListHandler) error {
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
