package main

import (
	"flag"
	"fmt"
	"github.com/carterpeel/go-latlong-parser/glp"
	"github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2/json"
	"os"
)

type ResponseData struct {
	FormattedAddress string `json:"formatted_address,omitempty"`
}


func main() {
	lat   := flag.String("lat", "", "Latitude")
	long  := flag.String("long", "", "Longitude")
	token := flag.String("apikey", "", "Your Google Maps API key")
	flag.Parse()
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "",
	})
	l.SetOutput(os.Stdout)
	switch {
	case *lat == "":
		flag.PrintDefaults()
		return
	case *long == "":
		flag.PrintDefaults()
		return
	case *token == "":
		flag.PrintDefaults()
		return
	}
	lData, err := glp.NewLatLong(*lat, *long, *token)
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



