package config

import "fmt"

type Config struct {
	Port     int
	Hostname string
	DBEngine string
	Database DatabaseConfig
}

func NewConfig() *Config {
	return &Config{
		Port:     8080,
		Hostname: "localhost",
		DBEngine: "edgedb",
	}
}

func (c *Config) SetPort(port int) {
	if port == 0 {
		c.Port = 8080
		return
	}
	c.Port = port
}

func (c *Config) SetHostname(hostname string) {
	if hostname == "" {
		c.Hostname = "localhost"
		return
	}
	c.Hostname = hostname
}

func (c *Config) SetDBEngine(dbEngine string) {
	if dbEngine != "edgedb" {
		fmt.Println("Invalid database engine specified. Defaulting to edgedb.")
		c.DBEngine = "edgedb"
		return
	}
	c.DBEngine = dbEngine
}

type DatabaseConfig struct {
	Host               string
	Port               int
	DSN                string
	EdgeDBInstanceName string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:               "localhost",
		Port:               5656,
		DSN:                "",
		EdgeDBInstanceName: "",
	}
}

func (c *DatabaseConfig) SetHost(host string) {
	if host == "" {
		c.Host = "localhost"
		return
	}
	c.Host = host
}

func (c *DatabaseConfig) SetPort(port int) {
	if port == 0 {
		c.Port = 5656
		return
	}
	c.Port = port
}

func (c *DatabaseConfig) SetDSN(dsn string) {
	if dsn == "" {
		c.DSN = ""
		return
	}
	c.DSN = dsn
}

func (c *DatabaseConfig) SetEdgeDBInstanceName(instanceName string) {
	c.EdgeDBInstanceName = instanceName
}
