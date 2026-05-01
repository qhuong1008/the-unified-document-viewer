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
    
    // Fixed: Use synced_at (or created_at) instead of non-existent event_date column
    result := r.db.Where("vin = ?", vin).Order("synced_at DESC").Find(&documents)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return documents, nil
}

// CheckIfVINExists checks if any records exist for a given VIN
func (r *PostgresRepository) CheckIfVINExists(vin string) (bool, int64, error) {
    var count int64
    
    result := r.db.Model(&models.VehicleDigitalVault{}).Where("vin = ?", vin).Count(&count)
    
    if result.Error != nil {
        return false, 0, result.Error
    }
    
    return count > 0, count, nil
}
