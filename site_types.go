package main

import (
	"reflect"
	"strings"
	"time"

	"github.com/gotidy/ptr"
)

type Site []Ipage

func (self Site) Posts() (out []PagePost) {
	for _, val := range self {
		switch val := val.(type) {
		case PagePost:
			out = append(out, val)
		}
	}
	return
}

func (self Site) ListedPosts() (out []PagePost) {
	for _, val := range self.Posts() {
		if val.IsListed {
			out = append(out, val)
		}
	}
	return
}

func (self Site) PageByType(ref interface{}) Ipage {
	for _, val := range self {
		if reflect.TypeOf(val) == reflect.TypeOf(ref) {
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

func (self Page) Write(body []byte) { writePublic(self.Path, body) }
func (self Page) Make(site Site)    { panic("implement in subclass") }

func (self Page) MdOnce(val interface{}) []byte {
	if self.MdTpl != nil && self.MdHtml == nil {
		self.MdHtml = self.Md(val, nil)
	}
	return self.MdHtml
}

func (self Page) Md(val interface{}, opt *MdOpt) []byte {
	return mdTplToHtml(self.MdTpl, opt, val)
}

type PagePost struct {
	Page
	RedirFrom   []string
	PublishedAt *time.Time
	UpdatedAt   *time.Time
	IsListed    bool
}

func (self PagePost) ExistsAsFile() bool {
	return self.PublishedAt != nil || !FLAGS.PROD
}

func (self PagePost) ExistsInFeeds() bool {
	return self.ExistsAsFile() && bool(self.IsListed)
}

func (self PagePost) UrlFromSiteRoot() string {
	return ensureLeadingSlash(trimExt(self.Path))
}

// Somewhat inefficient but shouldn't be measurable.
func (self PagePost) TimeString() string {
	var out []string

	if self.PublishedAt != nil {
		out = append(out, `published `+timeFmtHuman(*self.PublishedAt))
		if self.UpdatedAt != nil {
			out = append(out, `updated `+timeFmtHuman(*self.UpdatedAt))
		}
	}

	return strings.Join(out, ", ")
}

func (self PagePost) Make(site Site) {
	writePublic(self.Path, self.Render(site))

	for _, path := range self.RedirFrom {
		writePublic(path, Ebui(func(E E) {
			E(`meta`, A{{`http-equiv`, `refresh`}, {`content`, `0;URL='` + self.UrlFromSiteRoot() + `'`}})
		}))
	}
}

func (self PagePost) MakeMd() []byte {
	if self.MdTpl != nil && self.MdHtml == nil {
		self.MdHtml = mdTplToHtml(self.MdTpl, nil, self)
	}
	return self.MdHtml
}

func (self PagePost) FeedItem() FeedItem {
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

type Work struct {
	Name  string
	Link  string
	Start string
	End   string
	Role  string
	Tech  string
	Desc  string
}

func (self Work) Meta() string {
	return strJoin(`; `, self.Role, self.Tech, self.Range())
}

func (self Work) Range() string {
	if self.Start != "" && self.End != "" {
		return self.Start + EMDASH + self.End
	}
	if self.Start != "" && self.End == "" {
		return self.Start + `+`
	}
	return ``
}
