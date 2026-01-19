package api

import(
  "fmt"
  "log"
  "encoding/json"
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return true
},
}



func (api ApiV1) WsHello(ws *websocket.Conn) {
  // send some initial data here
	data := WebsocketHello{}
	data.MsgType = WS_LOGIN
	data.Inventory = append(data.Inventory, api.Hardware.Inventory...)
	data.Data = append(data.Data, api.Hardware.DataBuffer...)
  
	jsonResp, _ := json.Marshal(data)
	ws.WriteMessage(websocket.TextMessage, []byte(jsonResp))
}

func (api ApiV1) WsConnect (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  // Setup websocket
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println("upgrade:", err)
    return
  }
  fmt.Printf("WS connected\n")

  defer ws.Close()

  done := false
  current_message, _, _ :=api.MsgQueue.Last()

	api.WsHello(ws)

	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			res := WebsocketUpdate{}
			json.Unmarshal(message, &res)
			if(res.MsgType == "update") {
					api.Hardware.Update(res.Data)
			}
			log.Printf("recv: %s", message)
		}
	}()

  for {
    cnt, res, err := api.MsgQueue.Wait(current_message)		 
		current_message = cnt
		msg, ok := res.(WebsocketUpdate)
		if(ok && err == nil) {
			jsonResp, _ := json.Marshal(msg)
      ws.WriteMessage(websocket.TextMessage, []byte(jsonResp))
    }
    if done {
      break
    }
  }
}