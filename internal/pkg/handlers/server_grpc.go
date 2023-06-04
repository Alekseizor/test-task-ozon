package handlers

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"test-task-ozon/internal/pkg/repository/links"
	"test-task-ozon/internal/pkg/urlgeneration"
)

type ConverterServer struct {
	LinkRepo links.LinkRepo
	UnimplementedConverterServiceServer
}

func (c ConverterServer) GetLink(_ context.Context, requestLink *RequestGetLink) (*ResponseGetLink, error) {
	link, err := c.LinkRepo.GetInitialLink(requestLink.GetShortenUrl())
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, err
	}
	return &ResponseGetLink{
		InitialUrl: link.InitialURL,
	}, nil
}

func (c ConverterServer) Generation(_ context.Context, requestGeneration *RequestGeneration) (*ResponseGeneration, error) {
	link := new(links.Links)
	link.InitialURL = requestGeneration.GetInitialUrl()
	existingURL, err := c.LinkRepo.GetShortenLink(link.InitialURL)
	if err != nil {
		return nil, err
	}
	if existingURL != nil {
		link.ShortenURL = existingURL.ShortenURL
	} else {
		link.ShortenURL = urlgeneration.GenerationURL()
		err = c.LinkRepo.AddLink(link)
		if err != nil {
			return nil, err
		}
	}
	finalLink := prefixURL + link.ShortenURL
	return &ResponseGeneration{
		ShortenUrl: finalLink,
	}, nil
}

func StartConverterServer(linkRepo links.LinkRepo) {
	lis, err := net.Listen("tcp", "localhost:9879")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterConverterServiceServer(grpcServer, &ConverterServer{
		LinkRepo: linkRepo,
	})
	grpcServer.Serve(lis)
}
