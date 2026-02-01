# BasicBoard Backend

게시판 프로젝트 백엔드 서버 (Go + Gin + GORM + MySQL)

## 📋 목차

- [기술 스택](#기술-스택)
- [프로젝트 구조](#프로젝트-구조)
- [시작하기](#시작하기)
- [환경 설정](#환경-설정)
- [API 명세](#api-명세)
- [데이터 모델](#데이터-모델)
- [주요 기능](#주요-기능)

## 🛠 기술 스택

- **언어**: Go 1.25.1
- **웹 프레임워크**: Gin
- **ORM**: GORM
- **데이터베이스**: MySQL
- **환경 변수**: godotenv

## 📁 프로젝트 구조

```
BasicBoard_Backend/
├── main.go              # 서버 진입점 및 라우팅
├── database/
│   └── db.go           # 데이터베이스 연결 및 초기화
├── models/
│   └── post.go         # 데이터 모델 정의
├── handlers/
│   ├── post.go         # 게시글 관련 핸들러
│   └── search.go       # 검색 관련 핸들러
├── go.mod
├── go.sum
└── README.md
```

## 🚀 시작하기

### 1. 의존성 설치

```bash
go mod download
```

### 2. 환경 변수 설정

프로젝트 루트에 `.env` 파일을 생성하고 다음 내용을 추가하세요:

```env
DB_USER=root
DB_PASSWORD=your_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=basicboard_db
```

### 3. 데이터베이스 생성

MySQL에서 데이터베이스를 생성합니다:

```sql
CREATE DATABASE basicboard_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 서버 실행

```bash
go run main.go
```

서버는 기본적으로 `http://localhost:8080`에서 실행됩니다.

## ⚙️ 환경 설정

환경 변수는 `.env` 파일을 통해 설정하며, 기본값은 다음과 같습니다:

| 변수 | 기본값 | 설명 |
|------|-------|------|
| `DB_USER` | `root` | MySQL 사용자명 |
| `DB_PASSWORD` | `password` | MySQL 비밀번호 |
| `DB_HOST` | `127.0.0.1` | MySQL 호스트 |
| `DB_PORT` | `3306` | MySQL 포트 |
| `DB_NAME` | `basicboard_db` | 데이터베이스 이름 |

## 📡 API 명세

모든 API는 `/api` prefix를 사용합니다.

### 게시글 API

#### 1. 게시글 목록 조회

**웹 (페이지 기반)**
```
GET /api/posts?page=1&limit=5
```

**앱 (커서 기반)**
```
GET /api/posts?limit=5&cursorCreatedAt=2026-01-24T12:34:56.000Z&cursorId=10
```

**응답 (웹)**
```json
{
  "items": [Post[]],
  "page": 1,
  "limit": 5,
  "total": 57,
  "totalPages": 12
}
```

**응답 (앱)**
```json
{
  "items": [Post[]],
  "nextCursor": {
    "createdAt": "2026-01-24T12:34:56.000Z",
    "id": 10
  },
  "hasMore": true
}
```

**마지막 페이지일 경우**
```json
{
  "items": [Post[]],
  "nextCursor": null,
  "hasMore": false
}
```

#### 2. 게시글 상세 조회

```
GET /api/posts/:id
```

**응답**
```json
{
  "id": 3,
  "title": "제목",
  "content": "본문",
  "createdAt": "2026-01-24T12:34:56.000Z",
  "updatedAt": "2026-01-24T12:40:00.000Z"
}
```

**상태 코드**
- `200`: 성공
- `404`: 대상없음(없는id)
- `500`: 서버 내부 오류

#### 3. 게시글 작성

```
POST /api/posts
Content-Type: application/json
```

**요청**
```json
{
  "title": "제목",
  "content": "본문"
}
```

**응답**
```json
Post
```

**상태 코드**
- `201`: 성공
- `400`: 요청값 잘못됨(누락 또는 제목 길이 제한 초과)
- `500`: 서버 내부 오류

**제목 길이 제한**
- 한글: 최대 45자
- 영어: 최대 72자

#### 4. 게시글 수정

```
PUT /api/posts/:id
Content-Type: application/json
```

**요청**
```json
{
  "title": "수정제목",
  "content": "수정본문"
}
```

**응답**
```json
Post
```

**상태 코드**
- `201`: 성공
- `400`: 요청값 잘못됨(누락 또는 제목 길이 제한 초과)
- `404`: 대상없음(없는id)
- `500`: 서버 내부 오류

#### 5. 게시글 삭제

```
DELETE /api/posts/:id
```

**상태 코드**
- `204`: 삭제 성공
- `404`: 대상없음(없는id)
- `500`: 서버 내부 오류

### 검색 API

#### 1. 검색

**웹 (페이지 기반)**
```
GET /api/search?q=검색어&page=1&limit=5
```

**앱 (커서 기반)**
```
GET /api/search?q=검색어&limit=5&cursorCreatedAt=2026-01-24T12:34:56.000Z&cursorId=10
```

**검색 규칙**
- 검색어에서 **모든 공백을 제거**한 후 검색
- 예: "안녕 하세요" → "안녕하세요"로 변환하여 검색
- 공백만 입력한 경우: 빈 결과 반환 (검색 기록 저장 안 함)

**응답 (웹)**
```json
{
  "items": [Post[]],
  "page": 1,
  "limit": 5,
  "total": 10,
  "totalPages": 2
}
```

**응답 (앱)**
```json
{
  "items": [Post[]],
  "nextCursor": {
    "createdAt": "2026-01-24T12:34:56.000Z",
    "id": 10
  },
  "hasMore": true
}
```

**상태 코드**
- `200`: 성공
- `400`: 요청값 잘못됨(검색어 누락)
- `500`: 서버 내부 오류

#### 2. 최근 검색어 조회

```
GET /api/search/history
```

**응답**
```json
{
  "history": [
    {
      "id": 1,
      "keyword": "검색어",
      "searchedAt": "2026-01-24T12:34:56.000Z"
    }
  ]
}
```

- 최대 10개까지 반환
- 최신순으로 정렬

#### 3. 검색어 삭제

```
DELETE /api/search/history/:id
```

**응답**
```json
{
  "message": "Deleted"
}
```

#### 4. 검색어 전체 삭제

```
DELETE /api/search/history
```

**응답**
```json
{
  "message": "All cleared"
}
```

## 📊 데이터 모델

### Post (게시글)

```go
type Post struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title" gorm:"size:255;not null"`
    Content   string    `json:"content" gorm:"type:text;not null"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

### SearchHistory (검색 기록)

```go
type SearchHistory struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    Keyword    string    `json:"keyword" gorm:"size:255;not null;unique"`
    SearchedAt time.Time `json:"searchedAt"`
}
```

### 응답 타입

#### PageResponse (웹용 페이지 기반)

```go
type PageResponse[T any] struct {
    Items     []T   `json:"items"`
    Page      int   `json:"page"`
    Limit     int   `json:"limit"`
    Total     int64 `json:"total"`
    TotalPages int  `json:"totalPages"`
}
```

#### CursorResponse (앱용 커서 기반)

```go
type CursorResponse[T any] struct {
    Items      []T     `json:"items"`
    NextCursor *Cursor  `json:"nextCursor"`
    HasMore    bool     `json:"hasMore"`
}
```

#### Cursor (커서 정보)

```go
type Cursor struct {
    CreatedAt time.Time `json:"createdAt"`
    ID        uint      `json:"id"`
}
```

## ✨ 주요 기능

### 1. 페이지네이션

- **웹**: 페이지 기반 (`page`, `limit`)
- **앱**: 커서 기반 (`cursorCreatedAt`, `cursorId`)
- 정렬 기준: `createdAt DESC, id DESC`

### 2. 검색 기능

- 제목 또는 본문에서 검색어 포함 여부 검색
- 검색어 공백 자동 제거
- 검색 기록 자동 저장 (최대 10개)
- 검색 기록 조회/삭제 기능

### 3. 입력 검증

- 제목/본문 필수 입력
- 제목 길이 제한 (한글 45자/영어 72자)
- 공백만 입력된 경우 자동 제거 후 검증

### 4. CORS 지원

- 모든 Origin 허용 (`*`)
- OPTIONS 요청 자동 처리

## 🔧 개발 참고사항

### 정렬 기준

모든 목록 조회는 다음 기준으로 정렬됩니다:
- `createdAt DESC` (최신순)
- `id DESC` (동일 시간일 경우 ID 기준)

### 에러 처리

- `400`: 잘못된 요청 (누락된 필드, 길이 제한 초과 등)
- `404`: 리소스를 찾을 수 없음
- `500`: 서버 내부 오류

### 데이터베이스 마이그레이션

서버 시작 시 `AutoMigrate`를 통해 테이블이 자동으로 생성됩니다:
- `posts` 테이블
- `search_history` 테이블

## 📝 라이선스

이 프로젝트는 개인 프로젝트입니다.
