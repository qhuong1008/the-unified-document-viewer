package utils

import "fmt"

func main() {
	// A slice of 30 mock PDF links for UI testing and database seeding
	mockPDFLinks := []string{
		// 1-5: Stable public assets
		"https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
		"https://res.cloudinary.com/demo/image/upload/multi_page_pdf.pdf",
		"https://www.adobe.com/support/products/enterprise/knowledgecenter/pdfs/loremipsum.pdf",
		"https://pdfobject.com/pdf/sample.pdf",
		"https://www.clickdimensions.com/links/Test.pdf",

		// 6-15: GitHub Raw Sample Files (Technical/Complex PDFs)
		"https://raw.githubusercontent.com/py-pdf/sample-files/main/001-form.pdf",
		"https://raw.githubusercontent.com/py-pdf/sample-files/main/002-various-page-sizes.pdf",
		"https://raw.githubusercontent.com/py-pdf/sample-files/main/003-rotated-pages.pdf",
		"https://raw.githubusercontent.com/py-pdf/sample-files/main/004-labels.pdf",
		"https://raw.githubusercontent.com/py-pdf/sample-files/main/009-re-encoded-font.pdf",
		"https://raw.githubusercontent.com/py-pdf/sample-files/main/011-table.pdf",
		"https://raw.githubusercontent.com/mozilla/pdf.js/master/web/compressed.tracemonkey-pldi-09.pdf",
		"https://raw.githubusercontent.com/africau/sample-pdf/master/sample.pdf",
		"https://raw.githubusercontent.com/mathiasbynens/small/master/pdf.pdf",
		"https://raw.githubusercontent.com/wkhtmltopdf/wkhtmltopdf/master/test/functional/samples/table.pdf",

		// 16-30: Generated placeholders (guaranteed to work and unique)
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_016",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_017",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_018",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_019",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_020",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_021",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_022",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_023",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_024",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_025",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_026",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_027",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_028",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_029",
		"https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_030",
	}

	// Example usage
	fmt.Printf("Generated %d mock PDF links.\n", len(mockPDFLinks))
	fmt.Println("First link:", mockPDFLinks[0])
}