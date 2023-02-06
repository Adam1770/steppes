package main

import (
	"errors"
	"log"
	"os"
	"steppes/client"
	"steppes/router"

	"github.com/gin-gonic/gin"
)

const (
	serverAddress = "localhost:90"
	certFilePath  = "./servercert.pem"
	keyFilePath   = "./serverkey.pem"
)

var ErrCertificationFiles = errors.New("one of the certification files could not be opened")

func main() {
	if err := checkCertFiles(); err != nil {
		log.Fatal(err)
	}
	awsClient := client.NewAwsClient()
	r := router.NewRouter(awsClient)
	gin.SetMode(gin.ReleaseMode)
	r.Engine.RunTLS(serverAddress, certFilePath, keyFilePath)

}

func checkCertFiles() error {
	if file, err := os.Open(certFilePath); err != nil {
		return ErrCertificationFiles
	} else {
		file.Close()
	}
	if file, err := os.Open(keyFilePath); err != nil {
		return ErrCertificationFiles
	} else {
		file.Close()
	}
	return nil
}

//radek passy: jeezy gz1f=R@mM
