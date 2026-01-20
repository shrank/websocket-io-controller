//go:build linux
package io


import (
	"strings"
	"fmt"
	"gobot.io/x/gobot/v2/drivers/i2c"
)

var mcp12017_drivers = make(map[byte](*i2c.GenericDriver))

func MCP23017_init(data *Card)(*Card) {
	data.AddrCount=16
	data.WordSize=1
	data.Ready = false
	i2c_lock.Lock()
	defer i2c_lock.Unlock()
	
	d := i2c.NewGenericDriver(board, "mcp12017", int(0x20 + data.BusAddr), i2c.WithBus(1))
	err := d.Start()
	if(err != nil) {
		data.Status=err.Error()
		fmt.Printf("ERROR: %s\n", err.Error())
		return data
	} 

	err = d.WriteByteData(0x0a, 2 + 64)	// IOCON: MIRROR + INTPOL=high 

	if(err != nil) {
		data.Status=err.Error()
		fmt.Printf("ERROR: %s\n", err.Error())
		return data
	}

	if(strings.ToLower(data.Mode) == "out") {
		d.WriteWordData(0x12, 0xffff)  //set all pins to high=off
		d.WriteWordData(0x14, 0xffff)  //set all pins to high=off
		d.WriteWordData(0x0,0x0)	//set direction output
		d.WriteWordData(0x12, 0xffff)  //set all pins to high=off
	} else {
		d.WriteWordData(0x2,0xffff)	// Invert Polarity
		d.WriteWordData(0x4,0xffff)	// Enable Interrupts
		d.WriteWordData(0x0,0xffff)	//set direction output		
	}

	mcp12017_drivers[data.BusAddr] = d
	data.Ready = true
	data.Status="READY"
	return data
}

func MCP23017_update(card *Card, d []uint8)(error) {
	value := uint16(0)
	for i := (len(d) -1); i >= 0; i-=1 {
		value *= 2
		if(d[i] == 0 ) {
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
