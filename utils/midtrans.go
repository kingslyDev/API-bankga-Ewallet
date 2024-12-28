package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/kingslyDev/API-bankga-Ewallet/models"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// CallMidtrans - Fungsi untuk memanggil API Midtrans Snap
func CallMidtrans(req *snap.Request) (*snap.Response, error) {
	client := snap.Client{}
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	client.New(serverKey, midtrans.Sandbox)
	resp, err := client.CreateTransaction(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func InitMidtransClient() *snap.Client {
    serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
    isProduction := os.Getenv("MIDTRANS_IS_PRODUCTION") == "true"
    client := snap.Client{}
    if isProduction {
        client.New(serverKey, midtrans.Production)
    } else {
        client.New(serverKey, midtrans.Sandbox)
    }
    return &client
}

// BuildMidtransParams - Fungsi untuk membangun parameter untuk Snap API
func BuildMidtransParams(orderID string, amount float64, user models.User) *snap.Request {
	return &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			LName: "",
			Email: user.Email,
		},
		EnabledPayments: []snap.SnapPaymentType{snap.PaymentTypeGopay}, // Membatasi hanya untuk Gopay
	}
}

// MidtransNotification - Struct untuk notifikasi dari Midtrans
type MidtransNotification struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
// ParseMidtransNotification - Fungsi untuk memparsing notifikasi dari Midtrans
func ParseMidtransNotification(r *http.Request) (*MidtransNotification, error) {
	// Validasi Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, errors.New("invalid content type, expected application/json")
	}

	// Decode payload JSON
	var notif MidtransNotification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		return nil, err
	}

	// Validasi field wajib
	if notif.OrderID == "" || notif.TransactionStatus == "" || notif.PaymentType == "" {
		return nil, errors.New("missing required fields in notification payload")
	}

	return &notif, nil
}
