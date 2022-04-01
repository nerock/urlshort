package router

import (
	"context"

	"github.com/nerock/urlshort/grpc/proto"
	"google.golang.org/grpc"
)

type URLgRPC struct {
	proto.UnimplementedUrlShortenerServer
	svc URLService
}

func NewURLgRPC(svc URLService) *URLgRPC {
	return &URLgRPC{
		svc: svc,
	}
}

func (u *URLgRPC) Register(srv *grpc.Server) {
	proto.RegisterUrlShortenerServer(srv, u)
}

func (u URLgRPC) CreateURL(ctx context.Context, request *proto.CreateURLRequest) (*proto.URLResponse, error) {
	shortUrl, err := u.svc.CreateURL(ctx, request.Url)
	if err != nil {
		return nil, err
	}

	return &proto.URLResponse{
		Url:      request.Url,
		ShortUrl: shortUrl,
	}, nil
}

func (u URLgRPC) GetURL(ctx context.Context, request *proto.URLRequest) (*proto.URLResponse, error) {
	longURL, shortURL, err := u.svc.GetURL(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return &proto.URLResponse{
		Url:      longURL,
		ShortUrl: shortURL,
	}, nil
}

func (u URLgRPC) DeleteURL(ctx context.Context, request *proto.URLRequest) (*proto.DeleteURLResponse, error) {
	err := u.svc.DeleteURL(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteURLResponse{
		Ok: true,
	}, nil
}

func (u URLgRPC) GetRedirectionCount(ctx context.Context, request *proto.URLRequest) (*proto.RedirectionCountResponse, error) {
	count, err := u.svc.GetRedirectionCount(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return &proto.RedirectionCountResponse{
		Id:    request.Id,
		Count: int32(count),
	}, nil
}
