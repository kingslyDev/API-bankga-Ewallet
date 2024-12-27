package utils

import (
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
