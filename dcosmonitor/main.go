package main

import (
	"net/http"
	"os"
	"os/signal"
	// "fmt"
	// "time"

	"github.com/Sirupsen/logrus"

	"github.com/justinrlee/dcos-tools/dcosclient"
)

func setupLogger () {
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	logrus.SetFormatter(Formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

// I hate this.
var client *dcosclient.Client

// func init() {
// }

func main () {
	setupLogger()
	logrus.Infoln("Status monitor starting.")

	// To detect 
	signalChan := make(chan os.Signal, 1)

	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			logrus.Fatalln("Received an interrupt, stopping")
			// ManagerRemoveAllScalers()
			cleanupDone <- true
		}
	}()

	// client := dcos-Client.

	client = new(dcosclient.Client)

	client.Host = "34.221.25.119"
	client.UserAgent = "dcosmonitor"
	
	client.Setup()

	// Used for testing
	// body, err := client.Request("GET", 443, "ca/dcos-ca.crt")
	// if err != nil {
	// 	fmt.Println("Something broke")
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(body)
	// }

	router := NewRouter()
	logrus.Fatal(http.ListenAndServe(":8083", router))
}