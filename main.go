package main

import (
	"flag"
	"go/build"
	"net/http"
    "fmt"
    "log"
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	assets    = flag.String("assets", defaultAssetPath(), "path to assets")
)

func defaultAssetPath() string {
	p, err := build.Default.Import("github.com/gary.burd.info/go-websocket-chat", "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

func homeHandler(c http.ResponseWriter, req *http.Request) {
    fmt.Fprint(c, "Hello, world!")
	//homeTempl.Execute(c, req.Host)
}

var wsh *wsHandler
func main() {
	//flag.Parse()
	//homeTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "home.html")))
	wsh := NewwsHandler()
	_ = wsh
	http.HandleFunc("/ws", MyHandler)
    http.HandleFunc("/", homeHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
