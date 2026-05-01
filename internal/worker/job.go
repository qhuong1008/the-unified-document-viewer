package worker

type SourceType string

const (
	SourceSales   SourceType = "SALES"
	SourceService SourceType = "SERVICE"
)

// Job represents a unit of raw data waiting to be transformed and persisted.
type Job struct {
	Type    SourceType
	Payload interface{} // Will hold either models.RawSalesData or models.RawServiceData
}
