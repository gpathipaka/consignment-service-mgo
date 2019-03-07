package main

import (
	"log"
	"os"

	pb "consignment-service-mgo/proto/consignment"
	vesselPb "consignment-service-mgo/proto/vessel"

	"github.com/micro/go-micro"
)

const (
	defaultHost = "localhost:27017"
)

func main() {
	log.Println("Server starting....")
	//Get the DB host from the environment variable.
	host := os.Getenv("DB_HOST")
	if host == "" {
		log.Println("DB Host is empty and setting host to default...", defaultHost)
		host = defaultHost
	}

	session, err := CreateSession(host)
	// Mgo creates a 'master' session, we need to end that session
	// before the main function closes.
	if err != nil {
		// wrap the error from create session
		log.Panic("Could not connect to the Data Store with the host %s - %v", host, err)
	}
	defer session.Clone()

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	vesselClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	pb.RegisterShippingServiceHandler(srv.Server(), &service{session, vesselClient})
	srv.Init()

	if err := srv.Run(); err != nil {
		log.Printf("Could not run the server %v", err)
	}

	log.Println("Server about to go down.......")
}
