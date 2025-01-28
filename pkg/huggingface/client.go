package huggingface

import (
	"net/http"
	"path"
)

const (
	// DefaultURL is the default URL for the Hugging Face API.
	DefaultURL = "https://huggingface.co/api"

	// Endpoints
	datasetEndpoint = "datasets"
	parquetEndpoint = "parquet"
)

type Client struct {
	URL   string
	Token string
}

func NewDefaultClient() *Client {
	return &Client{
		URL:   DefaultURL,
		Token: "",
	}
}

func (c *Client) getHTTPClient() *http.Client {
	return &http.Client{}
}

func (c *Client) generateGetDatasetRequest(RepoID string) (*http.Request, error) {
	endpoint := path.Join(c.URL, datasetEndpoint, RepoID)
	return http.NewRequest("GET", endpoint, nil)
}

func (c *Client) generateGetDatasetParquetRequest(RepoID string) (*http.Request, error) {
	endpoint := path.Join(c.URL, datasetEndpoint, RepoID, parquetEndpoint)
	return http.NewRequest("GET", endpoint, nil)
}

func (c *Client) GetDataset(RepoID string) (*Dataset, error) {
	client := c.getHTTPClient()
	req, err := c.generateGetDatasetRequest(RepoID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return newDatasetFromReader(resp.Body)
}

func (c *Client) GetParquetFiles(RepoID string) (datasetParquetFiles, error) {
	client := c.getHTTPClient()
	req, err := c.generateGetDatasetParquetRequest(RepoID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return newDatasetParquetFilesFromReader(resp.Body)
}
