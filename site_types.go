package main

import (
	r "reflect"
	"strings"

	e "github.com/pkg/errors"
)

type Site struct {
	Pages []Ipage
	Posts []PagePost
}

func (self Site) All() (out []Ipage) {
	out = append(out, self.Pages...)
	for _, val := range self.Posts {
		out = append(out, val)
	}
	return
}

func (self Site) ListedPosts() (out []PagePost) {
	for _, val := range self.Posts {
		if val.IsListed {
			out = append(out, val)
		}
	}
	return
}

func (self Site) PageByType(ref interface{}) Ipage {
	for _, val := range self.Pages {
		if r.TypeOf(val) == r.TypeOf(ref) {
			return val
		}
	}
	return nil
}

type Ipage interface {
	GetPath() string
	GetTitle() string
	GetDescription() string
	GetType() string
	GetImage() string
	GetGlobalClass() string
	GetLink() string
	Make(Site)
}

type Page struct {
	Path        string
	Title       string
	Description string
	MdTpl       []byte
	Type        string
	Image       string
	GlobalClass string
	MdHtml      []byte // Compiled once and reused, if necessary.
}

func (self Page) GetPath() string        { return self.Path }
func (self Page) GetTitle() string       { return self.Title }
func (self Page) GetDescription() string { return self.Description }
func (self Page) GetType() string        { return self.Type }
func (self Page) GetImage() string       { return self.Image }
func (self Page) GetGlobalClass() string { return self.GlobalClass }

func (self Page) Make(site Site) {
	panic(e.Errorf(`"Make" is not implemented for page %#v`, self))
}

func (self Page) MdOnce(val interface{}) []byte {
	if self.MdTpl != nil && self.MdHtml == nil {
		self.MdHtml = self.Md(val, nil)
	}
	return self.MdHtml
}

func (self Page) Md(val interface{}, opt *MdOpt) []byte {
	return mdTplToHtml(self.MdTpl, opt, val)
}

func (self Page) GetLink() string {
	return ensureLeadingSlash(trimExt(self.GetPath()))
}

func pageWrite(page Ipage, body []byte) { writePublic(page.GetPath(), body) }

// TODO: add `.WrittenAt`, which often doesn't match `.PublishedAt`.
type PagePost struct {
	Page
	RedirFrom   []string
	PublishedAt Time
	UpdatedAt   Time
	IsListed    bool
}

func (self PagePost) ExistsAsFile() bool {
	return !self.PublishedAt.IsZero() || !FLAGS.PROD
}

func (self PagePost) ExistsInFeeds() bool {
	return self.ExistsAsFile() && bool(self.IsListed)
}

// Somewhat inefficient but shouldn't be measurable.
func (self PagePost) TimeString() string {
	var out []string

	if !self.PublishedAt.IsZero() {
		out = append(out, `published `+timeFmtHuman(self.PublishedAt))
		if !self.UpdatedAt.IsZero() {
			out = append(out, `updated `+timeFmtHuman(self.UpdatedAt))
		}
	}

	return strings.Join(out, `, `)
}

func (self PagePost) Make(site Site) {
	pageWrite(self, self.Render(site))

	for _, path := range self.RedirFrom {
		writePublic(path, F(
			E(`meta`, AP(`http-equiv`, `refresh`, `content`, `0;URL='`+self.GetLink()+`'`)),
		))
	}
}

func (self PagePost) MakeMd() []byte {
	if self.MdTpl != nil && self.MdHtml == nil {
		self.MdHtml = mdTplToHtml(self.MdTpl, nil, self)
	}
	return self.MdHtml
}

func (self PagePost) FeedItem() FeedItem {
	// Caution: `path.Join` breaks "//".
	href := siteBase() + self.GetLink()

	return FeedItem{
		XmlBase:     href,
		Title:       self.Page.Title,
		Link:        &FeedLink{Href: href},
		Author:      FEED_AUTHOR,
		Description: self.Page.Description,
		Id:          href,
		Published:   self.PublishedAt.MaybeTime(),
		Updated:     timeCoalesce(self.PublishedAt.MaybeTime(), self.UpdatedAt.MaybeTime(), timeNow().MaybeTime()),
		Content:     FeedPost(self).String(),
	}
}

type Work struct {
	Name string
	Link string
	Role string
	Tech string
	Desc string
	Lifecycle
}

func (self Work) Meta(b B) {
	SemiList{
		F(self.Role),
		F(self.Tech),
		F(self.Range()),
		F(self.StatusLink),
	}.Render(b)
}

type Lifecycle struct {
	Start  string
	End    string
	Status string
	Link   string
}

func (self Lifecycle) Range() string {
	if self.Start != `` && self.End != `` {
		return self.Start + EMDASH + self.End
	}
	if self.Start != `` && self.End == `` {
		return self.Start + `+`
	}
	return ``
}

func (self Lifecycle) StatusLink(b B) {
	if self.Link != `` {
		b.Child(Exta(self.Link, self.Status))
		return
	}

	if self.Status != `` {
		b.Child(self.Status)
		return
	}
}

func (self Lifecycle) StatusEnd(b B) {
	if !buiChild(b, self.StatusLink) {
		b.Child(self.End)
	}
}
