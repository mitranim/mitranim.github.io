package main

import (
	r "reflect"

	"github.com/mitranim/gg"
)

type Ipage interface {
	GetPath() string
	GetTitle() string
	GetDescription() string
	GetType() string
	GetImage() string
	GetGlobalClass() string
	GetLink() string
	Make(Site)
}

func PageByType[A Ipage](site Site) A {
	return gg.Find(site.Pages, func(val Ipage) bool {
		return r.TypeOf(val) == gg.Type[A]()
	}).(A)
}

func PageWrite[A Ipage](page A, body []byte) { writePublic(page.GetPath(), body) }
