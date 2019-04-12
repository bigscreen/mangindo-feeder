package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/domain"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/gojektech/heimdall"
	"io/ioutil"
	"time"
)

type MangaClient interface {
	GetMangaList(ctx context.Context) (*domain.MangaListResponse, error)
}

type mangaClient struct {
	httpClient heimdall.Client
}

func buildMangaListEndpoint() string {
	return config.BaseURL() + "/official/2016/main.php"
}

func (c *mangaClient) GetMangaList(ctx context.Context) (*domain.MangaListResponse, error) {
	res, err := c.httpClient.Get(buildMangaListEndpoint(), nil)
	if err != nil {
		errMsg := constants.ServerError + " " + err.Error()
		return nil, errors.New(errMsg)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response *domain.MangaListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Errorf("Error when unmarshalling origin response: %s", err.Error())
		return nil, errors.New(constants.InvalidJSONResponseError)
	}
	return response, err
}

func NewMangaClient() *mangaClient {
	hc := config.HystrixConfig()
	timeout := time.Duration(hc.Timeout) * time.Millisecond

	httpClient := heimdall.NewHystrixHTTPClient(timeout, heimdall.NewHystrixConfig(constants.GetMangaListCommand, hc))
	return &mangaClient{
		httpClient: httpClient,
	}
}
