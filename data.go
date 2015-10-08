package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
)

func dataHandler(w http.ResponseWriter, r *http.Request, collage Collage) {
	js, err := json.Marshal(collage)
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
	if num != 0 && num < len(out){
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
	ext := strings.ToLower(filepath.Ext(name))
	return (ext == ".gif" || ext == ".jpg" || ext == ".jpeg" || ext == ".png")
}
