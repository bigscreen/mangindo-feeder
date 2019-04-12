package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/domain"
	"github.com/gojektech/heimdall"
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
	_, err := c.httpClient.Get(buildContentListEndpoint(titleId, chapter), nil)
	if err != nil {
		errMsg := constants.ServerError + " " + err.Error()
		return nil, errors.New(errMsg)
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
