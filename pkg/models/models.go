package models

type Message struct {
	Id      string // 4 bytes
	Sign    string // 2 bytes
	Content string // 8 kiloBytes
}

func NewMessage(frame *[]byte) *Message {
	// id = first 4 byte ...
	return nil
}

type Connection struct {
	Destination string
	//id    string // for this connection
	Sign  string // use for messages
	Count int
}

