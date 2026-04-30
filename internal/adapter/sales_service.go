package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SalesServiceAdapter struct {
	baseURL string
	client  *http.Client
}

func NewSalesServiceAdapter(baseURL string) *SalesServiceAdapter {
	return &SalesServiceAdapter{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// FetchSalesData fetches sales data from API (or mock)
func (s *SalesServiceAdapter) FetchSalesData() (map[string]interface{}, error) {
	resp, err := s.client.Get(fmt.Sprintf("%s/api/sales", s.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// FetchServiceData fetches service data from API (or mock)
func (s *SalesServiceAdapter) FetchServiceData() (map[string]interface{}, error) {
	resp, err := s.client.Get(fmt.Sprintf("%s/api/service", s.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
