package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/httpapi"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/service"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

func main() {
	address := flag.String("listen", ":8080", "HTTP listen address")
	dataFile := flag.String("data", filepath.Join("data", "yardsense.json"), "persistent data file")
	flag.Parse()
	database, err := store.Open(*dataFile)
	if err != nil {
		log.Fatalf("open persistent store: %v", err)
	}
	application := service.New(database, time.Now)
	server := &http.Server{Addr: *address, Handler: httpapi.New(application), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("YardSense Dispatch listening on %s", *address)
	log.Fatal(server.ListenAndServe())
}
