import { defineStore } from "pinia";

export const SessionDataStore = defineStore("session", {
  state: () => ({
    // internal
    ws: null,
    connected: false,

    // data
    data: {},
    
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
      this.refreshAll();
    },
    handleUpdate(event) {
      console.log("handle Update", event);
      if(("MsgType" in event) == false){
        return
      }
      switch(event.MsgType) {
        case "login":
          this.data = event.Data
          return
        case "add":
          for (let a in event.Data) {
            this.data[a] = event.Data[a]
          }
          return
        case "update":
          for (let a in event.Data) {
            this.data[a] = event.Data[a]
          }
          return
      }
    },

    refreshAll() {
    },
  },
  getters: {
  },
});
