# docbase-cli

[docbase.io](https://docbase.io/) の記事を操作するコマンドラインツール

## Installation

### goコマンド でのインストール

(Go 1.16以降) `go install` でバイナリをインストール

```bash
$ go install github.com/micheam/docbase-cli/cmd/docbase@latest
```

(Go 1.15 以前) `go get` でバイナリをインストール

```bash
$ go get github.com/micheam/docbase-cli/...
```

### ソースコードからビルド

```bash
$ git clone github.com/micheam/docbase-cli
$ cd docbase-cli
$ make install
```

## Usage

```console
$ docbase [global options] command [command options] [arguments...]
```

## Commands

```
view     show post title and body
list     Search and list posts on docbase.io
new      Create new post.
edit     edit specified post.
tags     Show tags of group
help, h  Shows a list of commands or help for one command
```

## Global Options

```
--verbose, --vv       (default: false) [$DOCBASE_VERBOSE, $DOCBASE_DEBUG, $DEBUG]
--token ACCESS_TOKEN  ACCESS_TOKEN for docbase API [$DOCBASE_TOKEN]
--domain NAME         NAME on docbase.io [$DOCBASE_DOMAIN]
--help, -h            show help (default: false)
--version, -v         print the version (default: false)
```

## License
[MIT](./LICENSE)

## Author
Michito Maeda <https://github.com/micheam>
