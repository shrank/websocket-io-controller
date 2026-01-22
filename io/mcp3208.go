//go:build linux
package io


import (
	"gobot.io/x/gobot/v2/drivers/spi"
	"strconv"
)

var mcp3208_drivers = make(map[byte](*spi.MCP3208Driver))

func MCP3208_init(data *Card)(*Card) {
	data.AddrCount=8
	data.WordSize=1
	data.ReadEvery=10
	mcp3208_drivers[data.BusAddr] = spi.NewMCP3208Driver(board, spi.WithChipNumber(int(data.BusAddr)-20))
	err := mcp3208_drivers[data.BusAddr].Start()
	if(err != nil) {
		data.Status=err.Error()
		data.Ready = false
	} else {
		data.Ready = true
		data.Status="READY"
	}

	return data
}

func MCP3208_read(data *Card)(res []uint8, err error) {
	res = make([]uint8, 8)
	for k, _ := range res {
		r, err := MCP3208_read_one(data, k)
		if(err != nil) {
			return res, err
		}
		res[k] = r
	}
	return res, nil
}

func MCP3208_read_one(data *Card, channel int)(res uint8, err error) {
	r, err := mcp3208_drivers[data.BusAddr].AnalogRead(strconv.Itoa(channel))
	if(err != nil) {
		return byte(r), err
	}
	return byte(r/4), nil
}
