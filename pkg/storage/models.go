package storage

type connSend struct {
	Destination string
	Count       int
	Id          string // 20 char
	IsOpen      int    // open 1, closed 0
	Successful  int    //  1 - 0
}
