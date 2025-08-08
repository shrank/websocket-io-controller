//go:build windows
package io

func MCP23017_init(data *Card)(*Card) {
	data.AddrCount=16
	data.WordSize=1
	data.Status="READY"
	return data
}