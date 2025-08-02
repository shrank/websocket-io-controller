package api

//import()

type Card struct {
	Type string
	StartAddr int
	AddrCount int
	WordSize int
	Mode string
	BusAddr int
}

type WebsocketHello struct {
	MsgType string
	Inventory []Card
	Data []int
}
