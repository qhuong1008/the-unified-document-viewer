package worker

import (
	"strings"
	"the-unified-document-viewer/internal/models"
)

// EnrichVaultData thực hiện làm giàu dữ liệu sau khi đã Transform
func EnrichVaultData(vault *models.VehicleDigitalVault) {
	// 2. Phân loại mức độ ưu tiên (Priority Enrichment)[cite: 3]
	// Giả định: Tài liệu Technical (bảo trì) hoặc Commercial quan trọng sẽ có priority cao
	if vault.DocCategory == "Technical" || strings.Contains(strings.ToLower(vault.Title), "contract") {
		// Use asterisk instead of emoji to avoid UTF-8 encoding issues with database
		vault.Title = "[PRIORITY] " + vault.Title // Thêm badge cho UI dễ nhận diện[cite: 1]
	}

	// 3. Xử lý logic hiển thị cho Title[cite: 1]
	// Đảm bảo Title luôn bắt đầu bằng chữ hoa để chuyên nghiệp
	if len(vault.Title) > 0 {
		vault.Title = strings.ToUpper(vault.Title[:1]) + vault.Title[1:]
	}
	
	// 4. Sanitize the title to ensure proper UTF-8 encoding for database persistence
	vault.Title = sanitizeUTF8(vault.Title)
}