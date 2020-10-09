package main

import (
	"flag"
	"html/template"
	"log"

	"github.com/pkg/errors"
	"github.com/romanyx/home/internal/cv"
	"github.com/romanyx/home/internal/httprouter"
	"github.com/romanyx/home/internal/medium"
	"github.com/romanyx/home/internal/templates"
)

var (
	addr = flag.String("addr", ":3000", "server bind address")
)

func main() {
	flag.Parse()

	t := template.New("base")
	for _, name := range templates.AssetNames() {
		data, err := templates.Asset(name)
		if err != nil {
			log.Fatalf("asset: %v\n", err)
		}
		tmpl := t.New(name)
		if _, err := tmpl.Parse(string(data)); err != nil {
			log.Fatal(errors.Wrapf(err, "parse template %s", name))
		}
	}
	log.Println("load templates")

	logFunc := func(err error) {
		log.Printf("%+v\n", err)
	}

	cv := cv.MustAsset("roman_cv.pdf")
	log.Println("load cv")
	h := httprouter.NewHandler(medium.Stories, logFunc, t, cv)
	s := httprouter.NewServer(*addr, h, httprouter.GzipOn, httprouter.Letsencrypt)
	defer s.Close()

	log.Println("starting server")
	log.Fatal(s.ListenAndServeLetsencrypt())
}
