package models

import (
	"time"

	"github.com/google/uuid"
)

type VehicleDigitalVault struct {
    ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`                  
    ExternalID   string    `gorm:"uniqueIndex;not null" json:"external_id"`         
    VIN          string    `gorm:"index;not null" json:"vin"`                       
    SourceSystem string    `gorm:"type:varchar(20)" json:"source_system"`           // SALES or SERVICE
    Title        string    `gorm:"type:varchar(255)" json:"title"`                  
    DocCategory  string    `gorm:"type:varchar(50)" json:"doc_category"`            // Commercial, Technical, or Legal
    
    // Domain Specific Fields
    Technician   string    `gorm:"column:technician;type:varchar(50)" json:"technician"`   
    ServiceType  string    `gorm:"column:service_type;type:varchar(50)" json:"service_type"` 
    SalesPerson  string    `gorm:"column:sales_person;type:varchar(50)" json:"sales_person"` 
    
    // DateTime Types
    ServiceCompletionDate  time.Time `gorm:"column:service_completion_date;type:timestamp" json:"service_completion_date"`
    SalesDocumentIssueDate time.Time `gorm:"column:sales_document_issue_date;type:timestamp" json:"sales_document_issue_date"`
    
    FileURL      string    `gorm:"type:text" json:"access_url"`                     
    SyncedAt     time.Time `gorm:"autoCreateTime" json:"synced_at"`                 
}
