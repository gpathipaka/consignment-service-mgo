package main

import (
	"context"
	"log"

	"gopkg.in/mgo.v2"

	pb "consignment-service-mgo/proto/consignment"
	vesselPb "consignment-service-mgo/proto/vessel"
)

type service struct {
	session      *mgo.Session
	vesselClient vesselPb.VesselServiceClient
}

//GetRepo is
func (s *service) GetRepo() Repository {
	return &ConsignmentRepository{s.session.Clone()}
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	defer s.GetRepo().Close()
	log.Println("Server - Entered the CreateConsignment() with input weight.. ", req.Weight)
	spec := &vesselPb.Specification{
		Capacity:  req.Weight,
		MaxWeight: int32(len(req.Containers)),
	}
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), spec)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VesselId = vesselResponse.Vessel.Id
	if err := s.GetRepo().Create(req); err != nil {
		return err
	}
	res.Created = true
	res.Consignment = req
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	defer s.GetRepo().Close()
	consignments, err := s.GetRepo().GetAll()
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}
