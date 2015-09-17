package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// Flags to set at startup
var (
	version      = "0.1"
	versionFloat = float32(0.1)
	addr         = flag.String("addr", ":8080", "http service address")
)

// main will start the application.
func main() {
	// Parse the flags
	flag.Parse()

	// setup logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Run the server
	go server.run()

	// HTTP server
	http.HandleFunc("/", debugHandler)    // Debugger
	http.HandleFunc("/adcp", adcpHandler) // Debugger
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/ws", wsHandler)                // wsHandler in websocketConn.go.  Creates websocket
	http.HandleFunc("/wsAdcp", wsAdcpDisplayHandler) // wsHandler in websocketConn.go.  Creates websocket
	if err := http.ListenAndServe(*addr, nil); err != nil {
		fmt.Printf("Error trying to bind to port: %v, so exiting...", err)
		log.Fatal("Error ListenAndServe:", err)
	}

}

// adcpHandler passes the template
// to the http request.
func adcpHandler(c http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("adcp.html")
	t.Execute(c, nil)
}
