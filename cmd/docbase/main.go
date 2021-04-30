package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	docbasecli "github.com/micheam/docbase-cli"
	"github.com/micheam/docbase-cli/pointer"
	"github.com/micheam/go-docbase"
	"github.com/urfave/cli/v2"
)

var (
	Version = "0.1.0"
	Githash = "devel"
)

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetOutput(io.Discard)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "docbase"
	app.Usage = "CLI Client for DocBase API"
	app.Version = fmt.Sprintf("%s (rev: %s)", Version, Githash)
	app.Authors = []*cli.Author{
		{
			Name:  "Michito Maeda",
			Email: "michito.maeda@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"vv"},
			EnvVars: []string{"DOCBASE_VERBOSE", "DOCBASE_DEBUG", "DEBUG"},
		},
		&cli.StringFlag{
			Name:    "token",
			Usage:   "`ACCESS_TOKEN` for docbase API",
			EnvVars: []string{"DOCBASE_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "domain",
			EnvVars: []string{"DOCBASE_DOMAIN"},
			Usage:   "`NAME` on docbase.io",
		},
	}
	app.Commands = []*cli.Command{
		viewPost, listPosts,
		newPost, editPost,
		tags,
	}
	return app
}

var viewPost = &cli.Command{
	Name:      "view",
	Usage:     "show post title and body",
	ArgsUsage: "POST_ID",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "web",
			Aliases: []string{"w"},
			Usage:   "Open a post in the browser",
		},
		&cli.IntFlag{
			Name:    "lines",
			Aliases: []string{"l"},
			Usage:   "`NUM` to display body. set 0 to display full.",
			Value:   0,
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		postID, err := docbase.ParsePostID(c.Args().First())
		if err != nil {
			return err
		}
		req := docbasecli.GetPostRequest{
			Domain: c.String("domain"),
			ID:     postID,
		}
		var handle = docbasecli.OutputPostDetail(os.Stdout, c.Int("lines"))
		if c.Bool("web") {
			handle = docbasecli.OpenBrowser
		}
		return docbasecli.GetPost(c.Context, req, handle)
	},
}

var listPosts = &cli.Command{
	Name:  "list",
	Usage: "Search and list posts on docbase.io",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "`options` to narrow down the search. ex: groups,contributors, etc.",
		},
		&cli.IntFlag{
			Name:    "page",
			Aliases: []string{"p"},
			Value:   1,
			Usage:   "`num` of posts on a page",
		},
		&cli.IntFlag{
			Name:    "per-page",
			Aliases: []string{"pp"},
			Value:   20,
			Usage:   "`num` of page",
		},
		&cli.BoolFlag{
			Name:    "meta",
			Aliases: []string{"m"},
			Usage:   "Display META-Fields (Total,Previous,Next) on footer",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		req := docbasecli.ListPostsRequest{Domain: c.String("domain")}
		if c.String("query") != "" {
			req.Query = pointer.StringPtr(c.String("query"))
		}
		if c.Int("page") != 0 {
			req.Page = pointer.IntPtr(c.Int("page"))
		}
		if c.Int("per-page") != 0 {
			req.PerPage = pointer.IntPtr(c.Int("per-page"))
		}
		presenter, err := docbasecli.BuildPostCollectionHandler(c.Bool("meta"))
		if err != nil {
			return err
		}
		return docbasecli.ListPosts(c.Context, req, presenter)
	},
}

// TODO(micheam): 設定ファイルで指定可能にする
var defaultTitle = func() string {
	now := time.Now()
	return fmt.Sprintf("%s 作業メモ", now.Format("2006-01-02"))
}

var newPost = &cli.Command{
	Name:      "new",
	Usage:     "Create new post.",
	ArgsUsage: "-",
	Flags: []cli.Flag{
		// TODO(micheam): option `--dradt`
		// TODO(micheam): option `--notice`
		// TODO(micheam): option `--tags`
		// TODO(micheam): option `--scope`
		// TODO(micheam): option `--groups`
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "`STR-VAL` for title",
			Value:   defaultTitle(),
		},
		&cli.StringFlag{
			Name:    "body",
			Aliases: []string{"b"},
			Usage:   "`STR-VAL` for body",
		},
		&cli.StringFlag{
			Name:  "body-file",
			Usage: "`PATH` of input file",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		req := docbasecli.CreatePostRequest{
			Title:  c.String("title"),
			Domain: c.String("domain"),
		}

		// Body
		if len(c.String("body")) != 0 {
			req.Body = strings.NewReader(c.String("body"))
		} else if len(c.String("body-file")) != 0 {
			filepath := c.String("body-file")
			file, err := os.Open(filepath)
			if err != nil {
				return fmt.Errorf("cant open %q: %w", filepath, err)
			}
			defer func() { _ = file.Close() }()
			req.Body = file
		} else {
			b, err := docbasecli.CaptureInputFromEditor(
				docbasecli.GetPreferredEditorFromEnvironment,
				nil,
			)
			if err != nil {
				return fmt.Errorf("faild to capture input: %w", err)
			}
			req.Body = bytes.NewReader(b)
		}

		presenter := func(ctx context.Context, post docbase.Post) error {
			fmt.Println(post.URL)
			return nil
		}
		return docbasecli.CreatePost(c.Context, req, presenter)
	},
}

var editPost = &cli.Command{
	Name:      "edit",
	Usage:     "edit specified post.",
	ArgsUsage: "ID",
	Flags: []cli.Flag{
		// TODO(micheam): option `--dradt`
		// TODO(micheam): option `--notice`
		// TODO(micheam): option `--tags`
		// TODO(micheam): option `--scope`
		// TODO(micheam): option `--groups`
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "`STR-VAL` for title",
			Value:   defaultTitle(),
		},
		&cli.StringFlag{
			Name:    "body",
			Aliases: []string{"b"},
			Usage:   "`STR-VAL` for body",
		},
		&cli.StringFlag{
			Name:  "body-file",
			Usage: "`PATH` of input file",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		if !c.Args().Present() {
			return errors.New("need to specify target post id")
		}
		id, err := docbase.ParsePostID(c.Args().First())
		if err != nil {
			return fmt.Errorf("illegal post id: %w", err)
		}
		req := docbasecli.UpdatePostRequest{
			Domain: c.String("domain"),
			ID:     id,
		}

		// Get existing post
		var existing docbase.Post
		{
			r := docbasecli.GetPostRequest{
				Domain: c.String("domain"),
				ID:     id,
			}
			var h = func(_ context.Context, post docbase.Post) error {
				existing = post
				return nil
			}
			err := docbasecli.GetPost(c.Context, r, h)
			if err != nil {
				return fmt.Errorf("faild to get existing post(%d): %w", id, err)
			}
		}

		// Body
		if len(c.String("body")) != 0 {
			req.Body = strings.NewReader(c.String("body"))
		} else if len(c.String("body-file")) != 0 {
			filepath := c.String("body-file")
			file, err := os.Open(filepath)
			if err != nil {
				return fmt.Errorf("cant open %q: %w", filepath, err)
			}
			defer func() { _ = file.Close() }()
			req.Body = file
		} else {
			b, err := docbasecli.CaptureInputFromEditor(
				docbasecli.GetPreferredEditorFromEnvironment,
				[]byte(existing.Body),
			)
			if err != nil {
				return fmt.Errorf("faild to capture input: %w", err)
			}
			req.Body = bytes.NewReader(b)
		}
		h := func(ctx context.Context, post docbase.Post) error {
			fmt.Println("Updated.")
			fmt.Println(post.URL)
			return nil
		}
		return docbasecli.UpatePost(c.Context, req, h)
	},
}

var tags = &cli.Command{
	Name:  "tags",
	Usage: "Show tags of group",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		req := docbasecli.ListTagsRequest{
			Domain: c.String("domain"),
		}
		presenter := func(ctx context.Context, tags []docbase.Tag) error {
			for _, tag := range tags {
				fmt.Println(tag.Name)
			}
			return nil
		}
		return docbasecli.ListTags(c.Context, req, presenter)
	},
}
