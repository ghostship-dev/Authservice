package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/ghostship-dev/authservice/core"
	"github.com/ghostship-dev/authservice/core/config"
)

func main() {
	hostname := flag.String("hostname", "localhost", "Hostname to listen on")
	port := flag.Int("port", 8080, "Port to listen on")
	dbEngine := flag.String("dbengine", "edgedb", "Database engine")
	edgeDBInstance := flag.String("edgedb_instance", "", "EdgeDB instance name")
	databaseDSN := flag.String("database_dsn", "", "Database DSN")

	flag.Parse()

	config := config.Config{}

	if *dbEngine != "" {
		config.SetDBEngine(*dbEngine)
	} else {
		config.SetHostname(os.Getenv("DATABASE_ENGINE"))
	}

	if *port != 0 {
		config.SetPort(*port)
	} else {
		envPortInt, _ := strconv.Atoi(os.Getenv("PORT"))
		config.SetPort(envPortInt)
	}

	if *hostname != "" {
		config.SetHostname(*hostname)
	} else {
		config.SetHostname(os.Getenv("HOSTNAME"))
	}

	if *edgeDBInstance != "" {
		config.Database.SetEdgeDBInstanceName(*edgeDBInstance)
		os.Setenv("EDGEDB_INSTANCE", *edgeDBInstance)
	} else {
		config.Database.SetEdgeDBInstanceName(os.Getenv("EDGEDB_INSTANCE"))
	}

	if *databaseDSN != "" {
		config.Database.SetDSN(*databaseDSN)
	} else {
		config.Database.SetDSN(os.Getenv("DATABASE_DSN"))
	}

	core.RunService(&config)
}
