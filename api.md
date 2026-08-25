# API 명세서

> **관리 규칙:** API가 추가·변경·삭제될 때마다 이 문서를 반드시 갱신한다.

---

## 공통

### Base URL

| 환경 | URL |
|------|-----|
| 로컬 | `http://localhost:{SERVER_PORT}` |
| Docker | `http://localhost:{SERVER_PORT}` |

### 응답 포맷

```json
// 성공 — 래퍼 없이 데이터 직접 반환
{ "accessToken": "eyJ...", "refreshToken": "eyJ..." }

// 실패
{ "message": "...", "errorCode": "..." }
```

### 에러 코드

| errorCode | HTTP | 설명 |
|-----------|------|------|
| `INVALID_REQUEST` | 400 | 필수 파라미터 누락 또는 형식 오류 |
| `UNAUTHORIZED` | 401 | 인증 토큰 없음 또는 만료 |
| `FORBIDDEN` | 403 | 리소스 접근 권한 없음 |
| `NOT_FOUND` | 404 | 리소스 없음 |
| `CONFLICT` | 409 | 중복 등록 |
| `DEVICE_LIMIT_EXCEEDED` | 409 | 디바이스 5개 초과 |
| `GUARDIAN_LIMIT_EXCEEDED` | 409 | 보호자 5개 초과 |
| `INVALID_VERIFY_CODE` | 400 | 이메일 인증 코드 불일치 또는 만료 |
| `EMAIL_NOT_VERIFIED` | 403 | 이메일 미인증 계정 |
| `INTERNAL_ERROR` | 500 | 서버 내부 오류 |

---

## App API — `/api/v1/...`

JWT 인증이 필요한 엔드포인트는 `Authorization: Bearer <accessToken>` 헤더를 포함한다.

---

### 인증 (Auth)

#### `POST /api/v1/auth/signup`

회원가입. 이메일 인증 코드를 발송한다.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| email | string | ✓ | 이메일 |
| nickname | string | ✓ | 닉네임 |
| password | string | ✓ | 비밀번호 (최소 8자) |

**Response** `201`
```json
{ "message": "인증 코드가 발송되었습니다." }
```

---

#### `POST /api/v1/auth/verify-email`

이메일 인증 코드 확인.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| email | string | ✓ | 이메일 |
| code | string | ✓ | 6자리 인증 코드 |

**Response** `200`
```json
{ "message": "이메일 인증이 완료되었습니다." }
```

---

#### `POST /api/v1/auth/login`

이메일 로그인.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| email | string | ✓ | 이메일 |
| password | string | ✓ | 비밀번호 |

**Response** `200`
```json
{
  "accessToken": "eyJ...",
  "refreshToken": "eyJ..."
}
```

---

#### `POST /api/v1/auth/login/google`

Google OAuth ID Token으로 로그인.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| idToken | string | ✓ | Google ID Token |

**Response** `200`
```json
{
  "accessToken": "eyJ...",
  "refreshToken": "eyJ..."
}
```

---

#### `POST /api/v1/auth/refresh`

Access Token + Refresh Token 재발급 (token rotation).

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| refreshToken | string | ✓ | Refresh Token |

**Response** `200`
```json
{
  "accessToken": "eyJ...",
  "refreshToken": "eyJ..."
}
```

---

#### `POST /api/v1/auth/logout` 🔒

로그아웃. Redis의 Refresh Token을 삭제한다.

**Response** `200`
```json
{ "message": "로그아웃 되었습니다." }
```

---

### 디바이스 (Device) 🔒

#### `POST /api/v1/devices`

디바이스 등록 (최대 5개).

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| name | string | ✓ | 디바이스 이름 |

**Response** `201`
```json
{
  "deviceId": "uuid",
  "name": "My Cane",
  "isActive": true,
  "apiKey": "a3f8c2...",
  "createdAt": "2026-08-21T00:00:00Z"
}
```

> `apiKey` 는 등록 직후 1회만 응답에 포함된다. Orange Pi 설정 파일에 저장해야 한다.

---

#### `GET /api/v1/devices`

디바이스 목록 조회.

**Response** `200` — 배열

---

#### `DELETE /api/v1/devices/:deviceId`

디바이스 삭제.

**Response** `200`

---

### 보호자 (Guardian) 🔒

#### `POST /api/v1/guardians`

보호자 등록 (최대 5개).

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| name | string | ✓ | 보호자 이름 |
| phone | string | ✓ | 보호자 전화번호 |

**Response** `201`
```json
{ "guardianId": "uuid", "name": "홍길동", "phone": "010-0000-0000" }
```

---

#### `GET /api/v1/guardians`

보호자 목록 조회.

**Response** `200` — 배열

---

#### `DELETE /api/v1/guardians/:guardianId`

보호자 삭제.

**Response** `200`

---

### 목적지 (Destination) 🔒

#### `POST /api/v1/destinations`

목적지 등록.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| name | string | ✓ | 목적지 이름 |
| latitude | float64 | ✓ | 위도 |
| longitude | float64 | ✓ | 경도 |
| radiusM | int32 | ✓ | 도착 판정 반경(m) |
| targetText | string | ✓ | OCR 매칭 대상 텍스트 |

**Response** `201`
```json
{
  "destinationId": "uuid",
  "name": "회사",
  "latitude": 37.123,
  "longitude": 127.456,
  "radiusM": 50,
  "targetText": "회사 건물",
  "createdAt": "2026-08-21T00:00:00Z"
}
```

---

#### `GET /api/v1/destinations`

목적지 목록 조회.

**Response** `200` — 배열

---

#### `DELETE /api/v1/destinations/:destinationId`

목적지 삭제.

**Response** `200`

---

### SOS 🔒

#### `POST /api/v1/sos`

SOS 요청 (앱 → 서버).

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| deviceId | string | ✓ | 디바이스 ID |
| eventType | string | ✓ | `fall` \| `manual_sos` \| `fall_detected` \| `emergency_button` |
| latitude | float64 | ✓ | 위도 |
| longitude | float64 | ✓ | 경도 |
| battery | int32 | | 배터리 잔량 |
| timestamp | string | | ISO 8601 발생 시각 |

**Response** `201`
```json
{
  "sosId": "uuid",
  "deviceId": "uuid",
  "eventType": "fall",
  "latitude": 37.123,
  "longitude": 127.456,
  "sentToGuardian": false,
  "createdAt": "2026-08-21T00:00:00Z"
}
```

---

#### `GET /api/v1/sos`

SOS 이벤트 목록 조회.

**Response** `200` — 배열

---

### OCR 로그 🔒

#### `POST /api/v1/ocr-logs`

OCR 결과 저장.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| destinationId | string | ✓ | 목적지 ID |
| recognizedText | string | ✓ | OCR 인식 결과 |
| targetText | string | ✓ | 매칭 대상 텍스트 |
| matched | bool | | 일치 여부 |
| confidence | float32 | ✓ | 신뢰도 (0~1) |

**Response** `201`
```json
{ "logId": "uuid" }
```

---

### 길안내 세션 (Navigation) 🔒

#### `POST /api/v1/navigation/sessions`

길안내 세션 생성.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| destinationId | string | ✓ | 목적지 ID |
| deviceId | string | | 디바이스 ID |
| startLatitude | float64 | | 출발 위도 |
| startLongitude | float64 | | 출발 경도 |

**Response** `201`
```json
{
  "sessionId": "uuid",
  "status": "active",
  "startLatitude": 37.123,
  "startLongitude": 127.456,
  "createdAt": "2026-08-21T00:00:00Z"
}
```

---

#### `POST /api/v1/navigation/sessions/:sessionId/instructions`

길안내 명령 저장 (앱 → 서버, Orange Pi가 polling).

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| action | string | ✓ | `straight` \| `prepare_left` \| `left` \| `prepare_right` \| `right` \| `arrived` \| `reroute` |
| distanceMeters | int32 | | 다음 분기점까지 거리(m) |
| message | string | | 음성 안내 메시지 |

**Response** `201`
```json
{
  "instructionId": "uuid",
  "sessionId": "uuid",
  "action": "straight",
  "distanceMeters": 150,
  "message": "150m 직진하세요.",
  "createdAt": "2026-08-21T00:00:00Z"
}
```

---

#### `PATCH /api/v1/navigation/sessions/:sessionId/status`

세션 상태 변경.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| status | string | ✓ | `active` \| `paused` \| `arrived` \| `canceled` \| `error` |
| reason | string | | 변경 사유 |

**Response** `200`
```json
{
  "sessionId": "uuid",
  "status": "arrived",
  "createdAt": "2026-08-21T00:00:00Z"
}
```

---

## Embedded API — `/api/...`

Orange Pi에서 호출. JWT 없이 `userId`(문자열)를 body 또는 query로 전달한다.

---

#### `POST /api/users`

Orange Pi 사용자 최초 등록.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| userId | string | ✓ | 클라이언트 지정 사용자 ID |
| name | string | ✓ | 사용자 이름 |
| guardianName | string | ✓ | 보호자 이름 |
| guardianPhone | string | ✓ | 보호자 전화번호 |
| deviceId | string | ✓ | AI-Cane 장치 ID |

**Response** `201`
```json
{ "userId": "user_001", "deviceId": "aicane_001" }
```

---

#### `GET /api/users/:userId`

사용자 존재 확인.

**Response** `200`
```json
{ "userId": "user_001", "deviceId": "aicane_001" }
```

---

#### `POST /api/location`

GPS 위치 저장.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| userId | string | ✓ | 사용자 ID |
| deviceId | string | ✓ | 장치 ID |
| latitude | float64 | ✓ | 위도 |
| longitude | float64 | ✓ | 경도 |
| timestamp | string | ✓ | ISO 8601 측정 시각 |

**Response** `201`
```json
{ "message": "위치가 저장되었습니다." }
```

---

#### `POST /api/sos`

SOS 요청 (Orange Pi → 서버).

**Header**

| 헤더 | 필수 | 설명 |
|------|------|------|
| X-Device-Key | ✓ | 디바이스 등록 시 발급된 API Key |

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| eventType | string | ✓ | `fall` \| `manual_sos` \| `fall_detected` \| `emergency_button` |
| latitude | float64 | ✓ | 위도 |
| longitude | float64 | ✓ | 경도 |
| battery | int32 | | 배터리 잔량 |
| timestamp | string | | ISO 8601 발생 시각 |

**Response** `201`
```json
{ "sosId": "uuid", "sentToGuardian": false }
```

**Error** `401 UNAUTHORIZED` — X-Device-Key 누락 또는 유효하지 않은 키

---

#### `POST /api/device/status`

디바이스 상태 업데이트.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| deviceId | string | ✓ | 장치 ID |
| battery | int32 | ✓ | 배터리 잔량 (0~100) |
| lidarOk | bool | | LiDAR 정상 여부 |
| cameraOk | bool | | 카메라 정상 여부 |
| gpsOk | bool | | GPS 정상 여부 |
| networkOk | bool | | 네트워크 정상 여부 |

**Response** `200`
```json
{
  "deviceId": "aicane_001",
  "battery": 85,
  "lidarOk": true,
  "cameraOk": true,
  "gpsOk": true,
  "networkOk": true,
  "createdAt": "2026-08-21T00:00:00Z"
}
```

---

#### `GET /api/device/status/:deviceId`

최신 디바이스 상태 조회.

**Response** `200` — 위와 동일

---

#### `POST /api/ocr/results`

OCR 결과 저장.

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| destinationId | string | ✓ | 목적지 ID |
| recognizedText | string | ✓ | OCR 인식 결과 |
| targetText | string | ✓ | 매칭 대상 텍스트 |
| matched | bool | | 일치 여부 |
| confidence | float32 | ✓ | 신뢰도 (0~1) |

**Response** `201`
```json
{ "logId": "uuid" }
```

---

#### `POST /api/destinations`

목적지 등록 (Embedded).

**Request Body**

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| userId | string | ✓ | 사용자 ID |
| name | string | ✓ | 목적지 이름 |
| latitude | float64 | ✓ | 위도 |
| longitude | float64 | ✓ | 경도 |
| radiusM | int32 | ✓ | 도착 판정 반경(m) |
| targetText | string | ✓ | OCR 매칭 대상 텍스트 |

**Response** `201`

---

#### `GET /api/destinations?userId=xxx`

목적지 목록 조회 (Embedded).

**Query Parameter:** `userId` (필수)

**Response** `200` — 배열

---

#### `GET /api/navigation/sessions/:sessionId/instruction`

최신 길안내 명령 조회 (Orange Pi polling).

**Response** `200`
```json
{
  "instructionId": "uuid",
  "sessionId": "uuid",
  "action": "left",
  "distanceMeters": 30,
  "message": "30m 후 좌회전하세요.",
  "createdAt": "2026-08-21T00:00:00Z"
}
```

---

## WebSocket

#### `GET /ws/device?deviceId=xxx` 🔒

Orange Pi와 실시간 연결. JWT 인증 후 deviceId 소유권 검증.

**클라이언트 → 서버 메시지**

```json
{ "type": "location", "payload": { "latitude": 37.123, "longitude": 127.456 } }
{ "type": "device_status", "payload": { "battery": 80 } }
```

**서버 → 클라이언트 응답**

```json
{ "type": "ack", "refType": "location" }
{ "type": "error", "message": "unknown message type" }
```

> Ping/Pong: 50초 간격으로 서버가 Ping 전송. 60초 read deadline.

---

## 헬스체크

#### `GET /health`

```json
{ "status": "ok" }
```
