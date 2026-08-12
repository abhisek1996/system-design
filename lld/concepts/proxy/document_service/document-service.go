package documentservice

import (
	"context"
	"fmt"
)

type DocumentService interface {
	FetchDocument(ctx context.Context, documentID string) ([]byte, error)
	DownloadDocument(ctx context.Context, documentID string) ([]byte, error)
}

type documentServiceImpl struct {
}

func (d *documentServiceImpl) FetchDocument(ctx context.Context, documentID string) ([]byte, error) {
	fmt.Printf("Fetching document %s from storage\n", documentID)
	return []byte("document data"), nil
}

func (d *documentServiceImpl) DownloadDocument(ctx context.Context, documentID string) ([]byte, error) {
	fmt.Printf("Downloading document %s from storage\n", documentID)
	return []byte("binary data"), nil
}
