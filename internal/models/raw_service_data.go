package models

import "time"

// RawServiceData represents maintenance and repair documents from the Service System API.
// This structure is designed to hold raw logs before they are normalized for the UI.
type RawServiceData struct {
	ID             string    `json:"id"`              // Unique identifier from Service System
	VehicleVIN     string    `json:"vin"`             // VIN used as the key for history lookup
	ServiceType    string    `json:"service_type"`    // e.g., "Oil Change", "Brake Repair"
	Technician     string    `json:"technician"`      // The technician responsible for the work
	CompletionDate time.Time `json:"completion_date"` // Timestamp when the service was finalized
	ReportLink     string    `json:"report_link"`     // Path to the technical inspection report
}
