#!/bin/bash

BASE_URL="http://localhost:8080/webhooks"

# Danh sách 30 link PDF thật từ yêu cầu của bạn
PDF_LINKS=(
    "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
    "https://res.cloudinary.com/demo/image/upload/multi_page_pdf.pdf"
    "https://www.adobe.com/support/products/enterprise/knowledgecenter/pdfs/loremipsum.pdf"
    "https://pdfobject.com/pdf/sample.pdf"
    "https://www.clickdimensions.com/links/Test.pdf"
    "https://raw.githubusercontent.com/py-pdf/sample-files/main/001-form.pdf"
    "https://raw.githubusercontent.com/py-pdf/sample-files/main/002-various-page-sizes.pdf"
    "https://raw.githubusercontent.com/py-pdf/sample-files/main/003-rotated-pages.pdf"
    "https://raw.githubusercontent.com/py-pdf/sample-files/main/004-labels.pdf"
    "https://raw.githubusercontent.com/py-pdf/sample-files/main/009-re-encoded-font.pdf"
    "https://raw.githubusercontent.com/py-pdf/sample-files/main/011-table.pdf"
    "https://raw.githubusercontent.com/mozilla/pdf.js/master/web/compressed.tracemonkey-pldi-09.pdf"
    "https://raw.githubusercontent.com/africau/sample-pdf/master/sample.pdf"
    "https://raw.githubusercontent.com/mathiasbynens/small/master/pdf.pdf"
    "https://raw.githubusercontent.com/wkhtmltopdf/wkhtmltopdf/master/test/functional/samples/table.pdf"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_016"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_017"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_018"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_019"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_020"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_021"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_022"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_023"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_024"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_025"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_026"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_027"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_028"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_029"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_030"
)

# Hàm bổ trợ tạo UUID bằng Python để tránh lỗi 'command not found'
get_uuid() {
    python3 -c 'import uuid; print(uuid.uuid4())'
}

# Hàm lấy ngẫu nhiên 1 link từ mảng PDF_LINKS
get_random_pdf() {
    printf "%s\n" "${PDF_LINKS[$RANDOM % ${#PDF_LINKS[@]}]}"
}

send_sales_request() {
    local vin=$1
    local id=$(get_uuid)
    local pdf=$(get_random_pdf)
    
    curl -X POST "$BASE_URL/sales" \
        -H "Content-Type: application/json" \
        -d "{
          \"id\": \"$id\",
          \"vin\": \"$vin\",
          \"sales_person\": \"Senior Partner\",
          \"created_at\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
          \"file_url\": \"$pdf\"
        }" &
}

send_service_request() {
    local vin=$1
    local id=$(get_uuid)
    local pdf=$(get_random_pdf)
    
    curl -X POST "$BASE_URL/service" \
        -H "Content-Type: application/json" \
        -d "{
          \"id\": \"$id\",
          \"vin\": \"$vin\",
          \"service_type\": \"Maintenance\",
          \"technician\": \"Expert AI\",
          \"completion_date\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
          \"report_link\": \"$pdf\"
        }" &
}

echo "🚀 Starting Parallel Ingestion Test with REAL PDF Links..."

# Danh sách VIN để kiểm tra tính duy nhất và Upsert
VINS=("VIN1001" "VIN2002" "VIN3003" "VIN4004" "VIN5005")

for vin in "${VINS[@]}"
do
    send_sales_request "$vin"
    send_service_request "$vin"
done

wait
echo "✅ All requests dispatched. Check your Worker Pool logs!"