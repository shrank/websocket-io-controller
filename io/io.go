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
	InterruptPin string
}


type IoV1 struct {
	DataBuffer []uint8
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

		if(strings.ToLower(card.Type) == "mcp3208") {
			MCP3208_init(&self.Inventory[c])
		}

		end_addr := self.Inventory[c].StartAddr + self.Inventory[c].AddrCount
		if(end_addr > max_addr) {
			max_addr = end_addr
		}
	}
	self.BufferSize = max_addr
	fmt.Printf("initializing buffer of %d bytes\n", self.BufferSize)
	self.DataBuffer = make([]uint8, self.BufferSize, self.BufferSize)
  return nil
}

func (self *IoV1) Update(data map[int]uint8 ) map[int]uint8 {
	res := make(map[int]uint8)
	for key, value := range data {
		if(key < len(self.DataBuffer)) {
			if(self.DataBuffer[key] != value) {
				self.DataBuffer[key] = value
				res[key] = value
			}
		}	else	{
			fmt.Printf("address out of range: %d", key)
		}
	}
	for _, card := range self.Inventory {
		if(strings.ToLower(card.Mode) != "out") {
			continue
		}
		for key, _ := range res {
			if(key >= card.StartAddr && key < card.StartAddr + card.AddrCount) {
				self.doUpdate(card)
				break
			}
		}
	}
	return res
}

func (self *IoV1) doUpdate(card Card){
		if(strings.ToLower(card.Type) == "mcp23017") {
			MCP23017_update(&card, self.DataBuffer[card.StartAddr:card.StartAddr + card.AddrCount])
		}
}