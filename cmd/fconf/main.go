package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DeviceConfig interface {
	Name() string
	LoadJSON(io.Reader) error
	Save() error
	ToJSON(io.Writer) error
}

func main() {
	port := flag.Int64("port", 8970, "Specify the port in which to run the magager")
	flag.Parse()
	s := http.NewServeMux()
	w := newWeb()
	s.HandleFunc("/", w.Home)
	s.HandleFunc("/device/dongle", w.Dongle)
	log.Printf(" starting server on  localhost:%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), s))

}

type web struct {
}

func newWeb() *web {
	return &web{}
}

func (ww *web) Home(w http.ResponseWriter, r *http.Request) {
}

func (ww *web) Dongle(w http.ResponseWriter, r *http.Request) {
}
