package post

import (
	"context"
	"log"

	"github.com/micheam/go-docbase"
)

type GetRequest struct {
	ID     docbase.PostID
	Domain string
}

func Get(ctx context.Context, req GetRequest, handle PostHandler) error {
	log.Printf("get post with req: %v", req)
	post, err := docbase.GetPost(ctx, req.Domain, req.ID)
	if err != nil {
		return err
	}
	return handle(ctx, *post)
}
