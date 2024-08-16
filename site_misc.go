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
	Make()
}

func PageByType[A Ipage](src []Ipage) A {
	return gg.Find(src, func(val Ipage) bool {
		return r.TypeOf(val) == gg.Type[A]()
	}).(A)
}

func PageWrite[A Ipage](page A, body []byte) { writePublic(page.GetPath(), body) }
