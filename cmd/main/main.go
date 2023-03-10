package main

import (
	"go-chi-api/cmd/api"
	"os"
)

func main() {
	switch getServerType() {
	case "api":
		api.Run(getServerIdentifier(), getPort())
	default:
		panic("unsupported server type: \"" + getServerType() + "\"")
	}
}

func getServerType() string {
	return os.Getenv("SERVER_TYPE")
}

func getServerIdentifier() string {
	return os.Getenv("SERVER_IDENTIFIER")
}

func getPort() string {
	return os.Getenv("SERVER_PORT")
}
