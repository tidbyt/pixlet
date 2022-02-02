package fanout

// Fanout provides a structure for broadcasting messages to registered clients
// when an update comes in on a go channel.
type Fanout struct {
	broadcast  chan WebsocketEvent
	quit       chan bool
	register   chan *Client
	unregister chan *Client
}

// NewFanout creates a new Fanout structure and runs the main loop.
func NewFanout() *Fanout {
	fo := &Fanout{
		broadcast:  make(chan WebsocketEvent, channelSize),
		register:   make(chan *Client, channelSize),
		unregister: make(chan *Client, channelSize),
		quit:       make(chan bool, 1),
	}

	go fo.run()

	return fo
}

// Broadcast sends a message to all registered clients.
func (fo *Fanout) Broadcast(event WebsocketEvent) {
	fo.broadcast <- event
}

// RegisterClient registers a client to include in broadcasts.
func (fo *Fanout) RegisterClient(c *Client) {
	fo.register <- c
}

// UnregisterClient removes it from the broadcast.
func (fo *Fanout) UnregisterClient(c *Client) {
	fo.unregister <- c
}

// Quit stops broadcasting messages over the channel.
func (fo *Fanout) Quit() {
	fo.quit <- true
}

// run is the main loop. It provides a mechanism to register/unregister clients
// and will broadcast messages as they come in.
func (fo *Fanout) run() {
	clients := map[*Client]bool{}

	for {
		select {
		case <-fo.quit:
			for client := range clients {
				client.Quit()
			}
		case c := <-fo.register:
			clients[c] = true
		case c := <-fo.unregister:
			if _, ok := clients[c]; ok {
				delete(clients, c)
				c.Quit()
			}
		case broadcast := <-fo.broadcast:
			for client := range clients {
				client.Send(broadcast)
			}
		}
	}
}
