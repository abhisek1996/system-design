package cloudagnosticfileprocessor

import "fmt"

type cloudStorage interface {
	store(file string) error
}

type cloudQueue interface {
	send(message string) error
}

type cloudCompute interface {
	compute(file string) error
}

type cloudProvider interface {
	GetNewStorage(cloudType string) cloudStorage
	GetNewQueue(cloudType string) cloudQueue
	GetNewCompute(cloudType string) cloudCompute
}

// abstract factory
func NewCloudProviderFactory(cloudType string) cloudProvider {
	switch cloudType {
	case "aws":
		return &awsProviderFactory{}
	case "gcp":
		return &gcpProviderFactory{}
	default:
		fmt.Errorf("invalid cloud type")
	}
	return nil
}
