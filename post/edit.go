package post

import (
	"context"
	"fmt"
	"io"

	"github.com/micheam/go-docbase"
)

type UpdateRequest struct {
	Domain string
	ID     docbase.PostID
	Body   io.Reader
}

func Upate(ctx context.Context, req UpdateRequest, handle PostHandler) error {
	updated, err := docbase.UpdatePost(ctx, req.Domain, req.ID, req.Body, docbase.UpdateFields{})
	if err != nil {
		return fmt.Errorf("failed to create new post: %w", err)
	}
	return handle(ctx, *updated)
}
