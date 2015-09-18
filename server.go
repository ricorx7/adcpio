package main

import (
	"encoding/json"
	"log"

	"github.com/ricorx7/go-rti"
)

// adcpIO Connections and broadcast to
// send messages to all connections.
// A connection is either a serial port
// or websocket.
type adcpIO struct {
	websocketConn         map[*websocketConn]bool        // Registered connections.
	register              chan *websocketConn            // Register requests from the connections.
	unregister            chan *websocketConn            // Unregister requests from connections.
	wsAdcpDisplayConn     map[*websocketAdcpDisplay]bool // Registered Adcp Display connections.
	registerAdcpDisplay   chan *websocketAdcpDisplay     // Register requests from Adcp Display connections.
	unregisterAdcpDisplay chan *websocketAdcpDisplay     // Unregister requests from Adcp Display connections.
	broadcast             chan []byte                    // Broadcast data
	adcp                  map[string]*adcp               // List of ADCP data.  Key is the serial number of the ADCP
}

// echo initializes the values.
// This will hold all the registered websocket
// connections.  It will also hold the send and receive
// buffer from the websockets.
var server = adcpIO{
	register:              make(chan *websocketConn),            // Register a websocket connections
	unregister:            make(chan *websocketConn),            // Unregister a websocket connection
	websocketConn:         make(map[*websocketConn]bool),        // Websocket connection map
	registerAdcpDisplay:   make(chan *websocketAdcpDisplay),     // Register a websocket connections
	unregisterAdcpDisplay: make(chan *websocketAdcpDisplay),     // Unregister a websocket connection
	wsAdcpDisplayConn:     make(map[*websocketAdcpDisplay]bool), // Websocket connection map
	broadcast:             make(chan []byte),                    // Broadcast the data
	adcp:                  make(map[string]*adcp),               // ADCP Data map
}

// adcp will store all the ADCP it is monitoring and also the last ensemble.
type adcp struct {
	serialNum string       // Serial number
	lastEns   rti.Ensemble // Last ensemble
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

		// Register Adcp Display websocket
		case c := <-server.registerAdcpDisplay:
			log.Println("Registering Adcp Display websocket")
			// Register the websocket to the map
			server.wsAdcpDisplayConn[c] = true

			// Send a list of all ADCP
			sendAdcpList()

		// Unregister Adcp Display websocket
		case c := <-server.unregisterAdcpDisplay:
			if _, ok := server.wsAdcpDisplayConn[c]; ok {

				log.Println("Unregistering websocket")

				delete(server.wsAdcpDisplayConn, c) // Unregister the websocket from the map
				close(c.send)                       // Close the websocket send channel

				// Send a new list of all the ADCP
				sendAdcpList()
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
			// Pass the data to all the registered displays
			processEnsemble(server, ens)
		}
	}
}

// sendDataToDisplays will send data to all the registered displays.
func sendDataToDisplays(b []byte) {
	for c := range server.wsAdcpDisplayConn {
		select {
		case c.send <- b:
		default:
			log.Print("Close Adcp Display websocket send")
			close(c.send)
			delete(server.wsAdcpDisplayConn, c)
		}
	}
}

// sendAdcpList will send list of ADCP connected
// to all the registered displays.
func sendAdcpList() {
	var list []string
	for key, _ := range server.adcp {
		list = append(list, key)
	}

	adcps := adcpList{
		ID:            adcpListID, // ID
		SerialNumList: list,       // List of serial numbers
	}

	// Convert the JSON to byte array
	b, err := json.Marshal(adcps)
	if err != nil {
		log.Println(err)
		return
	}

	// Send the data to the display
	sendDataToDisplays(b)

}
