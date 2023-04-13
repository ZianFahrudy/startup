package payment

import (
	"bwastartup/user"
	"fmt"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var sl snap.Client

type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	sl.New("SB-Mid-server-aZHMJQxWrIDZRMQaJhb2ogVr", midtrans.Sandbox)

	resp, err := sl.CreateTransactionUrl(GenerateSnapReq(transaction, user))
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	fmt.Println("Response : ", resp)
	return resp, err
}

func GenerateSnapReq(transaction Transaction, user user.User) *snap.Request {

	// Initiate Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "#ORDER-" + strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},

		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
		EnabledPayments: snap.AllSnapPaymentType,
	}
	return snapReq
}
