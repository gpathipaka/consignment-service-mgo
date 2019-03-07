package main

import (
	pb "consignment-service-mgo/proto/consignment"

	"gopkg.in/mgo.v2"
)

const (
	dbName                = "shipping"
	consignmentCollection = "consignment"
)

//Repository is interface
type Repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
	Close()
}

//ConsignmentRepository is
type ConsignmentRepository struct {
	session *mgo.Session
}

func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(consignmentCollection)
}

//Close kills server cursor
func (repo *ConsignmentRepository) Close() {
	repo.session.Close()
}

// Create is
func (repo *ConsignmentRepository) Create(con *pb.Consignment) error {
	return repo.collection().Insert(con)
}

//GetAll returns eveything.
func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	var consignments []*pb.Consignment
	err := repo.collection().Find(nil).All(&consignments)
	return consignments, err
}
