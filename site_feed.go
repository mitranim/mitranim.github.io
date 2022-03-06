package main

import (
	"fmt"
	"time"

	"github.com/gotidy/ptr"
)

var (
	FEED_AUTHOR = &FeedAuthor{
		Name:  `Nelo Mitranim`,
		Email: EMAIL,
	}
)

func siteBase() string {
	if FLAGS.PROD {
		return `https://mitranim.com`
	}
	return fmt.Sprintf(`http://localhost:%v`, SERVER_PORT)
}

func siteFeed() Feed {
	base := siteBase()

	return Feed{
		Title:   `Software, Tech, Philosophy, Games`,
		XmlBase: base,
		AltLink: &FeedLink{
			Rel:  `alternate`,
			Type: `text/html`,
			Href: base + `/posts`,
		},
		SelfLink: &FeedLink{
			Rel:  `self`,
			Type: `application/atom+xml`,
			Href: base + `/feed.xml`,
		},
		Author:      FEED_AUTHOR,
		Published:   ptr.Time(time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)),
		Updated:     ptr.Time(time.Now()),
		Id:          base + `/posts`,
		Description: `Random thoughts about technology`,
		Items:       nil,
	}
}
