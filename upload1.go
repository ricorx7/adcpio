package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

const downloadDir1 = "//home//ubuntu//upload//"

const maxMemory = 1 * 1024 * 1024

// upload logic
func upload1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upload 1 method:", r.Method)
	if r.Method == "GET" {
		fmt.Println("upload 1 GET:")
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload1.html")
		t.Execute(w, token)
	} else {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
		}

		for key, value := range r.MultipartForm.Value {
			fmt.Fprintf(w, "%s:%s ", key, value)
			log.Printf("%s:%s", key, value)
		}

		for _, fileHeaders := range r.MultipartForm.File {
			for _, fileHeader := range fileHeaders {
				file, _ := fileHeader.Open()
				path := fmt.Sprintf("%s/%s", downloadDir1, fileHeader.Filename)
				buf, _ := ioutil.ReadAll(file)
				ioutil.WriteFile(path, buf, os.ModePerm)
			}
		}
	}

}
