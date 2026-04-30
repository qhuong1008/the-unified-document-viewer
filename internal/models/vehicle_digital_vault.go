package models

import (
	"time"

	"github.com/google/uuid"
)

// VehicleDigitalVault represents the unified record for the Model C Digital Vault.
type VehicleDigitalVault struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`        // Internal unique PK
	ExternalID   string    `gorm:"index;not null" json:"external_id"`     // Original ID from source system
	VIN          string    `gorm:"index;not null" json:"vin"`             // Vehicle identifier (Red Thread)
	SourceSystem string    `gorm:"type:varchar(20)" json:"source_system"` // Source: SALES or SERVICE
	Title        string    `gorm:"type:varchar(255)" json:"title"`        // Normalized title for UI display
	DocCategory  string    `gorm:"type:varchar(50)" json:"doc_category"`  // Commercial, Technical, or Legal
	EventDate    time.Time `json:"event_date"`                            // Actual date of the event/service
	FileURL    string    `gorm:"type:text" json:"access_url"`           // Direct link to source document
	SyncedAt     time.Time `gorm:"autoCreateTime" json:"synced_at"`       // Internal synchronization timestamp
}
