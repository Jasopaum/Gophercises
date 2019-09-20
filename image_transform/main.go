package main

import (
	"fmt"
	"gophercises/image_transform/primitive"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

type transformConfig struct {
	mode      primitive.Mode
	numShapes int
}

// Transformations with different shapes
var choiceModes = [4]primitive.Mode{
	primitive.ModeTriangle,
	primitive.ModePolygon,
	primitive.ModeRotatedRect,
	primitive.ModeCombo,
}
var choiceNums = [4]int{10, 20, 30, 40}

func configShapes(numShapes int) []transformConfig {
	ret := make([]transformConfig, len(choiceModes))
	for i, m := range choiceModes {
		ret[i] = transformConfig{m, numShapes}
	}
	return ret
}

func configNum(mode primitive.Mode) []transformConfig {
	ret := make([]transformConfig, len(choiceNums))
	for i, n := range choiceNums {
		ret[i] = transformConfig{mode, n}
	}
	return ret
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		html := `
<html><body>
<form action="/upload" method="post" enctype="multipart/form-data">
	<input type="file" name="image">
	<button type="submit">Upload Image</button>
</form>
</body></html>
`
		fmt.Fprint(w, html)
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, req *http.Request) {
		// Get file from request
		file, header, err := req.FormFile("image")
		ext := filepath.Ext(header.Filename)[1:]
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save image file on server
		origImgFile, err := ioutil.TempFile("./img/", fmt.Sprintf("orig-*.%s", ext))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer origImgFile.Close()
		_, err = io.Copy(origImgFile, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, fmt.Sprintf("/chooseshape/%s", filepath.Base(origImgFile.Name())), http.StatusFound)
	})
	mux.HandleFunc("/chooseshape/", func(w http.ResponseWriter, req *http.Request) {
		// Get file
		file, err := os.Open("./img/" + filepath.Base(req.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate trasformed images
		ext := filepath.Ext(file.Name())[1:]
		imgFiles, err := genImages(file, ext, configShapes(10))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type datastruct struct {
			Name string
			Mode primitive.Mode
		}
		data := make([]datastruct, len(imgFiles))
		for ind, img := range imgFiles {
			data[ind] = datastruct{img, choiceModes[ind]}
		}
		html := `<html><body>
			{{range .}}
				<a href="/choosenum/{{.Name}}?mode={{.Mode}}">
					<img style="width: 20%;" src="/img/{{.Name}}">
				</a>
			{{end}}
			</body></html>`
		tpl := template.Must(template.New("").Parse(html))
		err = tpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/choosenum/", func(w http.ResponseWriter, req *http.Request) {
		// Get file
		file, err := os.Open("./img/" + filepath.Base(req.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mode, err := strconv.Atoi(req.FormValue("mode"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate trasformed images
		ext := filepath.Ext(file.Name())[1:]
		imgFiles, err := genImages(file, ext, configNum(primitive.Mode(mode)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		html := `<html><body>
			{{range .}}
				<a href="/img/{{.}}">
					<img style="width: 20%;" src="/img/{{.}}">
				</a>
			{{end}}
			</body></html>`
		tpl := template.Must(template.New("").Parse(html))
		err = tpl.Execute(w, imgFiles)
		if err != nil {
			panic(err)
		}
	})

	mux.Handle("/img/", http.FileServer(http.Dir("./")))

	http.ListenAndServe(":8000", mux)
}

func genImages(file io.ReadSeeker, ext string, configs []transformConfig) ([]string, error) {
	var outFiles []string
	for _, conf := range configs {
		f, err := genImage(file, ext, conf)
		if err != nil {
			return nil, err
		}
		outFiles = append(outFiles, f)
	}
	return outFiles, nil
}

func genImage(file io.ReadSeeker, ext string, config transformConfig) (string, error) {
	// Create image file on server
	imgFile, err := ioutil.TempFile("./img/", fmt.Sprintf("img-*.%s", ext))
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	// Create primitive image
	file.Seek(0, 0)
	outPrimitive, err := primitive.Transform(file, config.numShapes, primitive.WithMode(config.mode))

	// Save primitive image in server file
	_, err = io.Copy(imgFile, outPrimitive)
	if err != nil {
		return "", err
	}
	return filepath.Base(imgFile.Name()), nil
}
