package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	"github.com/romanyx/home/internal/httprouter"
	"github.com/romanyx/home/internal/medium"
)

var (
	addr = flag.String("addr", "127.0.0.1:3000", "Server bind address")
)

func main() {
	flag.Parse()

	paths := []string{
		"../../templates/layout.html",
		"../../templates/index.html",
		"../../templates/meta.html",
		"../../templates/og.html",
	}

	cv, err := ioutil.ReadFile("../../static/roman_cv.pdf")
	if err != nil {
		log.Fatal(errors.Wrap(err, "read template"))
	}

	t, err := template.ParseFiles(paths...)
	if err != nil {
		log.Fatal(errors.Wrap(err, "template parse files"))
	}

	logFunc := func(err error) {
		log.Println(err)
	}

	h := httprouter.NewHandler(medium.Stories, logFunc, t, cv)
	s := httprouter.NewServer(*addr, h, httprouter.GzipOn, httprouter.Letsencrypt)
	defer s.Close()

	log.Fatal(s.ListenAndServeLetsencrypt())
}
