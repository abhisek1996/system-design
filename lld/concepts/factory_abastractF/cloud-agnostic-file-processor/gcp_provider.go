package cloudagnosticfileprocessor

import "fmt"

type gcpProviderFactory struct {
}

type gcpStorage struct {
}

func (a *gcpStorage) store(file string) error {
	fmt.Println("Storing file in AWS S3")
	return nil
}

type gcpQueue struct {
}

func (a *gcpQueue) send(message string) error {
	fmt.Println("Sending message to AWS SQS")
	return nil
}

type gcpCompute struct {
}

func (a *gcpCompute) compute(file string) error {
	fmt.Println("Computing file in AWS Lambda")
	return nil
}

func (a *gcpProviderFactory) GetNewStorage(cloudType string) cloudStorage {
	return &gcpStorage{}
}

func (a *gcpProviderFactory) GetNewQueue(cloudType string) cloudQueue {
	return &gcpQueue{}
}

func (a *gcpProviderFactory) GetNewCompute(cloudType string) cloudCompute {
	return &gcpCompute{}
}
