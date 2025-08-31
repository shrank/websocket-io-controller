//go:build linux
package io


import (
	"log"
	"github.com/googolgl/go-i2c"
	"github.com/googolgl/go-mcp23017"
	"strings"
)

func MCP23017_init(data *Card)(*Card) {
	data.AddrCount=16
	data.WordSize=1
	i2c, err := i2c.New(mcp23017.DefI2CAdr + data.BusAddr, 1)
	if err != nil {
			log.Print(err)
			data.Status="INIT Failed"
			return data
	}

	mcp, err := mcp23017.New(i2c)
	if err != nil {
			log.Print(err)
			data.Status="INIT Failed"
			return data
	}

	if(strings.ToLower(data.Mode) == "out") {
		mcp.Set(mcp23017.AllPins()).OUTPUT()
	}

	// set IOCON to 0x0 to enable banmk	
	i2c.WriteRegU8(byte(mcp23017.IOCON), 0x0)

	data.Status="READY"
	return data
}

func MCP23017_update(c *Card, d []uint8)(error) {
	i2c, err := i2c.New(mcp23017.DefI2CAdr + c.BusAddr, 1)
	if err != nil {
		return err
	}

	value := uint16(0)
	for _, v := range d {
		value *= 2
		if(v > 0 ) {
			value += 1
		}
	}
	
	i2c.WriteRegU16BE(byte(0x12), value)
	return nil
}