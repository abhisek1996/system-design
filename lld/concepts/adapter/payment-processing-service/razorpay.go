package paymentprocessingservice

import "fmt"

// RazorpayGateway

type razorpayGateway struct {
	client *razorpayClient
}

func NewRazorpayGateway() paymentGateway {
	return &razorpayGateway{
		client: &razorpayClient{},
	}
}

func (r *razorpayGateway) pay(amount int) {
	// convert paymentDetails to payuRequest
	r.client.Pay(amount)
}

type razorpayClient struct {
}

func (r *razorpayClient) Pay(amount int) error {
	// Simulated third-party call
	fmt.Printf("Paid ₹%.2f using Razorpay Client (Third Party)\n", amount)
	return nil
}
