package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
  "fmt"
  "net/http"
	"log"
	"flag"
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
	reader.FieldsPerRecord = -1
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
			InterruptPin: "",
		}
//		if(indices["Interrupt"] < len(record)) {
//			card.InterruptPin = record[indices["Interrupt"]]
//		}
		cards = append(cards, card)
	}

	return cards, nil
}



func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "Welcome!\n")
}

func main() {

	dbgDir := flag.String("debug-dir", "", "directory for debug files")
	config := flag.String("config", "/etc/io2websocket-gateway.conf", "config file")
	flag.Parse()
  router := httprouter.New()
	io_ctr := io_class.IoV1{}
	fmt.Printf("read config %s\n", *config)
  var csv_error error
	io_ctr.Inventory, csv_error = ReadCardsFromCSV(*config)
	if(csv_error != nil) {
		fmt.Printf("error reading config file %s: %s\n", *config, csv_error)
		return
	}
	io_ctr.Init()
	api := api_class.NewAPI(&io_ctr)

	// Index Page
  router.GET("/",Index)

	if(*dbgDir != "") {
		// static files
		router.ServeFiles("/debug/*filepath", http.Dir(*dbgDir))
		fmt.Printf("serving /debug UI from %s\n", *dbgDir)
	}
  // Websockets
  router.GET("/api/v1/live", api.WsConnect)
  
  fmt.Printf("Starting server at port 8000\n")
  if err := http.ListenAndServe(":8000", router); err != nil {
    log.Fatal(err)
  }
}
