package service

import (
	"the-unified-document-viewer/internal/adapter"
	"the-unified-document-viewer/internal/models"
)

type DocumentService struct {
	adapter *adapter.SalesServiceAdapter
}

func NewDocumentService(adapter *adapter.SalesServiceAdapter) *DocumentService {
	return &DocumentService{
		adapter: adapter,
	}
}

// ConsolidateData merges sales and service data
func (s *DocumentService) ConsolidateData() (*models.Document, error) {
	salesData, err := s.adapter.FetchSalesData()
	if err != nil {
		return nil, err
	}

	serviceData, err := s.adapter.FetchServiceData()
	if err != nil {
		return nil, err
	}

	// Merge data logic
	document := &models.Document{
		ID:      "consolidated",
		Title:   "Consolidated Document",
		Content: mergeContent(salesData, serviceData),
		Type:    "consolidated",
	}

	return document, nil
}

func mergeContent(sales, service map[string]interface{}) string {
	// TODO: Implement data merging logic
	return "Sales and Service data merged"
}
