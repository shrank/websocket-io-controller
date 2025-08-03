package io


func MCP23017_init(data *Card)(*Card) {
	data.Status="INIT Failed"
	data.AddrCount=16
	data.WordSize=1
	return data
}