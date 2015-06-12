package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/GeertJohan/go.rice"
)

var usage = `Usage: collage [options...] [FOLDER]

Options:
	-n 	Number of images to display (randomly) from folder
`

var (
	n = flag.Int("n", 100, "")
)

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()

	folderName := "."
	if flag.NArg() > 0 {
		folderName = flag.Args()[0]
	}

	folder, err := filepath.Abs(folderName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "filepath error")
		os.Exit(1)
	}

	num := *n
	if num < 0 {
		fmt.Fprintln(os.Stderr, "n cannot be smaller than 0")
		os.Exit(1)
	}

	http.HandleFunc("/i/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[3:])
	})
	http.HandleFunc("/data.json", func(w http.ResponseWriter, r *http.Request) {
		dataHandler(w, r, folder, num)
	})
	http.Handle("/", http.FileServer(rice.MustFindBox("static").HTTPBox()))

	http.ListenAndServe(":2222", nil)
}
