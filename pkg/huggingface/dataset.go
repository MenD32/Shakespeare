package huggingface

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"

	"github.com/xitongsys/parquet-go-source/http"
	"github.com/xitongsys/parquet-go/reader"
)

type Dataset struct {
	ID     string `json:"_id"`
	RepoID string `json:"id"`
}

type DatasetData []map[string]string

func newDatasetFromReader(body io.Reader) (*Dataset, error) {
	var dataset Dataset
	err := json.NewDecoder(body).Decode(&dataset)
	if err != nil {
		return nil, err
	}
	return &dataset, nil
}

type datasetParquetFiles map[string]map[string]string

func newDatasetParquetFilesFromReader(body io.Reader) (datasetParquetFiles, error) {
	var parquetFiles datasetParquetFiles
	err := json.NewDecoder(body).Decode(&parquetFiles)
	if err != nil {
		return nil, err
	}
	return parquetFiles, nil
}

func NewDatasetData(url string) (DatasetData, error) {
	httpReader, err := http.NewHttpReader(
		url,
		false,
		false,
		map[string]string{},
	)

	if err != nil {
		return nil, err
	}
	parquetReader, err := reader.NewParquetReader(httpReader, nil, 4)
	if err != nil {
		return nil, err
	}

	num := int(parquetReader.GetNumRows())
	res, err := parquetReader.ReadByNumber(num)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]string, 0)
	jsonBs, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBs, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetRandomPrompt(data DatasetData, column string, filter func(string) bool) (string, error) {
	prompts := make([]string, 0)

	for _, row := range data {
		if filter(row[column]) {
			prompts = append(prompts, row[column])
		}
	}

	len := len(prompts)
	if len == 0 {
		return "", fmt.Errorf("no compliant prompts found")
	}
	return prompts[rand.Intn(len)], nil
}
