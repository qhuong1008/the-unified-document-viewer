package repository

import (
	"the-unified-document-viewer/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) UpsertVehicleDigitalVaultRecord(data models.VehicleDigitalVault) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}}, 
		DoUpdates: clause.AssignmentColumns([]string{"title", "doc_category", "technician", "service_type", "sales_person", "service_completion_date", "sales_document_issue_date", "file_url", "synced_at"}),
	}).Create(&data).Error
}

func (r *PostgresRepository) UpsertVaultRecord(data models.VehicleDigitalVault) error {
	return r.UpsertVehicleDigitalVaultRecord(data)
}

func (r *PostgresRepository) GetVehicleDigitalVaultByVIN(vin string) ([]models.VehicleDigitalVault, error) {
    var documents []models.VehicleDigitalVault
    
    result := r.db.Where("vin = ?", vin).Order("event_date DESC").Find(&documents)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return documents, nil
}