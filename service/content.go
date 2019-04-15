package service

import (
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
)

type ContentService interface {
	GetContents(req contract.ContentRequest) (*[]contract.Content, error)
}

type contentService struct {
	cClient client.ContentClient
}

func (s *contentService) GetContents(req contract.ContentRequest) (*[]contract.Content, error) {
	cl, err := s.cClient.GetContentList(req.TitleId, req.Chapter)
	if err != nil {
		return nil, mErr.NewGenericError()
	}

	if len(cl.Contents) == 0 {
		return nil, mErr.NewNotFoundError("content")
	}

	var contents []contract.Content
	for _, dc := range cl.Contents {
		content := contract.Content{ImageURL: dc.ImageURL}
		contents = append(contents, content)
	}

	return &contents, nil
}

func NewContentService(cClient client.ContentClient) *contentService {
	return &contentService{cClient: cClient}
}
