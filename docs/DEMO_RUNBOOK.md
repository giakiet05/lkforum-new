# LKForum Demo Runbook

File này dùng để chạy nhanh backend, user web, Prometheus và kiểm tra một vài metric trước khi demo.

## 1. Chạy Backend

Từ thư mục gốc project:

```bash
cd backend
go run main.go
```

Backend mặc định chạy tại:

```text
http://localhost:8080
```

Kiểm tra backend:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl "http://localhost:8080/api/posts?sort=best&feed_type=home&page=1&limit=20"
```

Kết quả mong muốn:

```text
/health -> {"status":"ok"}
/ready  -> mongo ok, redis ok
/api/posts -> success true, có danh sách posts
```

## 2. Chạy User Web

Mở terminal khác, từ thư mục gốc project:

```bash
cd frontend/user-web
npm run dev
```

Nếu frontend gọi nhầm port backend, chạy rõ base URL:

```bash
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```

Mở web:

```text
http://localhost:5173/
```

Lưu ý:

- Dùng `localhost`, không dùng `127.0.0.1`, để tránh lỗi CORS.
- Nếu browser vẫn gọi `localhost:8081`, hard refresh bằng `Ctrl+Shift+R` hoặc đóng tab mở lại.

## 3. Chạy Prometheus

File config:

```text
monitoring/prometheus.yml
```

Config hiện tại scrape backend:

```text
target: localhost:8080
metrics_path: /metrics
scrape_interval: 15s
```

Chạy Prometheus từ thư mục gốc project:

```bash
prometheus --config.file=monitoring/prometheus.yml
```

Mở Prometheus UI:

```text
http://localhost:9090
```

Kiểm tra backend export metric:

```bash
curl http://localhost:8080/metrics
```

## 4. Query Prometheus Nên Demo

### Tổng số request theo endpoint

```promql
lkforum_http_requests_total
```

Ý nghĩa:

```text
Đếm số request HTTP theo method, path, status.
Ví dụ: GET /api/posts status 200.
```

### Thời gian xử lý request

```promql
lkforum_http_request_duration_seconds_count
```

Ý nghĩa:

```text
Đếm số mẫu latency đã ghi nhận theo endpoint.
```

Query latency trung bình:

```promql
rate(lkforum_http_request_duration_seconds_sum[1m])
/
rate(lkforum_http_request_duration_seconds_count[1m])
```

Ý nghĩa:

```text
Tính thời gian xử lý request trung bình trong 1 phút gần nhất.
```

### Feed cache hit/miss

```promql
lkforum_feed_cache_total
```

Ý nghĩa:

```text
result="miss": lần đầu chưa có cache, backend query MongoDB rồi lưu Redis.
result="hit": lần sau cùng query, backend lấy từ Redis.
```

Tạo cache miss/hit:

```bash
curl "http://localhost:8080/api/posts?sort=best&feed_type=home&page=1&limit=20"
curl "http://localhost:8080/api/posts?sort=best&feed_type=home&page=1&limit=20"
curl "http://localhost:8080/api/posts?sort=best&feed_type=home&page=1&limit=20"
```

Sau đó query lại:

```promql
lkforum_feed_cache_total
```

Kết quả mong muốn:

```text
miss tăng ở request đầu.
hit tăng ở các request sau cùng query.
```

### Message encrypted counter

```promql
lkforum_messages_sent_total
```

Ý nghĩa:

```text
Đếm số message đã gửi, có label encrypted.
encrypted="true" nghĩa là message được gửi lên backend dưới dạng ciphertext.
```

Sau khi demo E2EE bằng chat, query:

```promql
lkforum_messages_sent_total
```

Kết quả mong muốn:

```text
lkforum_messages_sent_total{encrypted="true"} tăng.
```

## 5. Checklist Trước Demo

```text
1. Redis đang chạy.
2. Backend chạy ở http://localhost:8080.
3. User web chạy ở http://localhost:5173.
4. User web gọi API sang localhost:8080, không phải 8081.
5. Prometheus chạy ở http://localhost:9090.
6. Prometheus target lkforum-backend ở trạng thái UP.
7. Query lkforum_http_requests_total có dữ liệu.
8. Query lkforum_feed_cache_total có hit/miss sau khi gọi /api/posts nhiều lần.
```
