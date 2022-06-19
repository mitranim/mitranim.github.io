package main

import (
	"fmt"
	"time"

	"github.com/mitranim/gg"
)

var (
	FEED_AUTHOR = &FeedAuthor{
		Name:  `Nelo Mitranim`,
		Email: EMAIL,
	}
)

func siteBase() Url {
	if FLAGS.PROD {
		return urlParse(`https://mitranim.com`)
	}
	return urlParse(fmt.Sprintf(`http://localhost:%v`, SERVER_PORT))
}

func siteFeed() Feed {
	base := siteBase()

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
