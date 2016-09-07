package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var rootDir = flag.String("d", "./", "document root directory")
var port = flag.String("p", "8000", "listen port")
var ch chan bool

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of: %s\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("statis server:")
	go func() {
		ch_signal := make(chan os.Signal)
		signal.Notify(ch_signal, os.Interrupt)
		<-ch_signal
		log.Println("Exit.")
		os.Exit(0)
	}()
	go httpServer()
	log.Printf("Listening on port: %s..., DocumentRoot: %s\n", *port, *rootDir)
	if <-ch == false {
		log.Printf("Listening on port %s failed\n", *port)
	}
}

func httpServer() {
	h := http.FileServer(http.Dir(*rootDir))
	if err := http.ListenAndServe(":"+*port, requestLog(h)); err != nil {
		log.Fatalln("ListenAndServe error:", err.Error())
		ch <- false
	}
}

func requestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s -- %s %s %s\n", time.Now(), r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
