package main

import (
	"encoding/xml"
	"time"

	"github.com/mitranim/try"
)

/*
Feed-related types are adapted from https://github.com/gorilla/feeds with minor
modifications, mainly support for `base`.

TODO: remake as a library, with support for streaming XML generation.
*/

type Feed struct {
	Title       string
	XmlBase     string
	AltLink     *FeedLink
	SelfLink    *FeedLink
	Description string
	Author      *FeedAuthor
	Published   *time.Time
	Updated     *time.Time
	Id          string
	Subtitle    string
	Items       []FeedItem
	Copyright   string
	Image       *FeedImage
}

type FeedLink struct {
	Rel    string
	Type   string
	Href   string
	Length string
}

type FeedAuthor struct {
	Name  string
	Email string
}

func (self FeedAuthor) RssAuthor() string {
	if len(self.Email) > 0 {
		str := self.Email
		if len(self.Name) > 0 {
			str += " (" + self.Name + ")"
		}
		return str
	}
	return self.Name
}

type FeedImage struct {
	Url    string
	Title  string
	Link   string
	Width  int64
	Height int64
}

type FeedItem struct {
	XmlBase     string
	Title       string
	Link        *FeedLink
	Source      *FeedLink
	Author      *FeedAuthor
	Description string // used as description in rss, summary in atom
	Id          string // used as guid in rss, id in atom
	Published   *time.Time
	Updated     *time.Time
	Enclosure   *FeedEnclosure
	Content     string
}

type FeedEnclosure struct {
	Url    string
	Length string
	Type   string
}

func (self Feed) AtomFeed() AtomFeed {
	feed := AtomFeed{
		Xmlns:    "https://www.w3.org/2005/Atom",
		XmlBase:  self.XmlBase,
		Title:    self.Title + " | Atom | mitranim",
		Subtitle: self.Description,
		Updated:  (*AtomTime)(self.Updated),
		Rights:   self.Copyright,
	}

	if self.AltLink != nil {
		feed.Id = self.AltLink.Href
		feed.Links = append(feed.Links, AtomLink{
			Rel:  self.AltLink.Rel,
			Type: self.AltLink.Type,
			Href: self.AltLink.Href,
		})
	}

	if self.SelfLink != nil {
		feed.Links = append(feed.Links, AtomLink{
			Rel:  self.SelfLink.Rel,
			Type: self.SelfLink.Type,
			Href: self.SelfLink.Href,
		})
	}

	if self.Author != nil {
		feed.Author = &AtomAuthor{
			AtomPerson: AtomPerson{
				Name:  self.Author.Name,
				Email: self.Author.Email,
			},
		}
	}

	for _, item := range self.Items {
		var name string
		var email string
		if item.Author != nil {
			name = item.Author.Name
			email = item.Author.Email
		}

		entry := AtomEntry{
			XmlBase:   item.XmlBase,
			Title:     item.Title,
			Id:        item.Id,
			Published: (*AtomTime)(item.Published),
			Updated:   (*AtomTime)(timeCoalesce(item.Updated, item.Published, timeNow().MaybeTime())),
			Summary:   &AtomSummary{Type: "html", Content: item.Description},
		}

		var linkRel string
		if item.Link != nil {
			link := AtomLink{
				Href: item.Link.Href,
				Rel:  item.Link.Rel,
				Type: item.Link.Type,
			}
			if link.Rel == `` {
				link.Rel = "alternate"
			}
			linkRel = link.Rel
			entry.Links = append(entry.Links, link)
		}

		if item.Enclosure != nil && linkRel != "enclosure" {
			entry.Links = append(entry.Links, AtomLink{
				Href:   item.Enclosure.Url,
				Rel:    "enclosure",
				Type:   item.Enclosure.Type,
				Length: item.Enclosure.Length,
			})
		}

		// If there's content, assume it's html
		if len(item.Content) > 0 {
			entry.Content = &AtomContent{Type: "html", Content: item.Content}
		}

		if len(name) > 0 || len(email) > 0 {
			entry.Author = &AtomAuthor{AtomPerson: AtomPerson{Name: name, Email: email}}
		}

		feed.Entries = append(feed.Entries, entry)
	}

	return feed
}

func (self Feed) RssFeed() RssFeed {
	var author string
	if self.Author != nil {
		author = self.Author.RssAuthor()
	}

	var image *RssImage
	if self.Image != nil {
		image = &RssImage{
			Url:    self.Image.Url,
			Title:  self.Image.Title,
			Link:   self.Image.Link,
			Width:  self.Image.Width,
			Height: self.Image.Height,
		}
	}

	feed := RssFeed{
		Version:          "2.0",
		ContentNamespace: "http://purl.org/rss/1.0/modules/content/",
		XmlBase:          self.XmlBase,
		Channel: &RssChannel{
			Title:          self.Title + " | RSS | mitranim",
			Description:    self.Description,
			ManagingEditor: author,
			PubDate:        (*RssTime)(timeCoalesce(self.Published, timeNow().MaybeTime())),
			LastBuildDate:  (*RssTime)(timeCoalesce(self.Updated, self.Published, timeNow().MaybeTime())),
			Copyright:      self.Copyright,
			Image:          image,
		},
	}

	if self.AltLink != nil {
		feed.Channel.AltLink = self.AltLink.Href
	}

	for _, item := range self.Items {
		rssItem := RssItem{
			XmlBase:     item.XmlBase,
			Title:       item.Title,
			Description: item.Description,
			Guid:        item.Id,
			PubDate:     (*RssTime)(item.Published),
		}

		if item.Link != nil {
			rssItem.Link = item.Link.Href
		}

		if len(item.Content) > 0 {
			rssItem.Content = &RssContent{Content: item.Content}
		}

		if item.Source != nil {
			rssItem.Source = item.Source.Href
		}

		if item.Enclosure != nil && item.Enclosure.Type != `` && item.Enclosure.Length != `` {
			rssItem.Enclosure = &RssEnclosure{
				Url:    item.Enclosure.Url,
				Type:   item.Enclosure.Type,
				Length: item.Enclosure.Length,
			}
		}

		if item.Author != nil {
			rssItem.Author = item.Author.RssAuthor()
		}

		feed.Channel.Items = append(feed.Channel.Items, rssItem)
	}

	return feed
}

type AtomFeed struct {
	XMLName     xml.Name    `xml:"feed"`
	Xmlns       string      `xml:"xmlns,attr"`
	XmlBase     string      `xml:"xml:base,attr,omitempty"`
	Title       string      `xml:"title"`   // required
	Id          string      `xml:"id"`      // required
	Updated     *AtomTime   `xml:"updated"` // required
	Category    string      `xml:"category,omitempty"`
	Icon        string      `xml:"icon,omitempty"`
	Logo        string      `xml:"logo,omitempty"`
	Rights      string      `xml:"rights,omitempty"` // copyright used
	Subtitle    string      `xml:"subtitle,omitempty"`
	Author      *AtomAuthor `xml:"author,omitempty"`
	Contributor *AtomContributor
	Links       []AtomLink
	Entries     []AtomEntry
}

// Multiple links with different rel can coexist
// Atom 1.0 <link rel="enclosure" type="audio/mpeg" title="MP3" href="https://www.example.org/myaudiofile.mp3" length="1234" />
type AtomLink struct {
	XMLName xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Href    string   `xml:"href,attr"`
	Length  string   `xml:"length,attr,omitempty"`
}

type AtomAuthor struct {
	XMLName xml.Name `xml:"author"`
	AtomPerson
}

type AtomContributor struct {
	XMLName xml.Name `xml:"contributor"`
	AtomPerson
}

type AtomPerson struct {
	Name  string `xml:"name,omitempty"`
	Uri   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`
}

type AtomEntry struct {
	XMLName     xml.Name     `xml:"entry"`
	Xmlns       string       `xml:"xmlns,attr,omitempty"`
	XmlBase     string       `xml:"xml:base,attr,omitempty"`
	Title       string       `xml:"title"` // required
	Id          string       `xml:"id"`    // required
	Category    string       `xml:"category,omitempty"`
	Content     *AtomContent ``
	Rights      string       `xml:"rights,omitempty"`
	Source      string       `xml:"source,omitempty"`
	Published   *AtomTime    `xml:"published,omitempty"`
	Updated     *AtomTime    `xml:"updated"` // required
	Contributor *AtomContributor
	Links       []AtomLink   // required if no content
	Summary     *AtomSummary // required if content has src or is base64
	Author      *AtomAuthor  // required if feed lacks an author
}

type AtomContent struct {
	XMLName xml.Name `xml:"content"`
	Content string   `xml:",cdata"`
	Type    string   `xml:"type,attr"`
}

type AtomSummary struct {
	XMLName xml.Name `xml:"summary"`
	Content string   `xml:",cdata"`
	Type    string   `xml:"type,attr"`
}

type AtomTime time.Time

func (self AtomTime) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	try.To(enc.EncodeToken(start))
	try.To(enc.EncodeToken(xml.CharData(time.Time(self).Format(time.RFC3339))))
	try.To(enc.EncodeToken(xml.EndElement{Name: start.Name}))
	return nil
}

type RssFeed struct {
	XMLName          xml.Name `xml:"rss"`
	XmlBase          string   `xml:"xml:base,attr,omitempty"`
	Version          string   `xml:"version,attr"`
	ContentNamespace string   `xml:"xmlns:content,attr"`
	Channel          *RssChannel
}

type RssChannel struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`       // required
	AltLink        string   `xml:"link"`        // required
	Description    string   `xml:"description"` // required
	Language       string   `xml:"language,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"` // Author used
	WebMaster      string   `xml:"webMaster,omitempty"`
	PubDate        *RssTime `xml:"pubDate,omitempty"`       // created or updated
	LastBuildDate  *RssTime `xml:"lastBuildDate,omitempty"` // updated used
	Category       string   `xml:"category,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	Ttl            int64    `xml:"ttl,omitempty"`
	Rating         string   `xml:"rating,omitempty"`
	SkipHours      string   `xml:"skipHours,omitempty"`
	SkipDays       string   `xml:"skipDays,omitempty"`
	Image          *RssImage
	TextInput      *RssTextInput
	Items          []RssItem
}

type RssImage struct {
	XMLName xml.Name `xml:"image"`
	Url     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Width   int64    `xml:"width,omitempty"`
	Height  int64    `xml:"height,omitempty"`
}

type RssTextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}

type RssItem struct {
	XMLName     xml.Name      `xml:"item"`
	XmlBase     string        `xml:"xml:base,attr,omitempty"`
	Title       string        `xml:"title"`       // required
	Link        string        `xml:"link"`        // required
	Description string        `xml:"description"` // required
	Content     *RssContent   ``
	Author      string        `xml:"author,omitempty"`
	Category    string        `xml:"category,omitempty"`
	Comments    string        `xml:"comments,omitempty"`
	Enclosure   *RssEnclosure ``
	Guid        string        `xml:"guid,omitempty"`    // Id used
	PubDate     *RssTime      `xml:"pubDate,omitempty"` // created or updated
	Source      string        `xml:"source,omitempty"`
}

type RssContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Content string   `xml:",cdata"`
}

// RSS 2.0 <enclosure url="https://example.com/file.mp3" length="123456789" type="audio/mpeg" />
type RssEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

type RssTime time.Time

func (self RssTime) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	try.To(enc.EncodeToken(start))
	try.To(enc.EncodeToken(xml.CharData(time.Time(self).Format(time.RFC1123Z))))
	try.To(enc.EncodeToken(xml.EndElement{Name: start.Name}))
	return nil
}
