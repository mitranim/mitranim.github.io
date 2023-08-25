package main

import (
	"time"

	"github.com/mitranim/gg"
)

var (
	FEED_AUTHOR = &FeedAuthor{
		Name:  `Nelo Mitranim`,
		Email: EMAIL,
	}
)

var siteBaseUrl = gg.NewLazy(func() (out Url) {
	if FLAGS.PROD {
		out.Scheme = `https`
		out.Host = `mitranim.com`
	} else {
		out.Scheme = `http`
		out.Host = gg.Str(`localhost:`, SERVER_PORT)
	}
	return
})

func siteFeed() Feed {
	base := siteBaseUrl.Get()

	return Feed{
		Title:   `Software, Tech, Philosophy, Games`,
		XmlBase: base.String(),
		AltLink: &FeedLink{
			Rel:  `alternate`,
			Type: `text/html`,
			Href: base.WithPath(`/posts`).String(),
		},
		SelfLink: &FeedLink{
			Rel:  `self`,
			Type: `application/atom+xml`,
			Href: base.WithPath(`/feed.xml`).String(),
		},
		Author:      FEED_AUTHOR,
		Published:   gg.Ptr(time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)),
		Updated:     gg.Ptr(time.Now()),
		Id:          base.WithPath(`/posts`).String(),
		Description: `Random thoughts about technology`,
		Items:       nil,
	}
}
