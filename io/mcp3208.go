//go:build linux
package io


import (
	"gobot.io/x/gobot/v2/drivers/spi"
)

mcp3208_drivers := make(map[string](*MCP3208Driver))

func MCP3208_enable(data *Card) {
	for k,v := range SPI_BUS_SELECTOR{
		Ouput_set(v, (k==card.BusAddr))
	}
}

func MCP3208_init(data *Card)(*Card) {
	data.AddrCount=8
	data.WordSize=1
	data.Status="READY"
	Ouput_init(SPI_BUS_SELECTOR[card.BusAddr])
	mcp3208_drivers[card.BusAddr] = spi.NewMCP3208Driver(board)
	return data
}

func MCP3208_read(data *Card)(res []uint8, err error) {
	res = make([]uint8, 8)
	MCP3208_enable(data)
	for k, _ := range res {
		r, err := mcp3208_drivers[card.BusAddr].AnalogRead(string(k))
		if(err != nil) {
			return err
		}
		res[k] = byte(r/4)
	}
	return res, nil
}
