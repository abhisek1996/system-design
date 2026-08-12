package cloudagnosticfileprocessor

import "fmt"

type awsStorage struct {
}

func (a *awsStorage) store(file string) error {
	fmt.Println("Storing file in AWS S3")
	return nil
}

type awsQueue struct {
}

func (a *awsQueue) send(message string) error {
	fmt.Println("Sending message to AWS SQS")
	return nil
}

type awsCompute struct {
}

func (a *awsCompute) compute(file string) error {
	fmt.Println("Computing file in AWS Lambda")
	return nil
}

// concrete factory
type awsProviderFactory struct {
}

func (a *awsProviderFactory) GetNewStorage(cloudType string) cloudStorage {
	return &awsStorage{}
}

func (a *awsProviderFactory) GetNewQueue(cloudType string) cloudQueue {
	return &awsQueue{}
}

func (a *awsProviderFactory) GetNewCompute(cloudType string) cloudCompute {
	return &awsCompute{}
}
