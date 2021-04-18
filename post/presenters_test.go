package post

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/micheam/go-docbase"
)

func timeMust(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}

func Test_marshal(t *testing.T) {
	data := struct {
		F1 string   `yaml:"title"`
		F2 int      `yaml:"id"`
		F3 []string `yaml:"tags"`
		F4 string   `yaml:"-"`
		F5 string   `yaml:"foo,omitempty"`
	}{
		F1: "Hello, Yaml World.",
		F2: 33224455,
		F3: []string{"dev", "golang"},
	}
	want := `title: Hello, Yaml World.
id: 33224455
tags:
- dev
- golang
`
	got := marshal(data)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):%s\n", diff)
	}
}

func TestWritePost(t *testing.T) {
	var buf = new(bytes.Buffer)
	sut := WritePost(buf, 5)
	ctx := context.Background()
	post := docbase.Post{
		ID:        11111,
		Title:     "Title For Test",
		Tags:      []docbase.Tag{{Name: "tag1"}, {Name: "tag2"}},
		CreatedAt: "1998-02-09T11:12:13",
		UpdatedAt: "1998-02-09T11:12:13",
		Body: `aaaaaaaaaaaaaaaaaaaa
bbbbbbbbbbbbbbbbbbbb
cccccccccccccccccccc
dddddddddddddddddddd
eeeeeeeeeeeeeeeeeeee
ffffffffffffffffffff
gggggggggggggggggggg`,
	}
	err := sut(ctx, post)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}
	want := `---
ID:        11111
Title:     Title For Test
Tags:      #tag1 #tag2 
CreatedAt: 1998-02-09T11:12:13
UpdatedAt: 1998-02-09T11:12:13
Draft:     false
Archived:  false
---

aaaaaaaaaaaaaaaaaaaa
bbbbbbbbbbbbbbbbbbbb
cccccccccccccccccccc
dddddddddddddddddddd
eeeeeeeeeeeeeeeeeeee


Showed 5 of 7
`
	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("result text mismatch (-want, +got):%s\n", diff)
	}
}
