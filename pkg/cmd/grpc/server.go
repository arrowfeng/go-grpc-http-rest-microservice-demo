package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	"github.com/arrowfeng/go-grpc-http-rest-microservice-demo/pkg/protocol/grpc"
	v1 "github.com/arrowfeng/go-grpc-http-rest-microservice-demo/pkg/service/v1"

	_ "github.com/go-sql-driver/mysql"
)

// Config is configuration for Server
type Config struct {

	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string

	// DB Datastore parameters section
	// DatastoreDBHost is host of database
	DatastoreDBHost string

	// DatastoreDBPort is port of database
	DatastoreDBPort string

	// DatastoreDBUser is username to connect to database
	DatastoreDBUser string

	// DatastoreDBPassword password to connect to database
	DatastoreDBPassword string

	// DatastoreDBSchema is schema of database
	DatastoreDBSchema string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Datastore host")
	flag.StringVar(&cfg.DatastoreDBPort, "db-port", "", "Datastore port")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Datastore user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Datastore password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Datastore schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// add MySQL driver specific parameter to parse date/time
	// Drop it for another database

	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost+":"+cfg.DatastoreDBPort,
		cfg.DatastoreDBSchema,
		param)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)

}
