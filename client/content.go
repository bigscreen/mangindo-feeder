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

type ContentClient interface {
	GetContentList(ctx context.Context, titleId string, chapter float32) (*domain.ContentListResponse, error)
}

type contentClient struct {
	httpClient heimdall.Client
}

func buildContentListEndpoint(titleId string, chapter float32) string {
	qParams := "?manga=%s&chapter=%f"
	qParams = fmt.Sprintf(qParams, titleId, chapter)
	return config.BaseURL() + "/official/2016/image_list.php" + qParams
}

func (c *contentClient) GetContentList(ctx context.Context, titleId string, chapter float32) (*domain.ContentListResponse, error) {
	res, err := c.httpClient.Get(buildContentListEndpoint(titleId, chapter), nil)
	if err != nil {
		errMsg := constants.ServerError + " " + err.Error()
		return nil, errors.New(errMsg)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response *domain.ContentListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Errorf("Error when unmarshalling origin response: %s", err.Error())
		return nil, errors.New(constants.InvalidJSONResponseError)
	}
	return nil, err
}

func NewContentClient() *contentClient {
	hc := config.HystrixConfig()
	timeout := time.Duration(hc.Timeout) * time.Millisecond

	httpClient := heimdall.NewHystrixHTTPClient(timeout, heimdall.NewHystrixConfig(constants.GetContentListCommand, hc))
	return &contentClient{
		httpClient: httpClient,
	}
}
