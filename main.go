package main

import (
	"context"
	"errors"
	"log"
	"os"

	pb "consignment-service-mgo/proto/consignment"
	vesselPb "consignment-service-mgo/proto/vessel"

	userService "github.com/gpathipaka/user-service/proto/user"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
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
		micro.WrapHandler(AuthWrapper),
	)
	vesselClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	pb.RegisterShippingServiceHandler(srv.Server(), &service{session, vesselClient})
	srv.Init()

	if err := srv.Run(); err != nil {
		log.Printf("Could not run the server %v", err)
	}

	log.Println("Server about to go down.......")
}

// AuthWrapper is
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		// Note this is now uppercase (not entirely sure why this is...)
		token := meta["Token"]
		log.Println("Authenticating with token: ", token)

		// Auth here
		authClient := userService.NewUserServiceClient("go.micro.srv.user", client.DefaultClient)
		_, err := authClient.ValidateToken(context.Background(), &userService.Token{
			Token: token,
		})
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}
