package tag

import (
	"context"
	"fmt"

	"github.com/micheam/go-docbase"
)

type ListRequest struct {
	Domain string
}

type ListResultHandler func(ctx context.Context, tags []docbase.Tag) error

func List(ctx context.Context, req ListRequest, handler ListResultHandler) error {
	tags, err := docbase.ListTags(ctx, req.Domain)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	return handler(ctx, tags)
}
