package main

import (
	"os"

	"github.com/hsyjkjkl/web/service"
	flag "github.com/spf13/pflag"
)

//PORT : the url port
const (
	PORT string = "8000"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}

	server := service.NewServer()
	server.Run(":" + port)
}
