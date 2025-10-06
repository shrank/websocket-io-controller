//go:build linux
package io


import (
)

func MCP3208_init(data *Card)(*Card) {
	data.AddrCount=8
	data.WordSize=1
	data.Status="READY"
	return data
}
