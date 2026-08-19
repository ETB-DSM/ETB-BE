# AI-Cane Backend TRD (Technical Requirements Document)

> **프로젝트명:** AI-Cane (ETB)
> **버전:** v1.0
> **작성일:** 2026-08-18
> **참조 문서:** PRD.md

---

## 1. 기술 스택

| 분류 | 기술 | 버전 |
|------|------|------|
| 언어 | Go | 1.22+ |
| HTTP 프레임워크 | Gin | v1 |
| DB | PostgreSQL | 15+ |
| DB 드라이버 | pgx/v5 | v5 |
| SQL 코드 생성 | sqlc | v1 |
| DB 마이그레이션 | golang-migrate | v4 |
| 캐시 | Redis | 7+ |
| Redis 클라이언트 | go-redis/v9 | v9 |
| JWT | golang-jwt/jwt/v5 | v5 |
| WebSocket | gorilla/websocket | v1 |
| SMTP 클라이언트 | gopkg.in/gomail.v2 | v2 |
| Google OAuth 검증 | google.golang.org/api/idtoken | latest |
| 환경 설정 | godotenv | v1 |
| 컨테이너 런타임 | Docker | 24+ |
| 컨테이너 오케스트레이션 | Docker Compose | v2 |

---

## 2. 프로젝트 구조

```
ETB-BE/
├── cmd/
│   └── server/
│       └── main.go              # 서버 진입점
├── internal/
│   ├── handler/                 # HTTP 핸들러 (요청/응답 처리)
│   │   ├── auth.go
│   │   ├── device.go
│   │   ├── destination.go
│   │   ├── guardian.go
│   │   ├── sos.go
│   │   ├── ocr.go
│   │   └── ws.go
│   ├── service/                 # 비즈니스 로직
│   │   ├── auth.go
│   │   ├── device.go
│   │   ├── destination.go
│   │   ├── guardian.go
│   │   ├── sos.go
│   │   ├── ocr.go
│   │   └── ws.go
│   ├── repository/              # DB 쿼리 (sqlc 생성 코드 위치)
│   │   └── sqlc/
│   │       ├── db.go            # sqlc 생성
│   │       ├── models.go        # sqlc 생성
│   │       └── query.sql.go     # sqlc 생성
│   ├── domain/                  # 엔티티, DTO, 에러 정의
│   │   ├── user.go
│   │   ├── device.go
│   │   ├── destination.go
│   │   ├── sos.go
│   │   ├── ws.go
│   │   └── errors.go
│   ├── middleware/
│   │   ├── auth.go              # JWT 검증 미들웨어
│   │   └── logger.go
│   └── infra/
│       ├── postgres.go          # DB 연결 초기화
│       ├── redis.go             # Redis 연결 초기화
│       └── smtp.go             # SMTP 클라이언트 초기화
├── config/
│   └── config.go               # 환경변수 → 구조체 바인딩
├── migrations/
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_devices.up.sql
│   └── ...
├── query/                       # sqlc 입력 SQL 파일
│   ├── user.sql
│   ├── device.sql
│   ├── destination.sql
│   ├── sos.sql
│   └── ocr.sql
├── sqlc.yaml
├── Dockerfile
├── docker-compose.yml
├── .env
├── .env.example
├── go.mod
└── go.sum
```

---

## 3. 환경 설정

### 3.1 .env 구조

```env
# Server
SERVER_PORT=8080

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=aicane
DB_SSLMODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_ACCESS_SECRET=your-access-secret
JWT_REFRESH_SECRET=your-refresh-secret
JWT_ACCESS_EXPIRE_MIN=15
JWT_REFRESH_EXPIRE_DAY=7

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=AI-Cane <your@gmail.com>

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id

# Email Verification
EMAIL_CODE_EXPIRE_MIN=5
```

### 3.2 config.go 구조

```go
type Config struct {
    Server   ServerConfig
    DB       DBConfig
    Redis    RedisConfig
    JWT      JWTConfig
    SMTP     SMTPConfig
    Google   GoogleConfig
}
```

`godotenv.Load()` 후 `os.Getenv()`로 각 필드를 채움. 서버 시작 시 필수 환경변수 누락 여부 검증 후 패닉.

---

## 4. 데이터베이스 스키마

### migrations/000001_create_users.up.sql
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           TEXT NOT NULL UNIQUE,
    nickname        TEXT NOT NULL,
    password_hash   TEXT,                    -- Google OAuth 사용자는 NULL
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### migrations/000002_create_devices.up.sql
```sql
CREATE TABLE devices (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### migrations/000003_create_guardians.up.sql
```sql
CREATE TABLE guardians (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    phone       TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### migrations/000004_create_destinations.up.sql
```sql
CREATE TABLE destinations (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    latitude     DOUBLE PRECISION NOT NULL,
    longitude    DOUBLE PRECISION NOT NULL,
    radius_m     INT NOT NULL DEFAULT 50,
    target_text  TEXT NOT NULL,             -- OCR 비교용 (예: "103동")
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### migrations/000005_create_sos_events.up.sql
```sql
CREATE TABLE sos_events (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id   UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    event_type  TEXT NOT NULL,              -- "fall" | "manual_sos"
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### migrations/000006_create_ocr_logs.up.sql
```sql
CREATE TABLE ocr_logs (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    destination_id   UUID NOT NULL REFERENCES destinations(id) ON DELETE CASCADE,
    recognized_text  TEXT NOT NULL,
    target_text      TEXT NOT NULL,
    matched          BOOLEAN NOT NULL,
    confidence       REAL NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 5. Redis 키 설계

| 키 패턴 | 값 | TTL | 용도 |
|---------|-----|-----|------|
| `verify:{email}` | 6자리 숫자 코드 (string) | 5분 | 이메일 인증 코드 |
| `refresh:{userId}` | Refresh Token (string) | 7일 | JWT Refresh Token |

---

## 6. 인증 설계

### 6.1 이메일 회원가입 + 2차 인증

```
1. POST /auth/signup
   Body: { email, nickname, password }

   → 비밀번호 bcrypt 해시
   → DB에 user 저장 (email_verified=false)
   → 6자리 랜덤 코드 생성
   → Redis SET "verify:{email}" = code, EX 300
   → SMTP로 인증 코드 발송
   ← 200 OK

2. POST /auth/verify-email
   Body: { email, code }

   → Redis GET "verify:{email}"
   → 코드 일치 확인
   → DB UPDATE users SET email_verified=true
   → Redis DEL "verify:{email}"
   ← 200 OK
```

### 6.2 JWT 로그인

```
POST /auth/login
Body: { email, password }

→ DB에서 user 조회 (email_verified=true 확인)
→ bcrypt 비밀번호 검증
→ Access Token 생성  (HS256, exp: 15분)
→ Refresh Token 생성 (HS256, exp: 7일)
→ Redis SET "refresh:{userId}" = refreshToken, EX 604800
← { accessToken, refreshToken }
```

```
POST /auth/refresh
Body: { refreshToken }

→ Refresh Token 서명 검증 → userId 추출
→ Redis GET "refresh:{userId}" 존재 확인
→ 저장된 값과 전달된 토큰 일치 확인
→ 새 Access Token 생성
← { accessToken }
```

```
POST /auth/logout
Header: Authorization: Bearer {accessToken}

→ JWT에서 userId 추출
→ Redis DEL "refresh:{userId}"
← 200 OK
```

### 6.3 Google OAuth

```
POST /auth/login/google
Body: { idToken }

→ idtoken.Validate(ctx, idToken, GOOGLE_CLIENT_ID)
→ payload에서 email, sub(googleId) 추출
→ DB에서 email로 user 조회
  → 없으면: user 생성 (password_hash=NULL, email_verified=true)
  → 있으면: 기존 user 사용
→ Access Token + Refresh Token 발급
→ Redis SET "refresh:{userId}" = refreshToken
← { accessToken, refreshToken }
```

### 6.4 JWT Payload 구조

```json
{
  "sub": "user-uuid",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### 6.5 JWT 미들웨어

`Authorization: Bearer {accessToken}` 헤더 파싱 → 서명 검증 → `gin.Context`에 `userId` 저장.  
Redis 조회 없이 서명만으로 검증 (Access Token 한정).

---

## 7. API 명세

### 공통 응답 포맷

**성공:**
```json
{
  "success": true,
  "data": { ... }
}
```

**실패:**
```json
{
  "success": false,
  "message": "destination not found",
  "errorCode": "DESTINATION_NOT_FOUND"
}
```

**주요 에러 코드:**

| errorCode | HTTP | 설명 |
|-----------|------|------|
| `INVALID_REQUEST` | 400 | 필수 필드 누락 또는 형식 오류 |
| `UNAUTHORIZED` | 401 | 인증 토큰 없음 또는 만료 |
| `FORBIDDEN` | 403 | 권한 없음 |
| `NOT_FOUND` | 404 | 리소스 없음 |
| `CONFLICT` | 409 | 중복 데이터 (이메일 중복 등) |
| `DEVICE_LIMIT_EXCEEDED` | 409 | 디바이스 5개 초과 |
| `GUARDIAN_LIMIT_EXCEEDED` | 409 | 보호자 5개 초과 |
| `INVALID_VERIFY_CODE` | 400 | 이메일 인증 코드 불일치 |
| `EMAIL_NOT_VERIFIED` | 403 | 이메일 미인증 상태로 로그인 시도 |
| `INTERNAL_ERROR` | 500 | 서버 내부 오류 |

---

### 인증 API

#### POST /auth/signup
```json
// Request
{ "email": "user@example.com", "nickname": "홍길동", "password": "P@ssw0rd!" }

// Response 200
{ "success": true, "data": { "message": "인증 코드가 이메일로 발송되었습니다." } }
```

#### POST /auth/verify-email
```json
// Request
{ "email": "user@example.com", "code": "391847" }

// Response 200
{ "success": true, "data": { "message": "이메일 인증이 완료되었습니다." } }
```

#### POST /auth/login
```json
// Request
{ "email": "user@example.com", "password": "P@ssw0rd!" }

// Response 200
{ "success": true, "data": { "accessToken": "...", "refreshToken": "..." } }
```

#### POST /auth/login/google
```json
// Request
{ "idToken": "google-id-token-string" }

// Response 200
{ "success": true, "data": { "accessToken": "...", "refreshToken": "..." } }
```

#### POST /auth/refresh
```json
// Request
{ "refreshToken": "..." }

// Response 200
{ "success": true, "data": { "accessToken": "..." } }
```

#### POST /auth/logout
```
// Header: Authorization: Bearer {accessToken}

// Response 200
{ "success": true, "data": { "message": "로그아웃 되었습니다." } }
```

---

### 디바이스 API (JWT 필요)

#### POST /devices
```json
// Request
{ "name": "내 AI-Cane" }

// Response 201
{ "success": true, "data": { "deviceId": "uuid", "name": "내 AI-Cane", "isActive": true } }
```

#### GET /devices
```json
// Response 200
{ "success": true, "data": [{ "deviceId": "uuid", "name": "내 AI-Cane", "isActive": true, "createdAt": "..." }] }
```

#### DELETE /devices/:deviceId
```json
// Response 200
{ "success": true, "data": { "message": "디바이스가 삭제되었습니다." } }
```

---

### 보호자 API (JWT 필요)

#### POST /guardians
```json
// Request
{ "name": "홍보호자", "phone": "010-1234-5678" }

// Response 201
{ "success": true, "data": { "guardianId": "uuid", "name": "홍보호자", "phone": "010-1234-5678" } }
```

#### GET /guardians
```json
// Response 200
{ "success": true, "data": [{ "guardianId": "uuid", "name": "홍보호자", "phone": "010-1234-5678" }] }
```

#### DELETE /guardians/:guardianId
```json
// Response 200
{ "success": true, "data": { "message": "보호자가 삭제되었습니다." } }
```

---

### 목적지 API (JWT 필요) — BE-002

#### POST /destinations
```json
// Request
{
  "name": "집",
  "latitude": 37.1234,
  "longitude": 127.5678,
  "radiusM": 50,
  "targetText": "103동"
}

// Response 201
{ "success": true, "data": { "destinationId": "uuid", "name": "집", ... } }
```

#### GET /destinations
```json
// Response 200
{
  "success": true,
  "data": [{
    "destinationId": "uuid",
    "name": "집",
    "latitude": 37.1234,
    "longitude": 127.5678,
    "radiusM": 50,
    "targetText": "103동"
  }]
}
```

#### DELETE /destinations/:destinationId
```json
// Response 200
{ "success": true, "data": { "message": "목적지가 삭제되었습니다." } }
```

---

### SOS API (JWT 필요) — BE-004

#### POST /sos
```json
// Request
{
  "deviceId": "uuid",
  "eventType": "fall",
  "latitude": 37.1234,
  "longitude": 127.5678
}

// Response 201
{ "success": true, "data": { "sosId": "uuid", "createdAt": "..." } }
```

#### GET /sos
```json
// Response 200
{
  "success": true,
  "data": [{
    "sosId": "uuid",
    "deviceId": "uuid",
    "eventType": "fall",
    "latitude": 37.1234,
    "longitude": 127.5678,
    "createdAt": "..."
  }]
}
```

---

### OCR 로그 API (JWT 필요) — BE-006 (선택)

#### POST /ocr-logs
```json
// Request
{
  "destinationId": "uuid",
  "recognizedText": "102동",
  "targetText": "103동",
  "matched": false,
  "confidence": 0.91
}

// Response 201
{ "success": true, "data": { "logId": "uuid" } }
```

---

## 8. WebSocket 설계 — BE-003, BE-005

### 연결

```
WS /ws/device
Header: Authorization: Bearer {accessToken}
Query:  ?deviceId={deviceId}
```

연결 시 JWT 검증 + deviceId가 해당 userId 소유인지 확인.

### 클라이언트 → 서버 메시지 포맷

```json
// 실시간 위치 전송 (BE-003)
{
  "type": "location",
  "payload": {
    "latitude": 37.1234,
    "longitude": 127.5678,
    "timestamp": "2026-08-18T12:00:00Z"
  }
}

// 디바이스 상태 전송 (BE-005)
{
  "type": "device_status",
  "payload": {
    "battery": 80,
    "sensorStatus": "ok",
    "networkStatus": "wifi"
  }
}
```

### 서버 → 클라이언트 메시지 포맷

```json
// 수신 확인
{ "type": "ack", "payload": { "receivedType": "location" } }

// 에러
{ "type": "error", "payload": { "message": "unknown message type" } }
```

### sensorStatus 허용 값
`"ok"` | `"disabled"` | `"error"`

### networkStatus 허용 값
`"wifi"` | `"hotspot"` | `"lte"` | `"unknown"`

### 연결 유지
- Ping/Pong: gorilla/websocket 기본 제공, 60초 read deadline
- 비정상 종료 시 Orange Pi에서 재연결

---

## 9. sqlc 설정

### sqlc.yaml
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./query"
    schema: "./migrations"
    gen:
      go:
        package: "repository"
        out: "./internal/repository/sqlc"
        emit_json_tags: true
        emit_prepared_queries: false
```

### query/destination.sql (예시)
```sql
-- name: CreateDestination :one
INSERT INTO destinations (user_id, name, latitude, longitude, radius_m, target_text)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListDestinationsByUser :many
SELECT * FROM destinations
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteDestination :exec
DELETE FROM destinations
WHERE id = $1 AND user_id = $2;

-- name: GetDestination :one
SELECT * FROM destinations
WHERE id = $1 AND user_id = $2;
```

---

## 10. 레이어 간 의존 방향

```
handler → service → repository(sqlc)
                  → infra(redis)
```

- handler: 요청 파싱, 응답 직렬화만 담당. 비즈니스 로직 없음.
- service: 비즈니스 규칙 처리 (디바이스 5개 제한 확인, 코드 생성, 토큰 발급 등)
- repository: sqlc 생성 코드 래핑. DB 쿼리만 담당.
- 각 레이어는 인터페이스로 의존해 테스트 시 mock 교체 가능.

---

## 11. 미들웨어 체인

```
Gin Engine
  └── Logger
  └── Recovery (패닉 → 500 응답)
  └── CORS
  └── /auth/* (인증 불필요)
  └── /ws/*   → WS 핸들러 (JWT 검증 내부 처리)
  └── JWT 미들웨어
      └── /devices/*
      └── /guardians/*
      └── /destinations/*
      └── /sos/*
      └── /ocr-logs/*
```

---

## 12. 개발 순서 (우선순위 기준)

| 단계 | 작업 |
|------|------|
| 1 | 프로젝트 초기 세팅 (go mod, 폴더 구조, .env, DB/Redis 연결) |
| 2 | DB 마이그레이션 작성 및 sqlc 코드 생성 |
| 3 | BE-001 사용자 회원가입·이메일 인증·로그인 (이메일 + Google OAuth) |
| 4 | BE-002 목적지 등록·조회·삭제 |
| 5 | BE-004 SOS 요청 저장 |
| 6 | BE-007 에러 응답 통일 (글로벌 에러 핸들러) |
| 7 | BE-003 WebSocket 실시간 위치 |
| 8 | BE-005 WebSocket 디바이스 상태 |
| 9 | BE-006 OCR 로그 저장 (선택) |

---

## 13. Docker

### 13.1 Dockerfile 구조

멀티스테이지 빌드를 사용한다. 빌더 스테이지에서 Go 바이너리를 컴파일하고, 최종 이미지에는 바이너리와 마이그레이션 파일만 포함한다.

| 스테이지 | 베이스 이미지 | 역할 |
|----------|--------------|------|
| builder | `golang:1.22-alpine` | 의존성 다운로드 + Go 바이너리 컴파일 |
| final | `alpine:3.19` | 최소 런타임 이미지 |

`CGO_ENABLED=0`으로 정적 바이너리를 생성해 alpine 위에서 외부 C 라이브러리 없이 실행 가능하게 한다.  
`ca-certificates`는 Google OAuth의 HTTPS 통신에 필요하다.

### 13.2 docker-compose 서비스 구성

| 서비스 | 이미지 | 역할 |
|--------|--------|------|
| `db` | `postgres:15-alpine` | PostgreSQL 데이터베이스 |
| `redis` | `redis:7-alpine` | 캐시 (JWT Refresh Token, 이메일 인증 코드) |
| `api` | `./Dockerfile` (로컬 빌드) | Go 백엔드 서버 |

### 13.3 서비스 시작 순서

```
db (healthy) ──┐
               ├── api (시작)
redis (healthy)┘
```

`api`는 `db`와 `redis`가 모두 헬스체크를 통과한 뒤 시작한다.  
마이그레이션은 docker compose 외부에서 `Makefile`로 별도 실행한다.

### 13.4 마이그레이션 스크립트 (Makefile)

`Makefile`이 `.env`를 읽어 DB 접속 정보를 자동으로 구성한다. 로컬에 `migrate` CLI를 설치하지 않아도 `docker run`으로 실행한다.

```bash
make migrate-up                        # 전체 마이그레이션 적용
make migrate-down                      # 1단계 롤백
make migrate-down-all                  # 전체 롤백
make migrate-version                   # 현재 마이그레이션 버전 확인
make migrate-create name=create_users  # 새 마이그레이션 파일 쌍 생성
```

### 13.5 주요 docker compose 명령어

```bash
make up       # 전체 스택 시작 (포그라운드)
make up-d     # 전체 스택 시작 (백그라운드)
make down     # 종료 (볼륨 유지)
make down-v   # 종료 + DB 볼륨 초기화
```

### 13.6 환경변수 처리 방식

docker compose는 프로젝트 루트의 `.env` 파일을 자동으로 읽어 Compose 파일 내 `${VAR}` 변수를 치환한다.  
`api` 서비스는 `env_file: .env`로 컨테이너 환경변수를 설정하되, `DB_HOST`와 `REDIS_ADDR`은 Docker 내부 서비스명으로 오버라이드한다.

```yaml
environment:
  DB_HOST: db          # .env의 DB_HOST(localhost)를 덮어씀
  REDIS_ADDR: redis:6379
```

Makefile의 마이그레이션 명령은 `--network host`로 실행되므로 DB_PORT를 통해 로컬호스트로 직접 접근한다.
