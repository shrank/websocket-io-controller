//go:build windows
package io

import (
	"fmt"
	"math/rand"
)

func Raspi_init()(error) {
	return nil
}


func Interrupt_init(pin string)(error) {
	fmt.Printf("Setup Interrupt Pin %s\n", pin)
	return nil
}

func Interrupt_Fired(pin string)( bool, error) {
	if(rand.Intn(255) == 0){
		return true, nil
	}
	return false, nil
}

func Output_init(pin string)(error) {
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
	data.Ready=true
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
	data.Ready=true
	data.ReadEvery=10
	return data
}

func MCP3208_read(data *Card)(res []uint8, err error) {
	res = make([]uint8, 8)
	for k, _ := range res {
		res[k] = uint8(rand.Intn(255))
	}
	return res, nil
}

func MCP3208_read_one(data *Card, channel int)(res uint8, err error) {
	return uint8(rand.Intn(255)), nil
}
