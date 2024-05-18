package database

import (
	"context"
	"os"

	"github.com/edgedb/edgedb-go"
)

func ConnectToEdgeDB() *edgedb.Client {
	os.Setenv("EDGEDB_INSTANCE", "Ghostship")
	ctx := context.Background()
	client, err := edgedb.CreateClient(ctx, edgedb.Options{})
	if err != nil {
		panic(err)
	}
	return client
}

var Client = ConnectToEdgeDB()
var Context = context.Background()
