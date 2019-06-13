package service

import (
	"github.com/bigscreen/mangindo-feeder/cache/manager"
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	"strings"
)

type ContentService interface {
	GetContents(req contract.ContentRequest) (*[]contract.Content, error)
}

type contentService struct {
	contentClient       client.ContentClient
	contentCacheManager manager.ContentCacheManager
	workerService       WorkerService
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
	cl, err := s.contentCacheManager.GetCache(req.TitleId, req.Chapter)
	if err != nil {
		cl, err = s.contentClient.GetContentList(req.TitleId, req.Chapter)
		if err != nil {
			return nil, mErr.NewGenericError()
		}

		err = s.workerService.SetContentCache(req.TitleId, req.Chapter)
		if err != nil {
			logger.Errorf("Failed to enqueue %s job, with error: %s", constants.SetContentCacheJob, err.Error())
		}
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

func NewContentService(cc client.ContentClient, ccm manager.ContentCacheManager, ws WorkerService) *contentService {
	return &contentService{
		contentClient:       cc,
		contentCacheManager: ccm,
		workerService:       ws,
	}
}
