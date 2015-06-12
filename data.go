package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
)

func dataHandler(w http.ResponseWriter, r *http.Request, folder string, num int) {
	images, err := fetchImages(folder, num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	options := struct {
		Folder string   `json:"folder"`
		Images []string `json:"images"`
	}{
		filepath.Base(folder),
		images,
	}
	js, err := json.Marshal(options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func fetchImages(folder string, num int) ([]string, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	out := []string{}
	for _, fp := range files {
		if fp.IsDir() || !isImage(fp.Name()) {
			continue
		}
		out = append(out, fp.Name())
	}
	Shuffle(out)
	if num != 0 {
		return out[:num], nil
	} else {
		return out, nil
	}
}

func Shuffle(a []string) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func isImage(name string) bool {
	ext := filepath.Ext(name)
	return (ext == ".gif" || ext == ".jpg" || ext == ".jpeg" || ext == ".png")
}
