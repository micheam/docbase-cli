package post

import (
	"strings"

	"github.com/micheam/go-docbase"
)

func summary(post docbase.Post) string {
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
