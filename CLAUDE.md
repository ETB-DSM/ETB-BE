# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Restrictions

- `.env` 파일은 절대 읽거나 수정하지 않는다. 환경변수 구조가 필요하면 `.env.example`을 참조한다.

## Project Reference Documents

현재 컨텍스트에 프로젝트 요구사항이나 기술 스택 정보가 없으면 아래 문서를 먼저 읽는다.

- `PRD.md` — 기능 요구사항, 사용자 흐름, API 목록
- `TRD.md` — 기술 스택, 아키텍처, DB 스키마, 인증 설계, WebSocket 설계

## Commands

```bash
make run          # 서버 실행
make build        # 빌드
make test         # 전체 테스트
make sqlc         # sqlc 코드 생성 (query/*.sql 변경 후 반드시 실행)
make lint         # 린트
```

```bash
make up           # Docker 전체 스택 시작 (포그라운드)
make up-d         # Docker 전체 스택 시작 (백그라운드)
make down         # 종료 (볼륨 유지)
make down-v       # 종료 + DB 볼륨 초기화
```

```bash
make migrate-up                        # 전체 마이그레이션 적용
make migrate-down                      # 1단계 롤백
make migrate-down-all                  # 전체 롤백
make migrate-version                   # 현재 버전 확인
make migrate-create name=<이름>         # 마이그레이션 파일 쌍 생성
```

---

## Work Flow

모든 작업은 아래 3단계를 순서대로 따른다.

### Step 1 — 설계

사용자 요구사항을 수신하면 곧바로 코드를 작성하지 않는다.

1. 요구사항을 분석해 작업 방식을 도출한다.
   - 어떤 파일을 생성·수정할지
   - 함수·메서드 시그니처
   - 데이터 흐름 및 레이어 간 책임 분리
   - API 요청/응답 포맷 변경이 있으면 명시
   - DB 스키마 변경이 있으면 명시
2. 도출한 내용을 `docs/YYYY.MM.DD_<description>.md` 에 작성한다.
   - 날짜는 작업 당일 기준
   - description은 영문 kebab-case (예: `2026.08.18_auth-signup-flow.md`)
3. 설계 문서를 사용자에게 공유하고 피드백을 기다린다.

### Step 2 — 피드백

사용자의 피드백을 반영해 설계 문서를 수정한다.  
사용자가 **"진행해줘"** 라고 하면 Step 3으로 넘어간다.

### Step 3 — 개발 및 검증

#### 브랜치 생성
`develop` 브랜치에서 새 브랜치를 만들어 작업한다.

```
feat/<description>   # 신규 기능
fix/<description>    # 버그 수정
chore/<description>  # 설정·문서·도구 변경
```

#### 개발
설계 문서(docs/)를 기준으로 코드를 작성한다.

#### 검증 (sub-agent)
코드 작성이 완료되면 **새 세션의 Claude Sonnet sub-agent**를 실행해 검증한다.

sub-agent에게 전달하는 내용:
- 해당 작업의 설계 문서 전문
- 작성된 코드 전문
- 기존 코드베이스 중 변경된 파일과 연관된 파일

sub-agent가 검증하는 항목:
- 구현이 설계 문서와 일치하는가
- 기존 코드베이스와 정합성이 맞는가 (레이어 의존 방향, 에러 포맷, 네이밍 컨벤션 등)

검증 결과 불일치가 있으면:
1. Claude가 코드를 수정한다.
2. sub-agent를 다시 실행해 재검증한다.
3. 검증이 통과될 때까지 반복한다.

#### 머지
검증 통과 후 작업 브랜치를 `develop`에 직접 merge한다.

```bash
git checkout develop
git merge --no-ff <브랜치명>
git branch -d <브랜치명>
```
