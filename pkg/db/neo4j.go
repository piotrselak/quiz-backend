package db

import (
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitNeo4j() neo4j.DriverWithContext {
	dbUri := "neo4j://localhost:7687"
	driver, err := neo4j.NewDriverWithContext(dbUri,
		neo4j.BasicAuth("username", "password", ""))
	if err != nil {
		panic(err)
	}
	return driver
}

func GetSessionFromContext(r *http.Request) neo4j.SessionWithContext {
	var session neo4j.SessionWithContext
	session = r.Context().Value("session").(neo4j.SessionWithContext)
	return session
}
