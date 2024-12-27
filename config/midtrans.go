package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"os"
)

type MidtransClient struct {
	Client snap.Client
}

var Midtrans MidtransClient

func InitMidtrans() {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	isProduction := os.Getenv("MIDTRANS_IS_PRODUCTION") == "true"

	Midtrans.Client.New(serverKey, midtrans.Sandbox)
	if isProduction {
		Midtrans.Client.New(serverKey, midtrans.Production)
	}
}
