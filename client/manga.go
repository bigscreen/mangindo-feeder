package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/domain"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/gojektech/heimdall"
)

type MangaClient interface {
	GetMangaList() (*domain.MangaListResponse, error)
}

type mangaClient struct {
	httpClient heimdall.Client
}

func buildMangaListEndpoint() string {
	return config.BaseURL() + "/official/2016/main.php"
}

func (c *mangaClient) GetMangaList() (*domain.MangaListResponse, error) {
	res, err := c.httpClient.Get(buildMangaListEndpoint(), nil)
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
