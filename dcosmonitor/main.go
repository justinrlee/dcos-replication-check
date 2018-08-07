package main

import (
	"net/http"
	"os"
	"os/signal"
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
	client := new(dcosclient.Client)

	client.Host = "localhost"
	req, err := client.NewRequest("GET", 443, "ca/dcos-ca.crt", nil)
	if err != nil {
		fmt.Println(req)
	}

	router := NewRouter()
	logrus.Fatal(http.ListenAndServe(":8083", router))
}