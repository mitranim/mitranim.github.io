package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/mitranim/afr"
	"github.com/mitranim/rout"
	"github.com/mitranim/srv"
	"github.com/mitranim/try"
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

func (self *Server) Watch() {
	watcher, err := fsnotify.NewWatcher()
	try.To(err)
	defer watcher.Close()

	dir := self.Dir
	walkDirs(dir, func(path string, ent fs.DirEntry) {
		try.To(watcher.Add(path))
	})

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			path := event.Name

			if event.Op&fsnotify.Create != 0 && isDir(path) {
				try.To(watcher.Add(path))
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
	try.To(http.ListenAndServe(fmt.Sprintf(`:%v`, port), self))
}

func (self *Server) ServeHTTP(rew http.ResponseWriter, req *http.Request) {
	preventCaching(rew.Header())
	rout.MakeRouter(rew, req).Serve(self.Route)
}

func (self *Server) Route(r rout.R) {
	r.Begin(`/afr`).Handler(&self.Broad)
	r.Get().Handler(srv.FileServer(self.Dir))
}

func preventCaching(head http.Header) {
	head.Set(`Cache-Control`, `no-store, max-age=0`)
}

func afrChangeMsg(base, path string) afr.Msg {
	return afr.Msg{
		Type: `change`,
		Path: try.String(filepath.Rel(base, path)),
	}
}
