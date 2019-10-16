package harvest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/imroc/req"
)

const (
	_get   = "GET"
	_patch = "PATCH"
	_post  = "POST"
)

var (
	accountID   = os.Getenv("HARVEST_ACCOUNT_ID")
	accessToken = os.Getenv("HARVEST_ACCESS_TOKEN")
)

var apiVersion = "2"

var (
	endpointHarvest      = "https://api.harvestapp.com/v" + apiVersion
	endpointTimeEntries  = endpointHarvest + "time_entries/"
	endpointStopTimer    = func(entryID string) string { return endpointTimeEntries + entryID + "/stop/" }
	endpointRestartTimer = func(entryID string) string { return endpointTimeEntries + entryID + "/restart/" }
)

func getAuthHeader() req.Header {
	return req.Header{
		"Authorization":      fmt.Sprintf("Bearer %s", accessToken),
		"Harvest-Account-Id": accountID,
	}
}

type params map[string]string

func (p params) toURLValues() url.Values {
	u := url.Values{}
	for k, v := range p {
		u.Set(k, v)
	}
	return u
}

type harvestAPI struct {
	baseURL string
	url     *url.URL
	request *http.Request
	client  *http.Client

	// Added iceKathon
	token        string
	refreshToken string
	accountID    string
}

func NewHarvestAPI(token, accountID string) harvestAPI {
	return harvestAPI{
		client:    http.DefaultClient,
		baseURL:   endpointHarvest,
		token:     token,
		accountID: accountID,
	}
}

func (h *harvestAPI) setURL(path string, p params) {
	h.url, _ = url.Parse(h.baseURL + path)
	h.url.RawQuery = p.toURLValues().Encode()
}

func (h *harvestAPI) setReqHeaderV2() {
	h.request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.token))
	h.request.Header.Set("Harvest-Account-Id", h.accountID)
}

func (h *harvestAPI) sendGetRequest(endpoint string, params params) (resp *http.Response, err error) {
	h.setURL(endpoint, params)

	h.request, err = http.NewRequest(_get, h.url.String(), nil)
	if err != nil {
		return nil, err
	}
	h.setReqHeaderV2()
	return h.client.Do(h.request)
}

func (h *harvestAPI) sendPatchRequest(endpoint string, params params) (resp *http.Response, err error) {
	h.setURL(endpoint, params)

	h.request, err = http.NewRequest(_patch, h.url.String(), nil)
	if err != nil {
		return nil, err
	}
	h.setReqHeaderV2()
	return h.client.Do(h.request)
}

func (h *harvestAPI) sendRequest(method, endpoint string, p params) (resp *http.Response, err error) {
	switch method {
	case _get:
		return h.sendGetRequest(endpoint, p)
	case _patch:
		return h.sendPatchRequest(endpoint, p)
	default:
		return nil, fmt.Errorf("http method not supported yet")
	}
}

func decodeResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}
