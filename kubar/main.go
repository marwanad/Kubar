package main

import (
	"flag"
	"log"
	"os"

	"github.com/marwanad/kubar/export"
	"github.com/marwanad/kubar/restore"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%v", err)
	}
	path := flag.String("path", dir, "path used for exporting/restoring Kubernetes resources")
	mode := flag.String("mode", "", "mode used by Kubar. accepted values: export and restore")

	flag.Parse()
	switch *mode {
	case "export":
		if err := export.Export(*path); err != nil {
			log.Fatalf("%v", err)
		}
	case "restore":
		if err := restore.Restore(*path); err != nil {
			log.Fatalf("%v", err)
		}
	default:
		log.Fatalf("%v", "Couldn't recognize the mode provided")
	}
}
