package main

import (
  "fmt"
  "net/http"
	"log"
	"time"
  "github.com/julienschmidt/httprouter"
  utils "msa/io-controller/utils"
	api_class "msa/io-controller/api"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "Welcome!\n")
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

func main() {

	MsgQueue := utils.CreateQueue()   // WebSocket message queue
	UNUSED(MsgQueue)
  router := httprouter.New()
	api := api_class.NewAPI()
  // static files
  router.ServeFiles("/debug/*filepath", http.Dir("./public"))
	
  // Websockets
  router.GET("/api/v1/live", api.WsConnect)

  // router.GET("/api/v1/", A.Auth(Index, dal.PERM_READ))

  
  
  // @TD: this is not nice, the event bus should be able to handle that somehow
  // go func() {
  //   current_message, _, _ :=MsgQueue.Last()
  
  //   for {
  //     cnt, res, _ := MsgQueue.Wait(current_message)
	// 		msg, _ := res.(api.WebsocketUpdate)
  //     current_message = cnt
  //   }
  // }()

	// Do your work here
  go func() {
		for  true {
			time.Sleep(time.Second)
			t := time.Now()
			api.SendUpdate(api_class.WS_UPDATE,t)
		}
	}()

  fmt.Printf("Starting server at port 8000\n")
  if err := http.ListenAndServe(":8000", router); err != nil {
    log.Fatal(err)
  }
}
