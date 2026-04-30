package worker

// SourceType defines where the data originated from
type SourceType string

const (
	SourceSales   SourceType = "SALES"
	SourceService SourceType = "SERVICE"
)

// Job represents a unit of raw data waiting to be transformed and persisted.
// Using an interface{} for Payload allows the JobQueue to be generic for different sources.
type Job struct {
	Type    SourceType
	Payload interface{} // Will hold either models.RawSalesData or models.RawServiceData
}
