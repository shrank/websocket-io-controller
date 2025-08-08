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
	data.Status="READY"
	return data
}