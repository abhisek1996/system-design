package paymentprocessingservice

type paymentGateway interface {
	pay(amount int)
}
