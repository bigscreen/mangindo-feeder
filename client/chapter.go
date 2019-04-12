package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/domain"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/gojektech/heimdall"
	"io/ioutil"
	"time"
)

type ChapterClient interface {
	GetChapterList(ctx context.Context, titleId string) (*domain.ChapterListResponse, error)
}

type chapterClient struct {
	httpClient heimdall.Client
}

func buildChapterListEndpoint(titleId string) string {
	qParam := "?manga=%s"
	qParam = fmt.Sprintf(qParam, titleId)
	return config.BaseURL() + "/official/2016/chapter_list.php" + qParam
}

func (c *chapterClient) GetChapterList(ctx context.Context, titleId string) (*domain.ChapterListResponse, error) {
	res, err := c.httpClient.Get(buildChapterListEndpoint(titleId), nil)
	if err != nil {
		errMsg := constants.ServerError + " " + err.Error()
		return nil, errors.New(errMsg)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response *domain.ChapterListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Errorf("Error when unmarshalling origin response: %s", err.Error())
		return nil, errors.New(constants.InvalidJSONResponseError)
	}
	return response, err
}

func NewChapterClient() *chapterClient {
	hc := config.HystrixConfig()
	timeout := time.Duration(hc.Timeout) * time.Millisecond

	httpClient := heimdall.NewHystrixHTTPClient(timeout, heimdall.NewHystrixConfig(constants.GetChapterListCommand, hc))
	return &chapterClient{
		httpClient: httpClient,
	}
}
