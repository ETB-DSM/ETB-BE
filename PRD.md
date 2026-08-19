# AI-Cane Backend PRD

> **프로젝트명:** AI-Cane (ETB)
> **버전:** v1.0
> **작성일:** 2026-08-18
> **상태:** In-progress

---

## 1. 프로젝트 개요

**AI-Cane은 시각장애인의 안전한 독립 보행을 지원하기 위한 Orange Pi 5 Plus 기반 엣지 AI 스마트 지팡이 시스템이다.**

백엔드 서버는 Orange Pi 5 Plus 디바이스와 REST API / WebSocket으로 통신하며, 목적지 정보 저장·조회, SOS 이벤트 기록, 디바이스 상태 모니터링, OCR 로그 저장의 역할을 담당한다.

---

## 2. 배경 및 문제 정의

### 2.1 기존 흰지팡이의 한계

| 문제 | 설명 |
|------|------|
| 상부 장애물 미감지 | 간판, 나뭇가지, 돌출 구조물 등 머리 높이 장애물 탐지 불가 |
| 최종 목적지 탐색 불가 | GPS로 목적지 근처 도달은 가능하나 건물명·동 번호 구분 불가 |
| 응급 상황 대응 부재 | 낙상 발생 시 보호자에게 위치 전달 수단 없음 |

### 2.2 AI-Cane의 해결 방식

- **LiDAR 5개** (정면·좌측·좌전방·우측·상단) 로 상부 포함 전방위 장애물 감지
- **OCR** 기반 건물명·동 번호 인식으로 최종 목적지 식별
- **MPU6050** 낙상 감지 + GPS 좌표 전송으로 SOS 대응
- **진동 피드백** 으로 소리·화면 없이 방향 안내

---

## 3. 시스템 아키텍처

```
사용자
  │
  ▼
AI-Cane 지팡이 (Orange Pi 5 Plus)
  ├── Embedded Layer
  │     ├── TF-Luna LiDAR x5 (장애물 거리 측정)
  │     ├── MPU-6050 (낙상 감지)
  │     ├── NEO-M8N GPS (위치 측정)
  │     └── YwRobot 진동모터 (사용자 피드백)
  │
  ├── AI Layer
  │     ├── YOLO (평상시 — 사람·차량·자전거 등 객체 탐지, NPU 경량 모델)
  │     └── OCR (목적지 근처에서만 조건부 실행 — 건물명·동 번호 인식)
  │
  └── Backend Communication (REST API / WebSocket)
        │
        ▼
      Backend Server (Go)
        ├── REST API  → 목적지 등록/조회, SOS 저장, OCR 로그
        └── WebSocket → 실시간 위치 스트리밍, 디바이스 상태
              │
              ▼
            Database (PostgreSQL)
```

### 3.1 AI 실행 전략 (상황별 활성화)

| 기능 | 실행 방식 |
|------|-----------|
| LiDAR 장애물 감지 | 항상 실행 |
| YOLO 객체 인식 | NPU 기반 경량 모델, 평상시 실행 |
| OCR 문자 인식 | 목적지 반경 진입 시에만 조건부 실행 |
| GPS 계산 | CPU 처리, 항상 실행 |
| 낙상 감지 | CPU 처리, 항상 실행 |
| 서버 통신 | 이벤트 발생 시 |

---

## 4. 사용자 및 주요 엔티티

### 4.1 사용자 (User)

- 회원가입: 이메일 + 닉네임 + 비밀번호 입력 후 이메일 2차 인증으로 계정 활성화
- 로그인: 이메일/비밀번호 방식 또는 Google OAuth 방식 중 선택
- 디바이스 최대 **5개** 등록 가능 (삭제 포함)
- 보호자 연락처 최대 **5개** 등록 가능 (삭제 포함)

### 4.2 디바이스 (Device)

- `deviceId` 로 Orange Pi 5 Plus 식별
- 디바이스 활성 상태 관리

### 4.3 목적지 (Destination)

- 목적지명, GPS 좌표(위도·경도), 도착 반경(m), OCR 비교용 `targetText` 보관
- 예시 `targetText`: `"103동"`

---

## 5. 기능 요구사항

### BE-001 — 사용자 및 장치 정보 관리 (⭐⭐ 중요)

**목적:** 서버에서 사용자와 AI-Cane 장치를 구분한다.

| 항목 | 내용 |
|------|------|
| 입력 | 사용자 ID, 장치 ID, 보호자 정보 |
| 출력 | 저장된 사용자/장치 정보 |
| 완료 기준 | 테스트 사용자와 장치 정보가 서버에 저장되고 조회 가능 |
| 구현 기준 | 테스트용 userId/deviceId 생성, 보호자 연락처와 장치 활성 상태 저장 |
| 의존 | 없음 |

**동작 순서:**
1. 사용자 정보 입력
2. 장치 ID 연결
3. 보호자 연락처 저장
4. 조회 API 제공

**예외:** 중복 사용자, 필수 정보 누락, 잘못된 전화번호 형식

> MVP에서는 `user_001`, `aicane_001` 고정 계정으로 시작 가능

---

### BE-002 — 목적지 등록 및 조회 (⭐⭐⭐ 필수)

**목적:** GPS 목적지 탐색과 OCR 비교에 필요한 데이터를 제공한다.

| 항목 | 내용 |
|------|------|
| 입력 | 목적지 이름, 좌표(위도·경도), 반경(m), targetText |
| 출력 | destinationId, 목적지 목록 |
| 완료 기준 | 목적지 등록 후 Orange Pi가 목록 조회 가능 |
| 구현 기준 | 목적지명, 위도/경도, 도착 반경, OCR 비교용 targetText 저장 및 목록 조회 |
| 의존 | BE-001 |

**동작 순서:**
1. 목적지 등록 요청 수신
2. 좌표와 targetText 검증
3. DB 저장
4. 목적지 목록 조회 응답

**예외:** 좌표 누락, targetText 누락, 존재하지 않는 userId

---

### BE-003 — 실시간 위치 전송 (WebSocket)

**목적:** 보행 중 현재 GPS 위치를 실시간으로 서버로 전달한다.

> DB 저장이 아닌 **WebSocket** 으로 처리한다. (질문사항 검토 결과: DB 저장 시 활용 가치 불명확, 실시간 스트리밍이 적합)

| 항목 | 내용 |
|------|------|
| 방식 | WebSocket 연결 유지 |
| 입력 | 위도, 경도, 타임스탬프 |
| 출력 | 수신 확인 |
| 주요 용도 | SOS 발생 시 즉시 위치 확인, 목적지 근처 진입 판단 보조 |

---

### BE-004 — SOS 요청 저장 (⭐⭐⭐ 필수)

**목적:** 낙상 감지 또는 긴급 요청 발생 시 위치와 이벤트 정보를 서버에 저장한다.

| 항목 | 내용 |
|------|------|
| 입력 | eventType, GPS 좌표, userId, deviceId |
| 출력 | sosId, 저장 결과 |
| 완료 기준 | 낙상 이벤트 발생 시 SOS 요청이 DB에 저장됨 |
| 구현 기준 | 낙상/SOS 이벤트를 userId, deviceId, GPS 좌표, eventType과 함께 저장 |
| 의존 | BE-001, EMB-007 |

**동작 순서:**
1. SOS payload 수신
2. userId와 deviceId 확인
3. 위치와 eventType 저장
4. sosId 반환

**예외:** 네트워크 요청 중복, 위치 정보 없음, DB 저장 실패

> 보호자 문자 발송은 후순위 (서버 저장 먼저 구현). 네트워크 실패 시 임베디드에서 재시도.

---

### BE-005 — 디바이스 상태 저장 (⭐ 선택)

**목적:** 시연 중 AI-Cane 장치 상태를 서버에서 확인한다 (디버깅용).

> DB 저장이 아닌 **WebSocket** 으로 처리한다. (질문사항 검토 결과)

| 항목 | 내용 |
|------|------|
| 입력 | 배터리 %, 센서 상태(ok/disabled/error), 네트워크 상태 |
| 출력 | 디바이스 상태 로그 |
| 완료 기준 | 상태 payload가 서버에 저장되고 조회 가능 |
| 의존 | BE-001, EMB-008 |

---

### BE-006 — OCR 결과 저장 (⭐ 선택)

**목적:** OCR 인식 결과와 목적지 일치 여부를 저장해 AI 성능 개선에 활용한다.

| 항목 | 내용 |
|------|------|
| 입력 | recognizedText, targetText, matched(bool), confidence(float) |
| 출력 | OCR 결과 로그 |
| 완료 기준 | OCR 성공/실패 결과가 서버에 저장됨 |
| 구현 기준 | OCR 결과 문자열, targetText, matched, confidence 저장 |
| 의존 | AI-003, AI-004, BE-002 |

---

### BE-007 — API 오류 응답 처리 (⭐⭐ 중요)

**목적:** Orange Pi가 서버 오류를 예측 가능한 방식으로 처리하게 한다.

| 항목 | 내용 |
|------|------|
| 입력 | 잘못된 요청, 없는 ID, 서버 오류 |
| 출력 | 일관된 오류 JSON |
| 완료 기준 | 400/404/500 상황에서 정해진 형식의 오류 응답 반환 |
| 구현 기준 | 모든 API에서 `success/message/errorCode` 형식의 일관된 JSON 응답 |
| 의존 | BE-001, BE-002, BE-003, BE-004 |

**공통 오류 응답 형식:**
```json
{
  "success": false,
  "message": "destination not found",
  "errorCode": "DESTINATION_NOT_FOUND"
}
```

---

## 6. 인증 및 계정

### 6.1 회원가입 (이메일)

이메일 + 닉네임 + 비밀번호로 가입하며, 가입 즉시 로그인은 불가하고 이메일 인증을 완료해야 계정이 활성화된다.

**흐름:**
1. 사용자가 이메일·닉네임·비밀번호 입력 후 가입 요청
2. 서버가 DB에 계정 생성 (`email_verified=false`)
3. SMTP를 통해 6자리 인증 코드를 이메일로 발송, 코드는 **5분간 유효**
4. 사용자가 인증 코드 입력
5. 코드 일치 시 계정 활성화 (`email_verified=true`), 이후 로그인 가능

**예외:**
- 이미 가입된 이메일로 재가입 시도 → 409 CONFLICT
- 인증 코드 불일치 → 400 INVALID_VERIFY_CODE
- 인증 코드 만료(5분 초과) → 400 INVALID_VERIFY_CODE (재발송 필요)
- 이메일 미인증 상태에서 로그인 시도 → 403 EMAIL_NOT_VERIFIED

---

### 6.2 로그인 (이메일/비밀번호)

이메일 인증이 완료된 계정에 한해 로그인 가능하다.

**흐름:**
1. 사용자가 이메일·비밀번호 입력
2. 서버가 비밀번호 검증 및 이메일 인증 여부 확인
3. **Access Token** (유효기간 15분) 과 **Refresh Token** (유효기간 7일) 발급
4. 이후 API 요청 시 `Authorization: Bearer {accessToken}` 헤더로 인증

**토큰 갱신:**
- Access Token 만료 시 Refresh Token으로 새 Access Token 발급
- Refresh Token은 서버(Redis)에 저장되며, 로그아웃 시 즉시 무효화

**예외:**
- 존재하지 않는 이메일 → 401 UNAUTHORIZED
- 비밀번호 불일치 → 401 UNAUTHORIZED
- 이메일 미인증 계정 → 403 EMAIL_NOT_VERIFIED

---

### 6.3 로그인 (Google OAuth)

Google 계정으로 간편 로그인한다. 별도 회원가입 절차 없이 최초 로그인 시 계정이 자동 생성된다.

**흐름:**
1. 클라이언트가 Google 로그인 후 ID Token 수신
2. 클라이언트가 서버로 ID Token 전달
3. 서버가 Google 공개키로 ID Token 서명 직접 검증
4. 검증된 이메일로 DB 사용자 조회
   - 없으면: 신규 계정 자동 생성 (비밀번호 없음, 이메일 인증 완료 처리)
   - 있으면: 기존 계정 사용
5. Access Token + Refresh Token 발급

**예외:**
- 유효하지 않은 ID Token → 401 UNAUTHORIZED
- Google Client ID 불일치 → 401 UNAUTHORIZED

---

### 6.4 계정 제한 사항

| 항목 | 제한 |
|------|------|
| 디바이스 등록 | 계정당 최대 5개, 삭제 가능 |
| 보호자 등록 | 계정당 최대 5개, 삭제 가능 |
| 이메일 인증 코드 유효시간 | 5분 |
| Access Token 유효시간 | 15분 |
| Refresh Token 유효시간 | 7일 |

---

## 7. 데이터 모델 (초안)

### users
| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK |
| email | TEXT UNIQUE | 이메일 |
| nickname | TEXT | 닉네임 |
| password_hash | TEXT | 비밀번호 해시 (OAuth 미사용 시) |
| email_verified | BOOLEAN | 이메일 인증 여부 |
| created_at | TIMESTAMP | |

### devices
| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK (deviceId) |
| user_id | UUID | FK → users |
| name | TEXT | 디바이스 이름 |
| is_active | BOOLEAN | 활성 상태 |
| created_at | TIMESTAMP | |

### guardians
| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK |
| user_id | UUID | FK → users |
| phone | TEXT | 보호자 연락처 |
| name | TEXT | 보호자 이름 |

### destinations
| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK (destinationId) |
| user_id | UUID | FK → users |
| name | TEXT | 목적지 이름 |
| latitude | FLOAT8 | 위도 |
| longitude | FLOAT8 | 경도 |
| radius_m | INT | 도착 판정 반경(m) |
| target_text | TEXT | OCR 비교용 문자열 (예: "103동") |
| created_at | TIMESTAMP | |

### sos_events
| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK (sosId) |
| user_id | UUID | FK → users |
| device_id | UUID | FK → devices |
| event_type | TEXT | "fall" \| "manual_sos" |
| latitude | FLOAT8 | 발생 위도 |
| longitude | FLOAT8 | 발생 경도 |
| created_at | TIMESTAMP | |

### ocr_logs (선택)
| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK |
| destination_id | UUID | FK → destinations |
| recognized_text | TEXT | OCR 인식 결과 |
| target_text | TEXT | 비교 대상 문자열 |
| matched | BOOLEAN | 일치 여부 |
| confidence | FLOAT4 | 신뢰도 |
| created_at | TIMESTAMP | |

---

## 8. API 목록 (초안)

### 인증
| Method | Path | 설명 |
|--------|------|------|
| POST | `/auth/signup` | 이메일 회원가입 |
| POST | `/auth/verify-email` | 이메일 인증 코드 확인 |
| POST | `/auth/login` | 이메일/비밀번호 로그인 |
| POST | `/auth/login/google` | Google OAuth 로그인 |
| POST | `/auth/refresh` | Access Token 갱신 |
| POST | `/auth/logout` | 로그아웃 (Refresh Token 무효화) |

### 디바이스
| Method | Path | 설명 |
|--------|------|------|
| POST | `/devices` | 디바이스 등록 |
| GET | `/devices` | 디바이스 목록 조회 |
| DELETE | `/devices/:deviceId` | 디바이스 삭제 |

### 보호자
| Method | Path | 설명 |
|--------|------|------|
| POST | `/guardians` | 보호자 등록 |
| GET | `/guardians` | 보호자 목록 조회 |
| DELETE | `/guardians/:guardianId` | 보호자 삭제 |

### 목적지 (BE-002)
| Method | Path | 설명 |
|--------|------|------|
| POST | `/destinations` | 목적지 등록 |
| GET | `/destinations` | 목적지 목록 조회 |
| DELETE | `/destinations/:destinationId` | 목적지 삭제 |

### SOS (BE-004)
| Method | Path | 설명 |
|--------|------|------|
| POST | `/sos` | SOS 이벤트 저장 |
| GET | `/sos` | SOS 이력 조회 |

### OCR 로그 (BE-006, 선택)
| Method | Path | 설명 |
|--------|------|------|
| POST | `/ocr-logs` | OCR 결과 저장 |

### WebSocket
| Path | 설명 |
|------|------|
| `WS /ws/location` | 실시간 GPS 위치 스트리밍 (BE-003) |
| `WS /ws/device-status` | 디바이스 상태 스트리밍 (BE-005) |

---

## 9. 개발 우선순위

### 1순위 (필수 — MVP)
- [ ] BE-002 목적지 등록 및 조회
- [ ] BE-004 SOS 요청 저장
- [ ] BE-007 API 오류 응답 처리

### 2순위 (중요)
- [ ] BE-001 사용자 및 장치 정보 관리
- [ ] BE-003 실시간 위치 WebSocket

### 3순위 (선택)
- [ ] BE-005 디바이스 상태 WebSocket
- [ ] BE-006 OCR 결과 저장

---

## 10. 비기능 요구사항

| 항목 | 기준 |
|------|------|
| 언어 | Go |
| 오류 응답 형식 | `{ "success": bool, "message": string, "errorCode": string }` 일관 유지 |
| 재시도 | 네트워크 실패 시 임베디드(Orange Pi)에서 재시도, 서버는 멱등성 보장 |
| SOS 중복 처리 | 동일 이벤트 중복 수신 방어 |
| 디바이스 제한 | 계정당 디바이스 최대 5개, 보호자 최대 5개 |

---

## 11. 미결 사항

| 항목 | 현황 |
|------|------|
| 보호자 SOS 알림 방식 | 문자 발송 후순위, 구체적 방식 미정 (SMS API vs 앱 푸시) |
| SOS LTE 통신 | Orange Pi 5 Plus는 LTE 미내장 → Wi-Fi/핫스팟/별도 LTE 모듈 중 선택 필요 |
| 목적지 반경 기본값 | 기본 도착 판정 반경(m) 미정 |
| OCR 신뢰도 임계값 | matched 판정 기준 confidence 값 미정 |
| DB 선택 | PostgreSQL 가정, 최종 확정 필요 |
