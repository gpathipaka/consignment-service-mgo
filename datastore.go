package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

//CreateSession creates a DB session.
func CreateSession(host string) (*mgo.Session, error) {
	session, err := mgo.Dial(host)
	if err != nil {
		log.Println("Failed to Create the session and the host ", host)
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	return session, nil
}
