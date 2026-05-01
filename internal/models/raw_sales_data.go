package models

import "time"

// RawSalesData represents commercial transaction documents from the Sales System API.
// This model captures the raw input before it is transformed into the Unified Vault.
type RawSalesData struct {
	ID           string    `json:"id"`            // Unique identifier from Sales System
	VehicleVIN   string    `json:"vin"`           // 17-character Vehicle Identification Number
	DocumentType string    `json:"document_type"` // e.g., "Sales Contract", "Invoice"
	SalesPerson  string    `json:"sales_person"`   // The staff member who handled the sale
	CreatedAt    time.Time `json:"created_at"`    // Timestamp of document issuance
	FileURL      string    `json:"file_url"`      // Path to the original source document
}
