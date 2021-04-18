package post

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/micheam/go-docbase"
)

func Test_postSummaryaaa(t *testing.T) {
	want := "Titile of post #tag-A #tag-B #tag-C"
	post := docbase.Post{
		Title: "Titile of post",
		Tags: []docbase.Tag{
			{"tag-A"}, {"tag-B"}, {"tag-C"},
		},
	}
	got := summary(post)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("postSummary mismatch (-want, +got):%s\n", diff)
	}
}

func Test_postSummary_emptyTags(t *testing.T) {
	want := "Titile of post"
	post := docbase.Post{
		Title: "Titile of post",
		Tags:  nil,
	}
	got := summary(post)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("postSummary mismatch (-want, +got):%s\n", diff)
	}
}

func Test_postSummary(t *testing.T) {
	tests := []struct {
		name string
		post docbase.Post
		want string
	}{
		{
			name: "empty post",
			want: "",
		},
		{
			name: "simple",
			want: "Titile of post",
			post: docbase.Post{
				Title: "Titile of post",
				Tags:  nil,
			},
		},
		{
			name: "with single tag",
			want: "Titile of post #tag-A",
			post: docbase.Post{
				Title: "Titile of post",
				Tags: []docbase.Tag{
					{"tag-A"},
				},
			},
		},
		{
			name: "with multi tags",
			want: "Titile of post #tag-A #tag-B",
			post: docbase.Post{
				Title: "Titile of post",
				Tags: []docbase.Tag{
					{"tag-A"}, {"tag-B"},
				},
			},
		},
		{
			name: "archived",
			want: "[archived] Titile of post",
			post: docbase.Post{
				Title:    "Titile of post",
				Archived: true,
			},
		},
		{
			name: "scope private",
			want: "[private] Titile of post",
			post: docbase.Post{
				Title: "Titile of post",
				Scope: docbase.ScopePrivate,
			},
		},
		{
			name: "scope group",
			want: "Titile of post",
			post: docbase.Post{
				Title: "Titile of post",
				Scope: docbase.ScopeGroup,
			},
		},
		{
			name: "complexed",
			want: "[archived][private] Titile of post #tag-A #tag-B",
			post: docbase.Post{
				Title:    "Titile of post",
				Archived: true,
				Scope:    docbase.ScopePrivate,
				Tags: []docbase.Tag{
					{"tag-A"}, {"tag-B"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := summary(tt.post)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("postSummary mismatch (-want, +got):%s\n", diff)
			}
		})
	}
}
