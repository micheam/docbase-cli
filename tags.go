package docbasecli

import (
	"context"
	"fmt"

	"github.com/micheam/go-docbase"
)

type ListTagsRequest struct {
	Domain string
}

type TagCollectionPresenter func(ctx context.Context, tags []docbase.Tag) error

func ListTags(ctx context.Context, req ListTagsRequest, presenter TagCollectionPresenter) error {
	tags, err := docbase.ListTags(ctx, req.Domain)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	return presenter(ctx, tags)
}
