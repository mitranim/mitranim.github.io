package main

func cmdPages() {
	defer timing("pages")()

	site := initSite()
	makePages(site)
	makeFeeds(site)
}

func makePages(site Site) {
	for _, val := range site {
		val.Make(site)
	}
}

func makeFeeds(site Site) {
	feed := siteFeed()

	for _, post := range site.Posts() {
		if post.ExistsInFeeds() {
			feed.Items = append(feed.Items, post.FeedItem())
		}
	}

	writePublic("feed.xml", xmlEncode(feed.AtomFeed()))
	writePublic("feed_rss.xml", xmlEncode(feed.RssFeed()))
}
