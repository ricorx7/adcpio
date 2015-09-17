package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func debugHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello adcp.IO!") // send data to client side
	for key, value := range server.adcp {
		fmt.Fprint(w, key)
		ensJSON, _ := json.Marshal(value.lastEns)
		fmt.Fprintf(w, string(ensJSON))
	}
}
