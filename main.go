package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
  "fmt"
  "net/http"
	"log"
  "github.com/julienschmidt/httprouter"
	api_class "msa/io2websocket-gateway/api"
	io_class "msa/io2websocket-gateway/io"
)

// ReadCardsFromCSV reads a CSV file and returns a slice of Card structs
func ReadCardsFromCSV(filename string) ([]io_class.Card, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("could not read header: %w", err)
	}

	// Index map for column headers
	indices := map[string]int{}
	for i, h := range headers {
		indices[h] = i
	}

	var cards []io_class.Card

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read record: %w", err)
		}

		busAddr, _ := strconv.Atoi(record[indices["BusAddr"]])
		startAddr, _ := strconv.Atoi(record[indices["StartAddr"]])
		card := io_class.Card{
			BusAddr:   byte(busAddr),
			Type:      record[indices["Type"]],
			StartAddr: startAddr,
			Mode:      record[indices["Mode"]],
			// Set default values
			AddrCount: 0,
			WordSize:  0,
			Status: "Configured",
		}
		cards = append(cards, card)
	}

	return cards, nil
}



func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "Welcome!\n")
}

func main() {

  router := httprouter.New()
	io_ctr := io_class.IoV1{}
	io_ctr.Inventory, _ = ReadCardsFromCSV("config.txt")
	io_ctr.Init()
	api := api_class.NewAPI(&io_ctr)

	// Index Page
  router.GET("/",Index)

	// static files
  router.ServeFiles("/debug/*filepath", http.Dir("./public"))
	
  // Websockets
  router.GET("/api/v1/live", api.WsConnect)
  
  fmt.Printf("Starting server at port 8000\n")
  if err := http.ListenAndServe(":8000", router); err != nil {
    log.Fatal(err)
  }
}
