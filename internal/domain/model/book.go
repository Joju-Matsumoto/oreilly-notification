package model

import (
	"bytes"
	"net/url"
	"strings"
	"time"

	"github.com/Joju-Matsumoto/oreilly-notification/pkg/oreillyapi"
	"golang.org/x/net/html"
)

type Book struct {
	Result oreillyapi.Result
}

func NewBookFromOreillyResult(result oreillyapi.Result) *Book {
	return &Book{
		Result: result,
	}
}

func (b *Book) ID() string {
	return b.Result.Id
}

func (b *Book) Title() string {
	return b.Result.Title
}

func (b *Book) Description() string {
	doc, _ := html.Parse(strings.NewReader(b.Result.Description))
	return htmlString(doc)
}

func (b *Book) Cover() string {
	return b.Result.CoverUrl
}

func (b *Book) AddedAt() time.Time {
	return b.Result.DateAdded
}

func (b *Book) PublishedAt() time.Time {
	return b.Result.Issued
}

func (b *Book) URL() string {
	u, _ := url.JoinPath("https://learning.oreilly.com", b.Result.WebUrl)
	return u
}

func (b *Book) Page() int {
	return b.Result.VirtualPages
}

func (b *Book) Authors() []string {
	return b.Result.Authors
}

func (b *Book) Publishers() []string {
	return b.Result.Publishers
}

// htmlタグの削除
func htmlString(node *html.Node) string {
	var buf bytes.Buffer

	var f func(node *html.Node)
	f = func(node *html.Node) {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	return buf.String()
}
