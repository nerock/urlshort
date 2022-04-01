package client

import (
	"context"
	"fmt"

	"github.com/nerock/urlshort/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// URLClient is a client to use the url shortener via gRPC
type URLClient struct {
	conn   *grpc.ClientConn
	client proto.UrlShortenerClient
}

// NewURLClient creates a new URLClient
func NewURLClient(ctx context.Context, url string) (URLClient, error) {
	conn, err := grpc.DialContext(ctx, url, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return URLClient{}, fmt.Errorf("create gRPC conn: %w", err)
	}

	return URLClient{
		conn:   conn,
		client: proto.NewUrlShortenerClient(conn),
	}, nil
}

// Close closes the URLClient connection
func (u URLClient) Close() error {
	return u.conn.Close()
}

// CreateURL sends a request to create a new shortened url
func (u URLClient) CreateURL(ctx context.Context, url string) (string, string, error) {
	res, err := u.client.CreateURL(ctx, &proto.CreateURLRequest{Url: url})
	if err != nil {
		return "", "", fmt.Errorf("could not create url: %w", err)
	}

	return res.Url, res.ShortUrl, nil
}

// GetURL sends a request to get a shortened url by its id
func (u URLClient) GetURL(ctx context.Context, id string) (string, string, error) {
	res, err := u.client.GetURL(ctx, &proto.URLRequest{Id: id})
	if err != nil {
		return "", "", fmt.Errorf("could not get url: %w", err)
	}

	return res.Url, res.ShortUrl, nil
}

// GetURL sends a request to delete a shortened url by its id
func (u URLClient) DeleteURL(ctx context.Context, id string) error {
	res, err := u.client.DeleteURL(ctx, &proto.URLRequest{Id: id})
	if err != nil {
		return fmt.Errorf("could not delete url: %w", err)
	}

	if !res.Ok {
		return fmt.Errorf("could not delete url")
	}

	return nil
}

// GetRedirectionCount sends a request to get a shortened url redirection count
func (u URLClient) GetRedirectionCount(ctx context.Context, id string) (string, int, error) {
	res, err := u.client.GetRedirectionCount(ctx, &proto.URLRequest{Id: id})
	if err != nil {
		return "", 0, fmt.Errorf("could not create url: %w", err)
	}

	return res.Id, int(res.Count), nil
}
