package main

import (
	"github.com/mitranim/gg"
)

type Site struct {
	Pages []Ipage
	Posts []PagePost
	Games GameColl
}

func (self *Site) Init() {
	self.Pages = initSitePages(self)
	self.Posts = initSitePosts(self)
	self.Games = initSiteGames()
}

func (self Site) All() (out []Ipage) {
	out = append(out, self.Pages...)
	for _, val := range self.Posts {
		out = append(out, val)
	}
	return
}

func (self Site) ListedPosts() []PagePost {
	return gg.Filter(self.Posts, PagePost.GetIsListed)
}
