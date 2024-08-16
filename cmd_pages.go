package main

import "github.com/mitranim/gg"

func init() { commands.Add(`pages`, cmdPages) }

func cmdPages() {
	defer gg.LogTimeNow(`pages`).LogStart().LogEnd()

	var site Site
	site.Init()

	makePages(site)
	makeFeeds(site)
}

func makePages(site Site) {
	for _, val := range site.All() {
		val.Make()
	}
}

func makeFeeds(site Site) {
	feed := siteFeed()

	for _, post := range site.Posts {
		if post.ExistsInFeeds() {
			feed.Items = append(feed.Items, post.FeedItem())
		}
	}

	writePublic(`feed.xml`, xmlEncode(feed.AtomFeed()))
	writePublic(`feed_rss.xml`, xmlEncode(feed.RssFeed()))
}
