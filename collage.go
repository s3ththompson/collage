package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var usage = `Usage: collage [options...] [FOLDER]
Options:
	-n 	Number of images to display (randomly) from folder (default 100)
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
		return
	}
	src, err := os.Stat(folder)
	if err != nil {
		fmt.Fprintln(os.Stderr, "directory does not exist")
		return
	}
	if !src.IsDir() {
		fmt.Fprintln(os.Stderr, "path is not a directory")
		return
	}

	num := *n
	if num < 0 {
		fmt.Fprintln(os.Stderr, "n cannot be smaller than 0")
		return
	}

	images, err := fetchImages(folder, num)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read directory files")
		return
	}
	if len(images) == 0 {
		fmt.Fprintln(os.Stderr, "no images in folder")
		return
	}

	collage := struct {
		Folder string   `json:"folder"`
		Images []string `json:"images"`
	}{
		filepath.Base(folder),
		images,
	}

	http.HandleFunc("/i/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(folder, r.URL.Path[3:]))
	})
	http.HandleFunc("/data.json", func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(collage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page)
	})

	fmt.Printf("Collage of [%d] images from [%s] 🎨\n", len(collage.Images), folder)
	fmt.Println("Serving on http://localhost:2222")
	http.ListenAndServe(":2222", nil)
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

	// shuffle array
	for i := range out {
		j := rand.Intn(i + 1)
		out[i], out[j] = out[j], out[i]
	}

	if num != 0 && num < len(out){
		return out[:num], nil
	} else {
		return out, nil
	}
}

func isImage(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return (ext == ".gif" || ext == ".jpg" || ext == ".jpeg" || ext == ".png")
}

var page = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Collage</title>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
		<script type="text/javascript">
			$(function(){
				$.getJSON( "data.json", function(data) {
					document.title = "Collage • " + data.folder;

					for (image of data.images) {
						$("#images").append("<img src='/i/"+image+"'></img>");
					}

					$images = $("#images img");
					$section = $("#images");
					$("#images img").click(function() {
						if ($section.hasClass("dimmed")) {
							$section.removeClass("dimmed");
							$images.removeClass("dim");
						} else {
							$section.addClass("dimmed");
							$images.not($(this)).addClass("dim");
						}
					});
				});
			});
		</script>
		<style type="text/css">
			#images {
				line-height: 0;
				-webkit-column-count: 3;
				-webkit-column-gap:   10px;
				-moz-column-count:    3;
				-moz-column-gap:      10px;
				column-count:         3;
				column-gap:           10px;  
			}

			#images img {
				max-width: 100%;
				height: auto;
				margin-bottom: 10px;
			}

			#images img.dim {
				opacity: 0.1;
			}
		</style>
	</head>
	<body>
		<section id="images">
		</section>
	</body>
</html>`