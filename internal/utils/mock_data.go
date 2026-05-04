package utils

import (
	"math/rand"
	"time"
)

var PDFLinks = []string{
	"https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	"https://res.cloudinary.com/demo/image/upload/multi_page_pdf.pdf",
	"https://www.adobe.com/support/products/enterprise/knowledgecenter/pdfs/loremipsum.pdf",
	"https://pdfobject.com/pdf/sample.pdf",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_001",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_002",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_003",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_005",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_006",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_007",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_008",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_009",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_012",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_013",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_014",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_015",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_016",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_017",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_018",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_019",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_020",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_021",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_022",
	"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_030",
}

// RandomPDFForSales returns random PDF link suitable for sales documents
func RandomPDFForSales() string {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(PDFLinks)/2) + 1 // Bias towards sales docs (indices 4-11)
	return PDFLinks[idx%len(PDFLinks)]
}

// RandomPDFForService returns random PDF link suitable for service reports
func RandomPDFForService() string {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(PDFLinks)/2) + 12 // Bias towards service docs (indices 12+)
	return PDFLinks[idx%len(PDFLinks)]
}

