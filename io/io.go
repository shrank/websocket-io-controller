package io

import (
	"fmt"
	"strings"
)


type Card struct {
	Type string
	StartAddr int
	AddrCount int
	WordSize int
	Mode string
	BusAddr byte
	Status string
}


type IoV1 struct {
	DataBuffer []byte
	Inventory []Card
	BufferSize int
}

func (self *IoV1) Init() error {
	max_addr := 0
	self.BufferSize = 0
	for c, card := range self.Inventory {
		fmt.Printf("Init Card #%d\n", c)
		if(strings.ToLower(card.Type) == "mcp23017") {
			MCP23017_init(&self.Inventory[c])
		}

		end_addr := self.Inventory[c].StartAddr + self.Inventory[c].AddrCount
		if(end_addr > max_addr) {
			max_addr = end_addr
		}
		self.BufferSize += (self.Inventory[c].AddrCount * self.Inventory[c].WordSize) / 8
		if( (self.Inventory[c].AddrCount * self.Inventory[c].WordSize) % 8 > 0) {
			self.BufferSize += 1
		}
	}
	fmt.Printf("initializing buffer of %d bytes\n", self.BufferSize)
	self.DataBuffer = make([]byte, self.BufferSize, self.BufferSize)
  return nil
}

