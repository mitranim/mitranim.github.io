package main

func init() { commands.Add(`pages`, cmdPages) }

func cmdPages() {
	defer timing("pages")()

	site := initSite()
	makePages(site)
	makeFeeds(site)
}

func makePages(site Site) {
	for _, val := range site.All() {
		val.Make(site)
	}
}

func makeFeeds(site Site) {
	feed := siteFeed()

	for _, post := range site.Posts {
		if post.ExistsInFeeds() {
			feed.Items = append(feed.Items, post.FeedItem())
		}
	}

	writePublic("feed.xml", xmlEncode(feed.AtomFeed()))
	writePublic("feed_rss.xml", xmlEncode(feed.RssFeed()))
}
