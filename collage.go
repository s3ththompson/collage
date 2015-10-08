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

type Collage struct {
	Folder string   `json:"folder"`
	Images []string `json:"images"`
}

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
	src, err := os.Stat(folder)
	if err != nil {
		fmt.Fprintln(os.Stderr, "directory does not exist")
		os.Exit(1)
	}
	if !src.IsDir() {
		fmt.Fprintln(os.Stderr, "path is not a directory")
		os.Exit(1)
	}

	num := *n
	if num < 0 {
		fmt.Fprintln(os.Stderr, "n cannot be smaller than 0")
		os.Exit(1)
	}

	images, err := fetchImages(folder, num)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read directory files")
		os.Exit(1)
	}
	if len(images) == 0 {
		fmt.Fprintln(os.Stderr, "no images in folder")
		os.Exit(1)
	}
	collage := Collage{
		filepath.Base(folder),
		images,
	}

	http.HandleFunc("/i/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(folder, r.URL.Path[3:]))
	})
	http.HandleFunc("/data.json", func(w http.ResponseWriter, r *http.Request) {
		dataHandler(w, r, collage)
	})
	http.Handle("/", http.FileServer(rice.MustFindBox("static").HTTPBox()))

	fmt.Printf("Collage of [%d] images from [%s] 🎨\n", num, folder)
	fmt.Println("Serving on http://localhost:2222")
	http.ListenAndServe(":2222", nil)
}
