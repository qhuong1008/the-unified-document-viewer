#!/bin/bash

BASE_URL="http://localhost:8080/webhooks"

# Danh sách PDF đã làm sạch, chỉ giữ lại các link hoạt động ổn định
PDF_LINKS=(
    "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
    "https://res.cloudinary.com/demo/image/upload/multi_page_pdf.pdf"
    "https://www.adobe.com/support/products/enterprise/knowledgecenter/pdfs/loremipsum.pdf"
    "https://pdfobject.com/pdf/sample.pdf"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_001"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_002"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_003"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_005"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_006"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_007"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_008"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Sales_Doc_009"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_012"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_013"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_014"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_015"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_016"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_017"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_018"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_019"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_020"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_021"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_022"
    "https://placehold.jp/24/cccccc/ffffff/1000x1400.pdf?text=Service_Report_030"
)

# Tạo ID ngẫu nhiên bằng cơ chế hệ thống để tránh phụ thuộc vào Python
get_uuid() {
    cat /proc/sys/kernel/random/uuid 2>/dev/null || echo "$RANDOM-$RANDOM"
}

get_random_pdf() {
    printf "%s\n" "${PDF_LINKS[$RANDOM % ${#PDF_LINKS[@]}]}"
}

send_sales_request() {
    local vin=$1
    local id=$(get_uuid)
    local pdf=$(get_random_pdf)
    
    # Kích hoạt background process (&) để tạo tính song song thực sự
    curl -s -X POST "$BASE_URL/sales" \
        -H "Content-Type: application/json" \
        -d "{
          \"id\": \"$id\",
          \"vin\": \"$vin\",
          \"sales_person\": \"Senior Partner\",
          \"created_at\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
          \"file_url\": \"$pdf\"
        }" > /dev/null &
}

send_service_request() {
    local vin=$1
    local id=$(get_uuid)
    local pdf=$(get_random_pdf)
    
    curl -s -X POST "$BASE_URL/service" \
        -H "Content-Type: application/json" \
        -d "{
          \"id\": \"$id\",
          \"vin\": \"$vin\",
          \"service_type\": \"Maintenance\",
          \"technician\": \"Expert AI\",
          \"completion_date\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
          \"report_link\": \"$pdf\"
        }" > /dev/null &
}

echo "🚀 Starting Parallel Ingestion Test - Clean Version..."

VINS=("VIN1001" "VIN2002" "VIN3003" "VIN4004" "VIN5005")

for vin in "${VINS[@]}"
do
    send_sales_request "$vin"
    send_service_request "$vin"
done

wait
echo -e "\n✅ All requests dispatched. Check your Worker Pool logs!"