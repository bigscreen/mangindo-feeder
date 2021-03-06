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

type ContentClient interface {
	GetContentList(titleID string, chapter float32) (*domain.ContentListResponse, error)
}

type contentClient struct {
	httpClient heimdall.Client
}

func buildContentListEndpoint(titleID string, chapter float32) string {
	qParams := "?manga=%s&chapter=%f"
	qParams = fmt.Sprintf(qParams, titleID, chapter)
	return config.BaseURL() + "/official/2016/image_list.php" + qParams
}

func (c *contentClient) GetContentList(titleID string, chapter float32) (*domain.ContentListResponse, error) {
	res, err := c.httpClient.Get(buildContentListEndpoint(titleID, chapter), nil)
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

	var response *domain.ContentListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Errorf("Error when unmarshalling origin response: %s", err.Error())
		return nil, errors.New(constants.InvalidJSONResponseError)
	}
	return response, err
}

func NewContentClient() *contentClient {
	hc := config.HystrixConfig()
	timeout := time.Duration(hc.Timeout) * time.Millisecond

	httpClient := heimdall.NewHystrixHTTPClient(timeout, heimdall.NewHystrixConfig(constants.GetContentListCommand, hc))
	return &contentClient{
		httpClient: httpClient,
	}
}
