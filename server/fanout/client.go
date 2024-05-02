package fanout

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// The size of the go channels for this client.
	channelSize = 10
)

// Client provides a structure for incomming websocket requests so they can be
// tracked by a Fanout.
type Client struct {
	fo   *Fanout
	conn *websocket.Conn
	send chan WebsocketEvent
	quit chan bool
}

// NewClient instantiates a client with a websocket connection. It spwans off
// go routines to send data to the client and send pings to ensure the client
// is still alive. We're passing a Fanout here so the client can unregister
// itself on error.
func (f *Fanout) NewClient(conn *websocket.Conn) *Client {
	c := &Client{
		fo:   f,
		conn: conn,
		send: make(chan WebsocketEvent, channelSize),
		quit: make(chan bool, 1),
	}

	f.RegisterClient(c)

	go c.writer()
	go c.reader()

	return c
}

// Send is used to send an image message to the client.
func (c *Client) Send(event WebsocketEvent) {
	c.send <- event
}

// Quit will close the connection and unregiseter it from the Fanout.
func (c *Client) Quit() {
	c.fo.UnregisterClient(c)
	c.quit <- true
	c.conn.Close()
}

// reader reads pong messages off of the connection. Once it recieves a message,
// the deadline is updated for when the next message must be recieved. If we
// don't get a message within the deadline, this method calls Quit to clean up
// the client.
func (c *Client) reader() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			c.Quit()
			break
		}
	}
}

// writer writes image events over the socket when it recieves messages via
// Send(). It also sends pings to ensure the connection stays alive.
func (c *Client) writer() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-c.quit:
			return
		case event := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteJSON(event); err != nil {
				c.Quit()
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Quit()
			}
		}
	}
}
