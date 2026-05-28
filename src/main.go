package main

import (
	"flag"
	"io"
	"log"
	"os"

	"virga-player/app"
	debugmgr "virga-player/debug/manager"
	"virga-player/version"
)

func main() {
	debugFlag := flag.Bool("debug", false, "enable debug logs and overlay")
	flag.Parse()

	dbg := debugmgr.NewManager(*debugFlag, *debugFlag)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(io.MultiWriter(os.Stderr, dbg.Writer()))
	log.Printf("Virga Player %s", version.AppVersion)

	if err := app.New(app.Options{Debug: *debugFlag}, dbg).Run(); err != nil {
		dbg.Errorf("application error: %v", err)
		log.Fatalf("application error: %v", err)
	}
}
