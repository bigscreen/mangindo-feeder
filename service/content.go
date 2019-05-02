package service

import (
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"strings"
)

type ContentService interface {
	GetContents(req contract.ContentRequest) (*[]contract.Content, error)
}

type contentService struct {
	cClient client.ContentClient
}

func getEncodedUrl(url string) string {
	return strings.Replace(url, " ", "%20", -1)
}

func isAdsContentUrl(url string) bool {
	for _, tag := range config.AdsContentTags() {
		if strings.Contains(url, tag) {
			return true
		}
	}
	return false
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
		if !isAdsContentUrl(dc.ImageURL) {
			content := contract.Content{ImageURL: getEncodedUrl(dc.ImageURL)}
			contents = append(contents, content)
		}
	}

	if contents == nil {
		return nil, mErr.NewNotFoundError("content")
	}

	return &contents, nil
}

func NewContentService(cClient client.ContentClient) *contentService {
	return &contentService{cClient: cClient}
}
