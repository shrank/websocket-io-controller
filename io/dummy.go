//go:build windows
package io

import (
	"fmt"
	"math/rand"
)

func Interrupt_init(pin string)(error) {
	fmt.Printf("Setup Interrupt Pin %s\n", pin)
	return nil
}

func Interrupt_Fired(pin string)( bool, error) {
	if(rand.Intn(2) > 0){
		return true, nil
	}
	return false, nil
}

func Ouput_init(pin string)(error) {
	fmt.Printf("Setup Output Pin %s\n", pin)
	return nil
}

func Output_set(pin string, value byte)(error) {
	return nil
}


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

func MCP23017_read(data *Card)(res []uint8, err error) {
	res = make([]uint8, 16)
	for k, _ := range res {
		res[k] = uint8(rand.Intn(2))
	}
	return res, nil
}


func MCP3208_init(data *Card)(*Card) {
	data.AddrCount=8
	data.WordSize=1
	data.Status="READY"
	return data
}

func MCP3208_read(data *Card)(res []uint8, err error) {
	res = make([]uint8, 8)
	for k, _ := range res {
		res[k] = uint8(rand.Intn(255))
	}
	return res, nil
}
