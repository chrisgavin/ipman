package edgeosclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const sessionRefreshTime = 10 * time.Minute

type EdgeOSClient struct {
	address          string
	username         string
	password         string
	httpClientCache  *http.Client
	sessionRefreshed time.Time
}

func NewEdgeOSClient(address, username, password string) *EdgeOSClient {
	return &EdgeOSClient{
		address:  address,
		username: username,
		password: password,
	}
}

func (client *EdgeOSClient) httpClient() (*http.Client, error) {
	if client.httpClientCache != nil && time.Since(client.sessionRefreshed) < sessionRefreshTime {
		return client.httpClientCache, nil
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.httpClientCache = &http.Client{
		Jar:       jar,
		Transport: &csrfTransport{RoundTripper: http.DefaultTransport},
	}

	loginValues := url.Values{
		"username": []string{client.username},
		"password": []string{client.password},
	}
	res, err := client.httpClientCache.PostForm(client.address, loginValues)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to login.")
	}
	defer res.Body.Close()

	client.sessionRefreshed = time.Now()
	return client.httpClientCache, err
}

func (client *EdgeOSClient) Get() (*ConfigurationRoot, error) {
	httpClient, err := client.httpClient()
	if err != nil {
		return nil, err
	}

	response, err := httpClient.Get(client.address + "/api/edge/get.json")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get configuration.")
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, errors.Errorf("Failed to get configuration. Response code was %d.", response.StatusCode)
	}

	getResponse := GetResponse{}
	if err := json.NewDecoder(response.Body).Decode(&getResponse); err != nil {
		return nil, errors.Wrap(err, "Failed to decode response.")
	}

	return &getResponse.Configuration, nil
}

func (client *EdgeOSClient) Set(updates interface{}) error {
	httpClient, err := client.httpClient()
	if err != nil {
		return err
	}

	body, err := json.Marshal(updates)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal updates.")
	}
	response, err := httpClient.Post(client.address+"/api/edge/set.json", "application/json", bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "Failed to update configuration.")
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return errors.Errorf("Failed to update configuration. Response code was %d.", response.StatusCode)
	}

	return nil
}

func (client *EdgeOSClient) Delete(updates interface{}) error {
	httpClient, err := client.httpClient()
	if err != nil {
		return err
	}

	body, err := json.Marshal(updates)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal delete.")
	}

	response, err := httpClient.Post(client.address+"/api/edge/delete.json", "application/json", bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "Failed to delete configuration.")
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return errors.Errorf("Failed to delete configuration. Response code was %d.", response.StatusCode)
	}

	return nil
}
