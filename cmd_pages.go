package main

import (
	"github.com/mitranim/try"
)

func cmdPages() (err error) {
	defer try.Rec(&err)
	defer timing("pages")()

	site := initSite()
	try.To(makePages(site))
	try.To(makeFeeds(site))
	return
}

func makePages(site Site) (err error) {
	defer try.Rec(&err)
	for _, val := range site.Ipages {
		try.To(val.Make(site))
	}
	return
}

func makeFeeds(site Site) (err error) {
	defer try.Rec(&err)

	feed := siteFeed()

	for _, post := range site.Posts() {
		if post.ExistsInFeeds() {
			feed.Items = append(feed.Items, post.FeedItem())
		}
	}

	try.To(writePublic("feed.xml", try.ByteSlice(xmlEncode(feed.AtomFeed()))))
	try.To(writePublic("feed_rss.xml", try.ByteSlice(xmlEncode(feed.RssFeed()))))
	return
}
