package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

const downloadDir1 = "//home//ubuntu//upload//"

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
		log.Printf("Received files upload1")
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Printf("Error Uploading File")
			fmt.Println(err)
			return
		}
		log.Printf("File : %s", file)
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(downloadDir1+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("Error Uploading File")
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
