//go:build linux
package io


import (
	"strings"
	"gobot.io/x/gobot/v2/drivers/i2c"
)

var mcp12017_drivers = make(map[byte](*i2c.GenericDriver))

func MCP23017_init(data *Card)(*Card) {
	data.AddrCount=16
	data.WordSize=1

	i2c_lock.Lock()
	defer i2c_lock.Unlock()
	
	d := i2c.NewMCP23017Driver(
		board,  
		i2c.WithBus(0), 
		i2c.WithAddress(int(0x20 + data.BusAddr)),
		i2c.WithMCP23017Bank(0),		// use bank 0
		i2c.WithMCP23017Intpol(1),	// interrupt polarity high
		i2c.WithMCP23017Seqop(0),		// enable sequencial operation
		i2c.WithMCP23017Mirror(1),		// mirror interrupt pins
	)

	if(strings.ToLower(data.Mode) == "out") {
		for i := 0; i < 8; i++ {
			d.SetPinMode(uint8(i), "A", 0)
			d.SetPinMode(uint8(i), "B", 0)
		}
	}

	mcp12017_drivers[data.BusAddr] = i2c.NewGenericDriver(board, "mcp12017", int(0x20 + data.BusAddr), i2c.WithBus(1))
	err := mcp12017_drivers[data.BusAddr].Start()

	if(err != nil) {
		data.Status=err.Error()
		data.Ready = false
	} else {
		data.Ready = true
		data.Status="READY"
	}
	return data
}

func MCP23017_update(card *Card, d []uint8)(error) {
	value := uint16(0)
	for _, v := range d {
		value *= 2
		if(v > 0 ) {
			value += 1
		}
	}
	i2c_lock.Lock()
	defer i2c_lock.Unlock()
	err := mcp12017_drivers[card.BusAddr].WriteWordData(0x12, value)
	return err
}

func MCP23017_read(card *Card)(res []uint8, err error) {
	res = make([]uint8, 16)

	i2c_lock.Lock()
	defer i2c_lock.Unlock()

	val, err := mcp12017_drivers[card.BusAddr].ReadWordData(0x12)
	if err != nil {
		return res, err
	}
	for k, _ := range res {
		res[k] = uint8(val % 2)
		val = val / 2
	}
	return res, nil
}
