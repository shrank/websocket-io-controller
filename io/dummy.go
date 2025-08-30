//go:build windows
package io

import (
	"fmt"
)

func MCP23017_init(data *Card)(*Card) {
	data.AddrCount=16
	data.WordSize=1
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
	fmt.Printf("update card %d: %x\n", c.BusAddr, value)
	return nil
}