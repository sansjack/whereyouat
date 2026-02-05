package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"whereyouat/internal/env"
	mmdbfetch "whereyouat/internal/mmdb-fetch"
	"whereyouat/internal/rpc"
	"whereyouat/internal/rpc/services"

	"github.com/oschwald/maxminddb-golang/v2"
)

func main() {
	log.Println("RPC servers starting...")

	if err := env.Load(); err != nil {
		log.Printf("Warning: failed to load .env file: %v (using defaults)", err)
	}

	cfg := env.Get()

	fetcher := mmdbfetch.New()
	dbPath, err := fetcher.EnsureDatabase()

	if err != nil {
		log.Fatalf("Failed to ensure database: %v", err)
	}

	db, err := maxminddb.Open(dbPath)
	if err != nil {
		log.Fatalf("IP Database failed to load: %v", err)
	}

	log.Println("IP Database loaded successfully.")

	locationService := &services.LocationService{IpDb: db}

	tcpServer := rpc.NewServer(cfg.TCPAddress, locationService)
	httpServer := rpc.NewHTTPServer(cfg.HTTPAddress, locationService)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down servers...")
		if err := tcpServer.Stop(); err != nil {
			log.Printf("Error stopping TCP server: %v", err)
		}
		if err := httpServer.Stop(); err != nil {
			log.Printf("Error stopping HTTP server: %v", err)
		}
		os.Exit(0)
	}()

	go func() {
		if err := tcpServer.Start(); err != nil {
			log.Fatalf("Failed to start TCP RPC server: %v", err)
		}
	}()

	if err := httpServer.Start(); err != nil {
		log.Fatalf("Failed to start HTTP RPC server: %v", err)
	}
}
