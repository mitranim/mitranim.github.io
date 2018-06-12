package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// Serves static files. Uses more or less the same HTML/URL resolution algorithm
// as GitHub Pages and Netlify.

const (
	PORT      = "11204"
	DIR       = "public"
	NOT_FOUND = DIR + "/404.html"
)

func main() {
	fmt.Printf("Starting server on %v\n", magenta(fmt.Sprintf("http://localhost:%v", PORT)))
	http.ListenAndServe(fmt.Sprintf(":%v", PORT), http.HandlerFunc(handle))
}

func handle(rew http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
	default:
		http.Error(rew, "", http.StatusMethodNotAllowed)
		return
	}

	rpath := req.URL.Path
	fpath := filepath.Join(DIR, rpath)

	// Has extension -> don't apply special rules
	if path.Ext(rpath) != "" {
		http.ServeFile(rew, req, fpath)
		return
	}

	stat, _ := os.Stat(fpath)

	if fileExists(stat) {
		http.ServeFile(rew, req, fpath)
		return
	}

	// Prioritize a possible path.html over a possible path/index.html
	candidatePath := fpath + ".html"
	candidateStat, _ := os.Stat(candidatePath)
	if fileExists(candidateStat) {
		http.ServeFile(rew, req, candidatePath)
		return
	}

	// Directory must exist, otherwise there's no index.html
	if stat == nil {
		http.NotFound(rew, req)
		return
	}

	fpath = filepath.Join(fpath, "index.html")
	stat, _ = os.Stat(fpath)
	if fileExists(stat) {
		http.ServeFile(rew, req, fpath)
		return
	}
	http.NotFound(rew, req)
}

func fileExists(stat os.FileInfo) bool {
	return stat != nil && !stat.IsDir()
}

func magenta(str string) string {
	const FG_MAGENTA = "\x1b[35m"
	const RESET = "\x1b[0m"
	return FG_MAGENTA + str + RESET
}
