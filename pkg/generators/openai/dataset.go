package openai

import (
	"encoding/json"
	"net/http"
	"path"
)

const (
	HUGGINGFACE_API_ENDPOINT = "https://huggingface.co/api"

	// Dataset API
	DATASET_API_PATH         = "datasets"
	DATASET_API_PARQUET_PATH = "parquet"
)

type Dataset struct {
	repo_id       string
	parquet_files map[string]map[string]string
}

func NewDataset(repo_id string) *Dataset {
	return &Dataset{
		repo_id: repo_id,
	}
}

func GetDataset(repo_id string) (*Dataset, error) {
	endpoint := path.Join(HUGGINGFACE_API_ENDPOINT, DATASET_API_PATH, repo_id, DATASET_API_PARQUET_PATH)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parquet_files := make(map[string]map[string]string)
	json.NewDecoder(resp.Body).Decode(&parquet_files)

	return &Dataset{
		repo_id:       repo_id,
		parquet_files: parquet_files,
	}, nil
}

func (d *Dataset) RandomPrompt() string {
	return "This is a random prompt"
}
