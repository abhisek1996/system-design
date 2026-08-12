package paymentprocessingservice

import "fmt"

// adapter
type payUGateway struct {
	client *payUClient
}

func NewPayUGateway() paymentGateway {
	return &payUGateway{
		client: &payUClient{},
	}
}

func (p *payUGateway) pay(amount int) {
	// implementation
	// convert paymentDetails to payURequest
	floatAmount := float64(amount)
	p.client.MakePayment(floatAmount)
}

// adaptee
type payUClient struct {
}

func (p *payUClient) MakePayment(amountInRupees float64) error {
	// Simulated third-party call
	fmt.Printf("Paid ₹%.2f using PayU Client (Third Party)\n", amountInRupees)
	return nil
}
