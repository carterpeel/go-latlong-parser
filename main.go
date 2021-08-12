package main

import (
	"./glp"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2/json"
	"os"
)

var (
	lat   = os.Getenv("lat")
	long  = os.Getenv("long")
	token = os.Getenv("apikey")
)

type ResponseData struct {
	FormattedAddress string `json:"formatted_address,omitempty"`
}

func main() {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "",
	})
	l.SetOutput(os.Stdout)
	switch {
	case lat == "":
		l.Logf(logrus.FatalLevel, "'lat' env field cannot be empty")
		return
	case long == "":
		l.Logf(logrus.FatalLevel, "'long' env field cannot be empty")
		return
	case token == "":
		l.Logf(logrus.FatalLevel, "'apikey' env field cannot be empty ")
		return
	}
	lData, err := glp.NewLatLong(lat, long, token)
	if err != nil {
		l.Logf(logrus.FatalLevel, "Could not generate new location structure from requested data: %s\n", err)
		return
	}
	addy, err := lData.GetAddress()
	if err != nil {
		l.Logf(logrus.FatalLevel, "Could not retrieve an address with the provided coordinates: %s\n", err)
		return
	}
	jsonData := &ResponseData{
		FormattedAddress: addy,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		l.Logf(logrus.FatalLevel, "Could not parse response data into a JSON struct: %s\n", err)
		return
	}
	fmt.Print(string(jsonBytes))
}


