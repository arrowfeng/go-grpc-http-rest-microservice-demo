package main

import (
	"fmt"
	"os"

	cmd "github.com/arrowfeng/go-grpc-http-rest-microservice-demo/pkg/cmd/rest"
)

func main() {
	if err := cmd.RunServerRest(); err != nil {
		fmt.Fprintf(os.Stderr, "%v/n", err)
		os.Exit(1)
	}
}
