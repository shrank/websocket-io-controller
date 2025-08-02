package main

import (
  "fmt"
  "net/http"
	"log"
	"time"
  "github.com/julienschmidt/httprouter"
	api_class "msa/io-controller/api"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "Welcome!\n")
}

func main() {

  router := httprouter.New()
	inventory := []api_class.Card {		
		api_class.Card{
			Type: "Test1",
			StartAddr: 0,
			AddrCount: 16,
			WordSize: 1,
			Mode: "Digital Input",
			BusAddr: 1,
		},
	}
	api := api_class.NewAPI(&inventory)

	// Index Page
  router.GET("/",Index)

	// static files
  router.ServeFiles("/debug/*filepath", http.Dir("./public"))
	
  // Websockets
  router.GET("/api/v1/live", api.WsConnect)
  

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
