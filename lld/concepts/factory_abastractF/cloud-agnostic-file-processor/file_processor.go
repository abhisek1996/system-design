package cloudagnosticfileprocessor

import "fmt"

type GCPStorage struct {
}

func (g *GCPStorage) store(file string) error {
	fmt.Println("Storing file in GCP")
	return nil
}

func GetNewStorage(cloudType string) cloudStorage {
	switch cloudType {
	case "aws":
		return &awsStorage{}
	case "gcp":
		return &GCPStorage{}
	default:
		fmt.Errorf("invalid cloud type")
	}
	return nil
}

func FileProcessor() {

	cloudType := "aws"

	cloudProvider := NewCloudProviderFactory(cloudType)
	storage := cloudProvider.GetNewStorage(cloudType)
	queue := cloudProvider.GetNewQueue(cloudType)
	compute := cloudProvider.GetNewCompute(cloudType)

	storage.store("file.txt")
	queue.send("file.txt")
	compute.compute("file.txt")
}
