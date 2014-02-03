package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	//err := r.ParseForm()
	err := r.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile

	if err != nil {
		log.Printf("handleFileUpload() could not ParseMultipartForm on request '%#+v': returned err: %v\n", r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Printf("handleFileUpload request r is: %#+v\n", r)
	/*
		fmt.Printf("handleFileUpload r.PostForm is: %#+v\n", r.PostForm)
		if r.MultipartForm != nil {
			fmt.Printf("r.MultipartForm is: %#+v\n", r.MultipartForm)
		} else {
			fmt.Printf("r.MultipartForm is nil\n")
		}
	*/

	/*
		    mpFile, mpHeader, err := r.FormFile("TheFile")
			if err != nil {
				log.Printf("Can't Find 'TheFile' in form ")
				return
			} else {
				log.Printf("successfully ran r.FormFile, getting mpHeader.Filename = %#+v\n", mpHeader.Filename)
			}
	*/

	fhs := r.MultipartForm.File["TheFile"]
	//log.Printf("r.MultipartForm.File[`TheFile`] == fhs: %#+v\n", fhs)

	/*
		if fhs != nil && len(fhs) > 0 && fhs[0] != nil {
			//log.Printf("r.MultipartForm.File[`TheFile`][0]: %#+v\n", fhs[0])
			log.Printf("  with .Filename: %#+v\n", fhs[0].Filename)
		}
	*/
	numFiles := len(fhs)
	receivedFiles := make(map[string]string)
	log.Printf("handleFileUpload(): received %d files associated with 'TheFile'\n", numFiles)

	for _, fh := range fhs {

		f, err := fh.Open()
		// f is one of the files
		defer f.Close()

		localFn := "uploads/" + fh.Filename
		os.Mkdir("uploads", 0744)
		t, err := os.Create(localFn)
		if err != nil {
			log.Printf("handleFileUpload, error trying to create file '%s': %s\n", localFn, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer t.Close()
		log.Printf("  received file named: browser:'%s' ---> local:'%s'\n", fh.Filename, localFn)

		_, err = io.Copy(t, f)

		if err != nil {
			panic(err)
		}

		//log.Printf("handleFileUpload(): uploaded one of the 'TheFile' form element to file: '%s;\n", t.Name())

		receivedFiles[fh.Filename] = localFn
		t.Close()
		f.Close()
	}

	fmt.Fprintf(w, "<html><body>Server received these files: %#+v\n</body></html>", receivedFiles)
}
