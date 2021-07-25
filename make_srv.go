package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mitranim/srv"
	"github.com/mitranim/try"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
)

func cmdSrv() error {
	events := make(chan notify.EventInfo, 1)
	defer notify.Stop(events)
	go srvWatch(events)
	return srvServe()
}

func srvWatch(events chan notify.EventInfo) {
	dir := filepath.Join(try.String(os.Getwd()), PUBLIC_DIR)
	try.To(os.MkdirAll(dir, os.ModePerm))
	try.To(notify.Watch(fpj(dir, `...`), events, notify.All))
	for event := range events {
		go afrMaybeSend(try.String(filepath.Rel(dir, event.Path())))
	}
}

func srvServe() error {
	const port = SERVER_PORT
	fmt.Printf("[srv] listening on http://localhost:%v\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), srv.FileServer(PUBLIC_DIR))
	return errors.WithStack(err)
}

func afrMaybeSend(path string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		`http://localhost:52692/afr/send`,
		bytes.NewReader(try.ByteSlice(json.Marshal(AfrWatchMsg{
			Type: `change`,
			Path: path,
		}))),
	)

	res, err := http.DefaultClient.Do(req)
	if res != nil && res.Body != nil {
		res.Body.Close()
	}

	if err != nil {
		// log.Println(`[srv] failed to send afr msg:`, err)
	}
}

type AfrWatchMsg struct {
	Type string `json:"type"`
	Path string `json:"path"`
}
