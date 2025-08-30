import { defineStore } from "pinia";




function base64ToArrayBuffer(base64) {
    var binaryString = atob(base64);
    var bytes = []
    for (var i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes;
}



export const SessionDataStore = defineStore("session", {
  state: () => ({
    // internal
    ws: null,
    connected: false,

    // data
    data: {},
    bytedata: [],
    inventory: [],
  }),
  actions: {
    // Websocket and live updates
    connect() {
      console.log("try WebSocket connect");
      let url = new URL(window.location.href);
      let proto = "wss://";
      if (url.protocol == "http:") {
        proto = "ws://";
      }
      this.ws = new WebSocket(proto + url.host + "/api/v1/live");
      this.ws.onopen = () => {
        this.connected = true
        console.log("WebSocket connection opened:");
        setInterval(() => { this.ws.send('__ping__'); }, 30000);
      };
      this.ws.onmessage = (event) => {
        this.handleUpdate(JSON.parse(event.data));
      };
      this.ws.onerror = (error) => {
        console.log("WebSocket error:", error);
      };
      this.ws.onclose = (event) => {
        this.connected = false
        console.log("WebSocket connection closed:", event.code);
      };
    },
    handleUpdate(event) {
      console.log("handle Update", event);
      if(("MsgType" in event) == false){
        return
      }
      switch(event.MsgType) {
        case "login":
          if ("Inventory" in event) {
            this.inventory = event.Inventory
          }
          if("Data" in event) {
            this.bytedata = base64ToArrayBuffer(event.Data)
            console.log(this.bytedata)
          }
          return
        case "add":
          for (let a in event.Data) {
            this.data[a] = event.Data[a]
          }
          return
        case "update":
          for (let a in event.Data) {
            this.bytedata[a] = event.Data[a]
          }
          return
      }
    },
    IoUpdate(data) {
      if(this.ws != null) {
        this.ws.send(JSON.stringify({"MsgType": "update", "data": data }))
      }
    }
  },
});
