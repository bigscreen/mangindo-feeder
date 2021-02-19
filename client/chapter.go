package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/domain"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/gojektech/heimdall"
)

type ChapterClient interface {
	GetChapterList(titleID string) (*domain.ChapterListResponse, error)
}

type chapterClient struct {
	httpClient heimdall.Client
}

func buildChapterListEndpoint(titleID string) string {
	qParam := "?manga=%s"
	qParam = fmt.Sprintf(qParam, titleID)
	return config.BaseURL() + "/official/2016/chapter_list.php" + qParam
}

func (c *chapterClient) GetChapterList(titleID string) (*domain.ChapterListResponse, error) {
	res, err := c.httpClient.Get(buildChapterListEndpoint(titleID), nil)
	if err != nil {
		errMsg := constants.ServerError + " " + err.Error()
		return nil, errors.New(errMsg)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if string(body) == constants.NullText {
		logger.Error("Origin response body is null")
		return nil, errors.New(constants.InvalidJSONResponseError)
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
