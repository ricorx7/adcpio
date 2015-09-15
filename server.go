package main

import (
	"encoding/json"
	"log"

	"github.com/ricorx7/go-rti"
)

type adcp struct {
	serialNum string
	lastEns   rti.Ensemble
}

// adcpIO Connections and broadcast to
// send messages to all connections.
// A connection is either a serial port
// or websocket.
type adcpIO struct {
	websocketConn map[*websocketConn]bool // Registered connections.
	register      chan *websocketConn     // Register requests from the connections.
	unregister    chan *websocketConn     // Unregister requests from connections.
	broadcast     chan []byte             // Broadcast data
	adcp          map[string]*adcp        // List of ADCP data.  Key is the serial number of the ADCP
}

// echo initializes the values.
// This will hold all the registered websocket
// connections.  It will also hold the send and receive
// buffer from the websockets.
var server = adcpIO{
	register:      make(chan *websocketConn),     // Register a websocket connections
	unregister:    make(chan *websocketConn),     // Unregister a websocket connection
	websocketConn: make(map[*websocketConn]bool), // Websocket connection map
	broadcast:     make(chan []byte),             // Broadcast the data
	adcp:          make(map[string]*adcp),        // ADCP Data map
}

// run the server process
// This will monitor websockets
// and serial ports for connections
// and disconnects.
func (server *adcpIO) run() {
	log.Print("Echo Hub running")
	for {
		select {

		// Register websocket
		case c := <-server.register:
			log.Println("Registering websocket")
			// Register the websocket to the map
			server.websocketConn[c] = true

		// Unregister websocket
		case c := <-server.unregister:
			if _, ok := server.websocketConn[c]; ok {

				log.Println("Unregistering websocket")

				delete(server.websocketConn, c) // Unregister the websocket from the map
				close(c.send)                   // Close the websocket send channel
			}

			// Broadcast message to all listeners
		case m := <-server.broadcast:

			//log.Printf("Message broadcast: %d", len(m))
			//log.Print(string(m))

			// Convert the message to JSON
			var ens rti.Ensemble
			err := json.Unmarshal(m, &ens)
			if err != nil {
				log.Print("Err converting JSON: ", err)
			}
			//log.Printf("Ensemble Number: %d", ens.EnsembleData.EnsembleNumber)

			// Process the ensemble
			processEnsemble(server, ens)
		}
	}
}

func processEnsemble(server *adcpIO, ens rti.Ensemble) {
	// See if the serial number exist in the map
	if val, ok := server.adcp[ens.EnsembleData.SerialNumber.SerialNumber]; ok {
		val.lastEns = ens
		log.Print("ADCP exist")
	} else {
		var data adcp
		data.lastEns = ens
		server.adcp[ens.EnsembleData.SerialNumber.SerialNumber] = &data
		log.Print("ADCP does not exist")
	}
}
