package main

import (
	"log"
	"net/rpc"
	"whereyouat/pkg/schemas/location"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Failed to connect to RPC server: %v", err)
	}
	defer client.Close()

	log.Println("Connected to RPC server")

	args := &location.CalculateArgs{}
	reply := &location.CalculateReply{}

	err = client.Call("LocationService.Calculate", args, reply)

	if err != nil {
		log.Fatalf("LocationService.Calculate error: %v", err)
	}

	log.Printf("Location: %s, Code: %s", reply.Location, reply.IsoCode)

	log.Println("All RPC calls completed successfully!")
}
