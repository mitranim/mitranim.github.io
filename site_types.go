package main

import (
	"strings"
	"time"

	"github.com/gotidy/ptr"
)

type Site []Ipage

func (self Site) Posts() (out []Post) {
	for _, val := range self {
		switch val := val.(type) {
		case Post:
			out = append(out, val)
		}
	}
	return
}

func (self Site) ListedPosts() (out []Post) {
	for _, val := range self.Posts() {
		if val.IsListed {
			out = append(out, val)
		}
	}
	return
}

type Ipage interface {
	GetPath() string
	GetTitle() string
	GetDescription() string
	GetType() string
	GetImage() string
	GetGlobalClass() string
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
	Fun         func(Site, Page) []byte
	MdToHtml    []byte // Compiled once and reused, if necessary.
}

func (self Page) GetPath() string        { return self.Path }
func (self Page) GetTitle() string       { return self.Title }
func (self Page) GetDescription() string { return self.Description }
func (self Page) GetType() string        { return self.Type }
func (self Page) GetImage() string       { return self.Image }
func (self Page) GetGlobalClass() string { return self.GlobalClass }

func (self Page) Make(site Site) {
	writePublic(self.Path, self.Fun(site, self))
}

func (self Page) MakeMd() []byte {
	if self.MdTpl != nil && self.MdToHtml == nil {
		self.MdToHtml = mdTplToHtml(self.MdTpl, self)
	}
	return self.MdToHtml
}

type Post struct {
	Page
	RedirFrom   []string
	PublishedAt *time.Time
	UpdatedAt   *time.Time
	IsListed    bool
}

func (self Post) ExistsAsFile() bool {
	return self.PublishedAt != nil || !FLAGS.PROD
}

func (self Post) ExistsInFeeds() bool {
	return self.ExistsAsFile() && bool(self.IsListed)
}

func (self Post) UrlFromSiteRoot() string {
	return ensureLeadingSlash(trimExt(self.Path))
}

// Somewhat inefficient but shouldn't be measurable.
func (self Post) TimeString() string {
	var out []string

	if self.PublishedAt != nil {
		out = append(out, `published `+timeFmtHuman(*self.PublishedAt))
		if self.UpdatedAt != nil {
			out = append(out, `updated `+timeFmtHuman(*self.UpdatedAt))
		}
	}

	return strings.Join(out, ", ")
}

func (self Post) Make(site Site) {
	writePublic(self.Path, PagePost(site, self))

	for _, path := range self.RedirFrom {
		writePublic(path, Ebui(func(E E) {
			E(`meta`, A{{`http-equiv`, `refresh`}, {`content`, `0;URL='` + self.UrlFromSiteRoot() + `'`}})
		}))
	}
}

func (self Post) MakeMd() []byte {
	if self.MdTpl != nil && self.MdToHtml == nil {
		self.MdToHtml = mdTplToHtml(self.MdTpl, self)
	}
	return self.MdToHtml
}

func (self Post) FeedItem() FeedItem {
	href := siteBase() + self.UrlFromSiteRoot()

	return FeedItem{
		XmlBase:     href,
		Title:       self.Page.Title,
		Link:        &FeedLink{Href: href},
		Author:      FEED_AUTHOR,
		Description: self.Page.Description,
		Id:          href,
		Published:   self.PublishedAt,
		Updated:     timeCoalesce(self.PublishedAt, self.UpdatedAt, ptr.Time(time.Now().UTC())),
		Content:     Ebui(func(E E) { FeedPostLayout(E, self) }).String(),
	}
}
