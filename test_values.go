package braintree

var testCreditCards = map[string]CreditCard{
	"visa":       CreditCard{Number: "4111111111111111"},
	"mastercard": CreditCard{Number: "5555555555554444"},
	"discover":   CreditCard{Number: "6011111111111117"},
}

var testGateway = New(
	Development,
	"integration_merchant_id",
	"8hghpwn86t9zffyy",
	"73ea657092c920d72bcd7dc6d09d103a",
)

// func init() {
// 	testGateway.Logger = log.New(os.Stdout, "\nBT: ", log.LstdFlags)
// }
