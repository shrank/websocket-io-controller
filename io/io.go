package io

import (
	"log"
	"fmt"
	"strings"
	"time"
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
	Ready bool
	ReadEvery int
}

type UpdateHandler func(t string, v map[int]uint8) error

type IoV1 struct {
	UpdateHandler UpdateHandler
	DataBuffer []uint8
	Inventory []Card
	BufferSize int
}

func (self *IoV1) Init() error {
	max_addr := 0
	self.BufferSize = 0
	Raspi_init()
	for c, card := range self.Inventory {
		fmt.Printf("Init Card #%d\n", c)
		if(strings.ToLower(card.Type) == "mcp23017") {
			MCP23017_init(&self.Inventory[c])
		}

		if(strings.ToLower(card.Type) == "mcp3208") {
			MCP3208_init(&self.Inventory[c])
		}
		if(card.InterruptPin != "") {
			Interrupt_init(card.InterruptPin)
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


func (self *IoV1) Run() error {
	count := 0
	fmt.Printf("Starting main IO loop\n")
	for true {
		count += 1
		for _, card := range self.Inventory {
			var err error
			var res []uint8
			update := make(map[int]uint8)

			if(card.Ready == false) {
				continue
			}
			if(strings.ToLower(card.Mode) != "in" && strings.ToLower(card.Mode) != "ain") {
				continue
			}

			if(card.InterruptPin != "") {
				i, _ := Interrupt_Fired(card.InterruptPin)
				if(i == false){
					continue
				}
			} else {
				if(card.ReadEvery > 1) {
					// we try to read only one card every run
					if(count % card.ReadEvery != 10 % int(card.BusAddr)) {
						continue
					}
				}
			}

			// handle Digial Input Cards
			if(strings.ToLower(card.Type) == "mcp23017") {
				res, err = MCP23017_read(&card)
			}

			// handle Analog Input Cards
			if(strings.ToLower(card.Type) == "mcp3208") {
				res, err = MCP3208_read(&card)
			}

			if(err != nil) {
				log.Fatal(err)
			}

			// Send Update
			i := 0
			for i < card.AddrCount {
				if(i>= len(res)) {
					break
				}
				if(self.DataBuffer[i+card.StartAddr] != res[i]) {
					self.DataBuffer[i+card.StartAddr] = res[i]
					update[i+card.StartAddr] = res[i]			
				}
				i += 1
			}
			if(len(update) > 0 ) {
				self.UpdateHandler("update", update)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}	
	return nil
}

func (self *IoV1) doUpdate(card Card){
		if(card.Ready == false) {
			return
		}

		if(strings.ToLower(card.Type) == "mcp23017") {
			MCP23017_update(&card, self.DataBuffer[card.StartAddr:card.StartAddr + card.AddrCount])
		}
}