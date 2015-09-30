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
	http.HandleFunc("/", debugHandler)      // Debugger
	http.HandleFunc("/adcp", adcpHandler)   // Adcp Display
	http.HandleFunc("/adcp1", adcp1Handler) // Adcp1 Display
	http.HandleFunc("/adcp2", adcp2Handler) // Adcp2 Display
	http.HandleFunc("/adcp3", adcp3Handler) // Adcp3 Display
	http.HandleFunc("/adcp4", adcp4Handler) // Adcp4 Display
	http.HandleFunc("/adcp5", adcp5Handler) // Adcp5 Display
	http.HandleFunc("/adcp6", adcp6Handler) // Adcp6 Display
	http.HandleFunc("/adcp7", adcp7Handler) // Adcp7 Display
	http.HandleFunc("/adcp8", adcp8Handler) // Adcp8 Display
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

// adcp1Handler passes the template
// to the http request.
func adcp1Handler(c http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("adcp1.html")
	t.Execute(c, nil)
}

// adcp2Handler passes the template
// to the http request.
func adcp2Handler(c http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("adcp2.html")
	t.Execute(c, nil)
}

// adcp3Handler passes the template
// to the http request.
func adcp3Handler(c http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("adcp3.html")
	t.Execute(c, nil)
}

// adcp4Handler passes the template
// to the http request.
func adcp4Handler(c http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("adcp4.html")
	t.Execute(c, nil)
}

// adcp5Handler passes the template
// to the http request.
func adcp5Handler(c http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("adcp5.html")
	t.Execute(c, nil)
}

// adcp6Handler passes the template
// to the http request.
func adcp6Handler(c http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("adcp6.html")
	t.Execute(c, nil)
}

// adcp7Handler passes the template
// to the http request.
func adcp7Handler(c http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("adcp7.html")
	t.Execute(c, nil)
}

// adcp8Handler passes the template
// to the http request.
func adcp8Handler(c http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("adcp8.html")
	t.Execute(c, nil)
}
