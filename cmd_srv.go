package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mitranim/afr"
	"github.com/mitranim/rout"
	"github.com/mitranim/srv"
	"github.com/mitranim/try"
	"github.com/rjeczalik/notify"
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
	events := make(chan notify.EventInfo, 1)
	defer notify.Stop(events)

	dir := filepath.Join(try.String(os.Getwd()), self.Dir)

	try.To(os.MkdirAll(dir, os.ModePerm))
	try.To(notify.Watch(fpj(dir, `...`), events, notify.All))

	for event := range events {
		go self.Broad.SendMsg(afrChangeMsg(dir, event.Path()))
	}
}

func (self *Server) Serve(port int) {
	log.Printf("[srv] listening on http://localhost:%v", port)
	try.To(http.ListenAndServe(fmt.Sprintf(":%v", port), self))
}

func (self *Server) ServeHTTP(rew http.ResponseWriter, req *http.Request) {
	preventCaching(rew.Header())
	rout.MakeRouter(rew, req).Serve(self.Route)
}

func (self *Server) Route(r rout.R) {
	r.Reg(`^/afr/`).Handler(&self.Broad)
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
