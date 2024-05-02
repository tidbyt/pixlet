package fanout

const (
	// EventTypeImage is used to signal what type of message we are sending over
	// the socket.
	EventTypeImage = "img"

	// EventTypeSchema is used to signal that the schema for a given app has
	// changed.
	EventTypeSchema = "schema"

	// EventTypeErr is used to signal there was an error encountered rendering
	// the image.
	EventTypeErr = "error"
)

// WebsocketEvent is a structure used to send messages over the socket.
type WebsocketEvent struct {
	// Message is the contents of the message. This is a webp or gif, base64 encoded.
	Message string `json:"message"`

	// ImageType indicates whether the Message is webp or gif image.
	ImageType string `json:"img_type"`

	// Type is the type of message we are sending over the socket.
	Type string `json:"type"`
}
