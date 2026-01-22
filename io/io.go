package io

import (
	"log"
	"fmt"
	"strings"
	"time"
  "sync"
	"strconv"
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
	ActivityPin int
	IActivityPin int
	ActivityDiv int
	updateMap map[int]uint8
	updateLock sync.Mutex
}

func (self *IoV1) Init() error {
	max_addr := 0
	self.BufferSize = 0
	self.ActivityDiv = 50
	self.ActivityPin = -1
	self.IActivityPin = -1
	self.updateMap = make(map[int]uint8)
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

func (self *IoV1) Update(data map[int]uint8 ) {
	self.updateLock.Lock()
  defer self.updateLock.Unlock()
	for key, value := range data {
		self.updateMap[key] = value
	}
}



func (self *IoV1) Run() error {
	count := 0
	act := 1 
	if(self.ActivityPin > 0) {
		Output_init(strconv.Itoa(self.ActivityPin))
	}
	if(self.IActivityPin > 0) {
		Output_init(strconv.Itoa(self.IActivityPin))
	}
	fmt.Printf("Starting main IO loop\n")
	for true {
		count += 1
		update := make(map[int]uint8)

		if(self.IActivityPin > 0) {
			Output_set(strconv.Itoa(self.IActivityPin), 0)
		}
		if(self.ActivityPin > 0 && count % self.ActivityDiv == 0 ) {
			act = (act + 1) % 2
			Output_set(strconv.Itoa(self.ActivityPin), byte(act))
		}

		self.updateLock.Lock()
		if(len(self.updateMap) > 0) {
			for key, value := range self.updateMap {
				if(key < len(self.DataBuffer)) {
					if(self.DataBuffer[key] != value) {
						self.DataBuffer[key] = value
						update[key] = value
					}
				}	else	{
					fmt.Printf("address out of range: %d", key)
				}
			}
			self.updateMap = make(map[int]uint8)
		}
		self.updateLock.Unlock()

		for _, card := range self.Inventory {
			var err error
			var res []uint8

			if(card.Ready == false) {
				continue
			}

			for key, _ := range update {
				if(key >= card.StartAddr && key < card.StartAddr + card.AddrCount) {
					if(strings.ToLower(card.Mode) == "out") {
						self.doUpdate(card)
					}
					break
				}
			}

			if(strings.ToLower(card.Mode) != "in" && strings.ToLower(card.Mode) != "ain") {
				continue
			}


			if(card.ReadEvery > 1) {
				// we try to read only one card every run
				if(count % card.ReadEvery != int(card.BusAddr) % card.ReadEvery) {
					if(card.InterruptPin != "") {
						i, _ := Interrupt_Fired(card.InterruptPin)
						if(i == false){
							continue
						}
						if(self.IActivityPin > 0) {
							Output_set(strconv.Itoa(self.IActivityPin), 1)
						}
						fmt.Printf("interrupt card %d\n", card.BusAddr)					
					} else {
						continue
					}
				}
			}

			fmt.Printf("read card %d\n", card.BusAddr)

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
		}
		if(len(update) > 0 ) {
			self.UpdateHandler("update", update)
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