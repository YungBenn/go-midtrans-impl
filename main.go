package main

import (
	"context"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/example"
	"github.com/midtrans/midtrans-go/snap"
)

var s snap.Client

func setupGlobalMidtransConfig() {
	midtrans.ServerKey = ""
	midtrans.Environment = midtrans.Sandbox

	// Optional : here is how if you want to set append payment notification globally
	midtrans.SetPaymentAppendNotification("https://example.com/append")
	// Optional : here is how if you want to set override payment notification globally
	midtrans.SetPaymentOverrideNotification("https://example.com/override")

	//// remove the comment bellow, in cases you need to change the default for Log Level
	// midtrans.DefaultLoggerLevel = &midtrans.LoggerImplementation{
	//	 LogLevel: midtrans.LogInfo,
	// }
}

func initializeSnapClient() {
	s.New("SB-Mid-server-9aBp709Y6zK6ZMuhKSBoHniu", midtrans.Sandbox)
}

func createTransactionWithGlobalConfig() {
	res, err := snap.CreateTransactionWithMap(example.SnapParamWithMap())
	if err != nil {
		fmt.Println("Snap Request Error", err.GetMessage())
	}
	fmt.Println("Snap response", res)
}

func createTransaction(item ItemDetail) {
	// Optional : here is how if you want to set append payment notification for this request
	s.Options.SetPaymentAppendNotification("https://example.com/append")

	// Optional : here is how if you want to set override payment notification for this request
	s.Options.SetPaymentOverrideNotification("https://example.com/override")
	// Send request to Midtrans Snap API

	resp, err := s.CreateTransaction(GenerateSnapReq(item))
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	fmt.Println("Response : ", resp)
}

func createTokenTransactionWithGateway(item ItemDetail) {
	s.Options.SetPaymentOverrideNotification("https://example.com/url2")

	resp, err := s.CreateTransactionToken(GenerateSnapReq(item))
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	fmt.Println("Response : ", resp)
}

func createUrlTransactionWithGateway(item ItemDetail) {
	s.Options.SetContext(context.Background())

	resp, err := s.CreateTransactionUrl(GenerateSnapReq(item))
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	fmt.Println("Response : ", resp)
}

func main() {
	item := ItemDetail{
		Name:  "Futsal bersama",
		Price: 5000,
	}

	fmt.Println("================ Request with global config ================")
	setupGlobalMidtransConfig()
	createTransactionWithGlobalConfig()

	fmt.Println("================ Request with Snap Client ================")
	initializeSnapClient()
	createTransaction(item)

	fmt.Println("================ Request Snap token ================")
	createTokenTransactionWithGateway(item)

	fmt.Println("================ Request Snap URL ================")
	createUrlTransactionWithGateway(item)
}

type ItemDetail struct {
	Name  string
	Price int64
}

func GenerateSnapReq(item ItemDetail) *snap.Request {

	// Initiate Customer address
	custAddress := &midtrans.CustomerAddress{
		FName:       "John",
		LName:       "Doe",
		Phone:       "081234567890",
		Address:     "Baker Street 97th",
		City:        "Jakarta",
		Postcode:    "16000",
		CountryCode: "IDN",
	}

	// Initiate Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "MID-GO-ID-" + example.Random(),
			GrossAmt: item.Price,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    "John",
			LName:    "Doe",
			Email:    "john@doe.com",
			Phone:    "081234567890",
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeGopay,
			snap.PaymentTypeShopeepay,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "ITEM1",
				Price: item.Price,
				Qty:   1,
				Name:  item.Name,
			},
		},
	}
	return snapReq
}
