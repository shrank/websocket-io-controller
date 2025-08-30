// The api package contains all methods accessible by the web API
package api

import (
  "strconv"
  "encoding/json"
  "net/http"
  "github.com/julienschmidt/httprouter"
  utils "msa/io-controller/utils"
  io "msa/io-controller/io"
)

const (
  WS_LOGIN string  = "login"
  WS_UPDATE        = "update"
  WS_ADD	  	     = "add"
  WS_DELETE        = "delete"
)

type WebsocketUpdate struct {
	MsgType string
  Data map[int]uint8 
}

type WebsocketHello struct {
	MsgType string
	Inventory []io.Card
	Data []uint8
}

// Main Api class, accessable by all API methods
type ApiV1 struct {
  MsgQueue *utils.Queue   // WebSocket message queue
	Hardware *io.IoV1        // Card inventory
}

func NewAPI(hwobj *io.IoV1) (ApiV1) {
  res := ApiV1{}
	res.MsgQueue = utils.CreateQueue()
	res.Hardware = hwobj
  return res
}


// send JSON reply to client
func ResponseJson(w http.ResponseWriter, v any) {
  w.Header().Set("Content-Type", "application/json")
  jsonResp, err := json.Marshal(v)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Write(jsonResp)
  return
}

// send 404 Error reply to client
func Response404(w http.ResponseWriter) {
  w.WriteHeader(http.StatusNotFound)
  resp := make(map[string]string)
  resp["message"] = "Resource Not Found"
  ResponseJson(w, resp)
  return
}

// send Sucess JSON reply to client
func ResponseSuccess(w http.ResponseWriter) {
  resp := make(map[string]string)
  resp["message"] = "success"
  ResponseJson(w, resp)
  return
}

func unpackArray[T any](arr any) (r []any) {
  o := arr.([]T)
  r = make([]any, len(o))

  for i, v := range o {
      r[i] = any(v)
  }

  return r
}

// send update to WebSocket Clients
func (api ApiV1) SendUpdates(t string, v map[int]uint8) error {
	api.MsgQueue.Insert(WebsocketUpdate {
		MsgType: t,
		Data: v,
		})
  return nil
}

// parse an ID from url path
func parseId(ps httprouter.Params) (success bool, res uint) {
  idstr := ps.ByName("id")
  if idstr == "" {
    return false, 0
  }
  id, err := strconv.ParseUint(idstr, 10, 64)
  if err != nil {
    return false, 0
  }
  return true, uint(id)
}
