# io2websocket-gateway

This is a small webserver that lets you control I/O ports using a websocket based API. It currently supports MCP23017

```
                            ______________      i2c
Websocket Client --------->| io2websocket | --------------> MCP23017 Card #1
                           |______________|        |------> MCP23017 Card #2
                                                   +------> MCP23017 Card #3
```

## config

The config is done in a csv file:

```
BusAddr,Type,StartAddr,Mode
0,MCP23017,0,IN
1,MCP23017,16,IN
2,MCP23017,32,IN
3,MCP23017,48,IN
4,MCP23017,64,OUT
5,MCP23017,80,OUT
6,MCP23017,96,OUT
7,MCP23017,112,OUT
```

Bus Addresses are assiged as follows:
0-7   i2c Bus #1
10-17 i2c Bus #2 (planned)
20-27 SPI Bus #1 (planned)

## Websocket Protocoll

When opening a new connection, the server sends a complete dump of all cards and I/O data:

```
WebsocketHello {
	MsgType: "login",
	Inventory: [],      # array of card data
	Data: ""            # base64 encoded dump of all I/O data
}
```

Here is an example to convert the base64 string to an array:
```
function base64ToArray(base64) {
    var binaryString = atob(base64);
    var bytes = []
    for (var i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes;
}
```

The server sends an update message for every change. The same message format can be used to change data on the server.
```
WebsocketUpdate {
	MsgType: "update"
  Data: { "1": 1, "10": 1}      # Dictonary of updated I/O data 
}
```

## Build

The project is built with `build_image.sh`. It builds the frontend, cross-compiles the server for Raspberry Pi, packages everything as a Tiny Core overlay, and optionally embeds it into a Tiny Core image.

```
./build_image.sh <path-to-tinycore-image>
```

What the script does:
1. Build the frontend: `npm install && npm run build`
2. Cross-compile the Go server for `linux/arm64` (default, configurable)
3. Create a gzip-compressed cpio overlay containing the binary, `config.txt`, frontend assets, an init script, and a `bootlocal.sh` hook
4. Embed the overlay into the supplied Tiny Core image by copying it to `/overlays/` and patching `/config.txt`

Environment variables:
```
VERSION     package version (default: 0.0.1)
BUILD       build number (default: current timestamp)
GOARCH      target Go architecture (default: arm64)
GOARM       target ARM variant when GOARCH=arm (default: unset)
TC_OUT      output image name (default: tc-io2websocket-gateway-<version>-<build>.img)
OVERLAY_OUT output overlay name (default: overlay-io2websocket-gateway-<version>-<build>.gz)
```

For a 32-bit ARM build use:
```
GOARCH=arm GOARM=7 ./build_image.sh <path-to-tinycore-image>
```

The built overlay and image are written to `./build/`.
