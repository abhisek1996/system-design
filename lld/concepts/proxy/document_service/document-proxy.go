package documentservice

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type userKey string

const UserIDKey userKey = "userID"

type proxyDocumentService struct {
	realService DocumentService
}

func NewDocumentProxy(realService DocumentService) DocumentService {
	return &proxyDocumentService{
		realService: realService,
	}
}

func (p *proxyDocumentService) FetchDocument(ctx context.Context, documentID string) ([]byte, error) {
	user, err := p.checkAccess(ctx)
	if err != nil {
		return nil, err
	}

	p.auditLog(ctx, "FETCH", documentID, user)
	return p.realService.FetchDocument(ctx, documentID)
}

func (p *proxyDocumentService) DownloadDocument(ctx context.Context, documentID string) ([]byte, error) {
	user, err := p.checkAccess(ctx)
	if err != nil {
		return nil, err
	}

	p.auditLog(ctx, "DOWNLOAD", documentID, user)
	return p.realService.DownloadDocument(ctx, documentID)
}

func (p *proxyDocumentService) checkAccess(ctx context.Context) (string, error) {
	user, ok := ctx.Value(UserIDKey).(string)
	if !ok || user == "" {
		return "", errors.New("unauthenticated")
	}

	if user != "admin" {
		return "", fmt.Errorf("user %s not authorized", user)
	}
	return user, nil
}

func (p *proxyDocumentService) auditLog(ctx context.Context, action, docID, user string) {
	go func() {
		log.Printf("[AUDIT] User: %s, Action: %s, DocID: %s\n", user, action, docID)
	}()
}

func ProxyDemo() {
	realService := &documentServiceImpl{}
	proxy := NewDocumentProxy(realService)

	// Case 1: Unauthorized
	fmt.Println("--- Case 1: Regular User ---")
	ctxUser := context.WithValue(context.Background(), UserIDKey, "abhishek")
	_, err := proxy.FetchDocument(ctxUser, "doc_123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Case 2: Authorized
	fmt.Println("\n--- Case 2: Admin User ---")
	ctxAdmin := context.WithValue(context.Background(), UserIDKey, "admin")
	data, _ := proxy.FetchDocument(ctxAdmin, "doc_123")
	fmt.Printf("Result: %s\n", string(data))
}
