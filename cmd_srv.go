package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/mitranim/afr"
	"github.com/mitranim/gg"
	"github.com/mitranim/goh"
	"github.com/mitranim/rout"
	"github.com/mitranim/srv"
)

func init() { commands.Add(`srv`, cmdSrv) }

func cmdSrv() {
	srv := Server{Dir: PUBLIC_DIR}
	go srv.Watch()
	srv.Serve(SERVER_PORT)
}

type Server struct {
	afr.Broad
	Dir string
}

/*
Implementation note: technically it would be preferable to use the FS watching
library "github.com/rjeczalik/notify" which uses various OS-specific APIs that
allow to listen for FS events and avoid polling. However, in this codebase, we
use "github.com/fsnotify/fsnotify" because the "better" watcher mentioned
earlier depends on CGo and slows down compilation, which can be slightly
annoying here.
*/
func (self *Server) Watch() {
	watcher := gg.Try1(fsnotify.NewWatcher())
	defer watcher.Close()

	dir := self.Dir
	walkDirs(dir, func(path string, ent fs.DirEntry) { gg.Try(watcher.Add(path)) })

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			path := event.Name

			if event.Op&fsnotify.Create != 0 && isDir(path) {
				gg.Try(watcher.Add(path))
				continue
			}

			go self.Broad.SendMsg(afrChangeMsg(dir, path))

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf(`[srv] watch error: %+v`, err)
		}
	}
}

func (self *Server) Serve(port int) {
	log.Printf(`[srv] listening on http://localhost:%v`, port)
	gg.Try(http.ListenAndServe(fmt.Sprintf(`:%v`, port), self))
}

func (self *Server) ServeHTTP(rew http.ResponseWriter, req *http.Request) {
	preventCaching(rew.Header())
	rout.MakeRou(rew, req).Serve(self.Route)
}

func (self *Server) Route(rou rout.Rou) {
	rou.Sta(`/afr`).Handler(&self.Broad)
	rou.Get().Func(self.Fallback)
}

func (self *Server) Fallback(rew http.ResponseWriter, req *http.Request) {
	// Allows local import overrides.
	if (goh.File{Path: req.URL.Path}).ServedHTTP(rew, req) {
		return
	}
	srv.FileServer(self.Dir).ServeHTTP(rew, req)
}

func preventCaching(head http.Header) {
	head.Set(`Cache-Control`, `no-store, max-age=0`)
}

func afrChangeMsg(base, path string) afr.Msg {
	return afr.Msg{
		Type: `change`,
		Path: gg.Try1(filepath.Rel(base, path)),
	}
}
