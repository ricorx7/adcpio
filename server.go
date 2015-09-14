package main

import "log"

// adcpIO Connections and broadcast to
// send messages to all connections.
// A connection is either a serial port
// or websocket.
type adcpIO struct {
	websocketConn map[*websocketConn]bool // Registered connections.
	register      chan *websocketConn     // Register requests from the connections.
	unregister    chan *websocketConn     // Unregister requests from connections.
	broadcast     chan []byte             // Broadcast data
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
}

// run the Echo process
// This will monitor websockets
// and serial ports for connections
// and disconnects.
func (server *adcpIO) run() {
	log.Print("Echo Hub running")
	for {
		select {

		// Register websocket
		case c := <-server.register:
			// Register the websocket to the map
			server.websocketConn[c] = true

			log.Println("Registering websocket")

		// Unregister websocket
		case c := <-server.unregister:
			if _, ok := server.websocketConn[c]; ok {

				log.Println("Unregistering websocket")

				delete(server.websocketConn, c) // Unregister the websocket from the map
				close(c.send)                   // Close the websocket send channel
			}

			// Broadcast message to all listeners
		case m := <-server.broadcast:
			log.Print(m)

			for c := range server.websocketConn {
				log.Print(c.adcpSerialNum)
				select {
				case c.send <- m:
				default:
					// If connection not good, remove the connection
					close(c.send)
					delete(server.websocketConn, c)
				}
			}

		}
	}
}
