package api

import(
		io "msa/io-controller/io"
)


type WebsocketHello struct {
	MsgType string
	Inventory []io.Card
	Data []byte
}
