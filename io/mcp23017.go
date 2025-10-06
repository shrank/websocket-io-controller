//go:build linux
package io


import (
	"log"
	"strings"
)

mcp12017_drivers := make(map[string](*MCP23017Driver))

func MCP23017_init(data *Card)(*Card) {
	data.AddrCount=16
	data.WordSize=1

	i2c_lock.Lock()
	defer m.mutex.Unlock()
	
	d := i2c.NewMCP23017Driver(
		board,  
		i2c.WithBus(0), 
		i2c.WithAddress(0x20 + data.BusAddr),
		i2c.WithMCP23017Bank(0),		// use bank 0
		i2c.WithMCP23017Intpol(1),	// interrupt polarity high
		i2c.WithMCP23017Seqop(0),		// enable sequencial operation
		i2c.WithMCP23017Mirror(1)		// mirror interrupt pins
	)

	if(strings.ToLower(data.Mode) == "out") {
		for i := 0; i < 8; i++ {
			d.PinMode(i, 1, "A")
			d.PinMode(i, 1, "B")
		}
	}

	mcp12017_drivers[card.BusAddr] = d
	data.Status="READY"
	return data
}

func MCP23017_update(c *Card, d []uint8)(error) {
	value := uint16(0)
	for _, v := range d {
		value *= 2
		if(v > 0 ) {
			value += 1
		}
	}
	i2c_lock.Lock()
	defer m.mutex.Unlock()
	err := mcp12017_drivers[card.BusAddr].WriteWordData(0x12, value)
	return err
}

func MCP23017_read(data *Card)(res []uint8, err error) {
	res = make([]uint8, 16)
	c := mcp12017_drivers[card.BusAddr].Connection()

	i2c_lock.Lock()
	defer m.mutex.Unlock()

	val, err = c.ReadWordData(0x12)
	if err != nil {
		return res, err
	}
	for k, _ := range res {
		res[k] = val % 2
		val = val / 2
	}
	return res, nil
}
