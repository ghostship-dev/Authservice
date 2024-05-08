package database

import (
	"context"
	"github.com/edgedb/edgedb-go"
)

func ConnectToEdgeDB() *edgedb.Client {
	ctx := context.Background()
	client, err := edgedb.CreateClient(ctx, edgedb.Options{})
	if err != nil {
		panic(err)
	}
	return client
}

var Client = ConnectToEdgeDB()
var Context = context.Background()
