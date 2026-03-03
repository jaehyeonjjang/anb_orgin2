# 보수점검 백엔드 구조 문서

## 📁 전체 구조

```
repair/back/
├── main.go                 # 진입점
├── .env.yml               # 환경 설정 (develop/production)
├── go.mod                 # Go 모듈 의존성
├── index.sql              # 인덱스 DDL
├── projectile.sql         # 테이블 DDL
├── global/                # 전역 설정 및 유틸리티
│   ├── config/           # 설정 관리
│   ├── log/              # 로깅
│   ├── time/             # 시간 유틸
│   ├── export.go         # 엑셀/PDF 내보내기
│   ├── zip.go            # 압축 처리
│   ├── image.go          # 이미지 처리
│   ├── pdf.go            # PDF 생성
│   └── notify.go         # 알림 처리
├── models/                # 데이터 모델 (DB 레이어)
│   ├── db.go             # DB 연결 기본 클래스
│   ├── cache.go          # 캐시 관리
│   ├── data.go           # 공통 데이터 처리
│   ├── periodic/         # 정기점검 모델
│   ├── repair/           # 보수 모델
│   ├── patrol/           # 순회점검 모델
│   ├── user/             # 사용자 모델
│   ├── apt/              # 아파트 모델
│   └── ... (총 100개 이상의 도메인 모델)
├── controllers/           # 컨트롤러 (비즈니스 로직)
│   ├── controllers.go    # 기본 컨트롤러
│   ├── api/              # API 컨트롤러
│   │   ├── periodic.go   # 정기점검 API
│   │   ├── repair.go     # 보수 API
│   │   ├── patrol.go     # 순회점검 API
│   │   └── ... (25개 이상의 API 컨트롤러)
│   └── rest/             # REST API 컨트롤러
├── router/                # 라우팅
│   ├── router.go         # 라우트 정의 (7600+ 라인)
│   └── auth.go           # JWT 인증
├── services/              # 백그라운드 서비스
│   ├── http.go           # HTTP 서버 (Gin)
│   ├── cron.go           # 스케줄러
│   ├── notify.go         # 알림 서비스
│   └── report.go         # 리포트 생성
├── chat/                  # 채팅 (Socket.IO)
├── estimate/              # 견적
├── repair/                # 보수 관련 추가 로직
└── periodic/              # 정기점검 관련 추가 로직
```

---

## 🚀 시작 흐름 (main.go)

```go
func main() {
    log.Init()                          // 1. 로깅 초기화
    models.InitCache()                  // 2. 캐시 초기화
    os.MkdirAll("webdata/detail", 0755) // 3. 업로드 디렉토리 생성
    
    go services.Cron()                  // 4. 스케줄러 시작 (백그라운드)
    go services.Notify()                // 5. 알림 서비스 시작 (백그라운드)
    services.Http()                     // 6. HTTP 서버 시작 (메인)
}
```

### 실행 모드
- **develop**: 포트 9108, 개발 DB 사용
- **production**: 포트 8080, 운영 DB 사용
- **test**: 포트 9109, 테스트 DB 사용 (SQLite)

---

## ⚙️ 설정 (.env.yml)

```yaml
develop:
  port: 9108
  cors:
    - http://localhost:4000
  path: webdata
  database:
    type: mysql
    host: dev.netb.co.kr
    port: 3306
    name: repair
    user: yuhki
    password: dkTkfl123!

production:
  port: 8080
  path: /usr/local/main/webdata
  database:
    type: mysql
    host: 10.34.96.4
    port: 3306
    name: repair
    user: repair
    password: repairdb
```

### 설정 로딩 (global/config/config.go)
- YAML 파일 파싱
- `MODE` 환경변수로 모드 결정 (develop/production/test)
- `ENV_FILE` 환경변수로 설정 파일 지정 가능

---

## 🗄️ 데이터베이스 레이어

### 1. DB 연결 (models/db.go)

```go
type Connection struct {
    Conn        *sql.DB
    Tx          *sql.Tx
    Transaction bool
    Isolation   bool
}
```

**주요 기능:**
- MySQL/PostgreSQL/SQLite 지원
- 트랜잭션 관리
- 쿼리 실행 추상화

### 2. 모델 구조

각 도메인별로 독립적인 패키지:
- `models/periodic/` - 정기점검
- `models/repair/` - 보수
- `models/patrol/` - 순회점검
- `models/user/` - 사용자
- `models/apt/` - 아파트
- ... 총 100개 이상

### 3. 모델 예시 (periodic)

```go
type Periodic struct {
    Id              int64   // 정기점검 ID
    Name            string  // 점검명
    Apt             int64   // 아파트 ID
    Status          int     // 상태
    Category        int     // 카테고리
    Reportdate      int     // 보고일
    Startdate       int     // 시작일
    Enddate         int     // 종료일
    User            int64   // 작성자
    // ... 50개 이상의 필드
}
```

---

## 🎮 컨트롤러 레이어

### 기본 컨트롤러 (controllers/controllers.go)

```go
type Controller struct {
    Ctx     *gin.Context
    Session *Session
    Code    int
    Result  gin.H
    DB      *models.Connection
}
```

**주요 메서드:**
- `Init(c *gin.Context)` - 초기화 및 세션 로드
- `NewConnection()` - DB 연결 생성
- `Close()` - 리소스 정리

### API 컨트롤러 예시 (api/periodic.go)

```go
type PeriodicController struct {
    controllers.Controller
}

func (c *PeriodicController) Pre_Insert(item *models.Periodic)
func (c *PeriodicController) Post_Insert(item *models.Periodic)
func (c *PeriodicController) Search()
func (c *PeriodicController) Update()
func (c *PeriodicController) Delete()
```

---

## 🛣️ 라우팅 (router/router.go)

### 인증 흐름

```
1. /api/jwt              → JWT 토큰 발급 (로그인)
2. /api/* (protected)    → JWT 검증 미들웨어
3. API 엔드포인트        → 컨트롤러 실행
```

### 라우트 패턴

```go
apiGroup := r.Group("/api")
apiGroup.Use(JwtAuthRequired())  // JWT 인증 미들웨어
{
    apiGroup.GET("/periodic/search", func(c *gin.Context) {
        var controller api.PeriodicController
        controller.Init(c)        // 초기화 (세션, DB 연결)
        controller.Search()       // 비즈니스 로직
        controller.Close()        // 리소스 정리
        c.JSON(controller.Code, controller.Result)
    })
}
```

### 주요 엔드포인트 그룹

- `/api/periodic/*` - 정기점검
- `/api/repair/*` - 보수
- `/api/patrol/*` - 순회점검
- `/api/apt/*` - 아파트
- `/api/user/*` - 사용자
- `/api/estimate/*` - 견적
- `/api/breakdown/*` - 하자
- `/api/file/*` - 파일 업로드/다운로드
- `/api/report/*` - 리포트

**총 라우트 수:** 7600+ 라인에 수백 개의 엔드포인트

---

## 🌐 HTTP 서버 (services/http.go)

### Gin 웹 프레임워크 사용

```go
r := gin.Default()
r.Use(CORSMiddleware())              // CORS 설정

// Socket.IO (채팅)
r.GET("/socket.io/", gin.WrapH(c.Server))
r.POST("/socket.io/", gin.WrapH(c.Server))

// 정적 파일
r.Static("/assets", "./dist/assets")
r.Static("/webdata", config.UploadPath)  // 업로드된 파일

// SPA 라우팅 (Vue.js)
r.NoRoute(func(c *gin.Context) {
    c.File("./dist/index.html")
})

router.SetRouter(r)  // API 라우트 등록
```

### CORS 설정
- 설정 파일의 `cors` 배열에 정의된 origin 허용
- `Access-Control-Allow-Credentials: true`
- 메서드: GET, POST, PUT, DELETE

---

## 📦 주요 서비스

### 1. 스케줄러 (services/cron.go)
- 주기적 작업 실행
- 데이터 동기화
- 리포트 자동 생성

### 2. 알림 서비스 (services/notify.go)
- 실시간 알림 처리
- WebSocket 통신

### 3. 리포트 생성 (services/report.go)
- PDF 리포트 생성
- Excel 내보내기
- 데이터 집계

---

## 💾 데이터베이스 스키마

### 주요 테이블

#### 1. 사용자 & 회사
- `user_tb` - 사용자
- `company_tb` - 회사
- `aptgroup_tb` - 아파트 그룹

#### 2. 정기점검
- `periodic_tb` - 정기점검 기본
- `p_damage_tb` - 균열 손상
- `p_joint_tb` - 신축이음
- `p_leak_tb` - 누수
- `p_slope_tb` - 기울기
- `p_strength_tb` - 강도
- `p_detail_tb` - 상세 데이터
- `p_picture_tb` - 사진

#### 3. 보수
- `repair_tb` - 보수 기본
- `repairlist_tb` - 보수 목록
- `breakdown_tb` - 하자

#### 4. 순회점검
- `patrol_tb` - 순회점검
- `patrolimage_tb` - 순회점검 이미지

#### 5. 기타
- `standard_tb` - 표준 데이터
- `category_tb` - 카테고리
- `estimate_tb` - 견적
- `contract_tb` - 계약

**총 테이블 수:** 100개 이상

---

## 🔐 인증 (JWT)

### JWT 토큰 발급

```
POST /api/jwt?loginid=test&passwd=test123
→ { "token": "eyJhbGc..." }
```

### JWT 검증 미들웨어

```go
func JwtAuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // Bearer 토큰 검증
        // 세션 정보 컨텍스트에 저장
    }
}
```

### 세션 구조

```go
type Session struct {
    Id       int64   // 사용자 ID
    Loginid  string  // 로그인 ID
    Name     string  // 이름
    Level    int     // 권한 레벨
    Company  int64   // 회사 ID
}
```

---

## 📤 파일 업로드/다운로드

### 업로드 디렉토리
- `webdata/` - 기본 업로드 경로
- `webdata/detail/` - 상세 데이터

### 파일 API
- `POST /api/file/upload` - 파일 업로드
- `GET /webdata/{filename}` - 파일 다운로드
- `POST /api/image/*` - 이미지 처리

---

## 🧪 테스트 환경

### 테스트 모드 특징
1. **SQLite 사용** - MySQL 대신 로컬 SQLite DB
2. **독립적인 데이터** - test.db 파일로 분리
3. **별도 포트** - 9109 (develop: 9108, production: 8080)
4. **별도 업로드 경로** - webdata_test/

### 테스트 서버 시작

```bash
# 1. 초기 데이터로 시작
./start_test_server.sh

# 2. 프로덕션 데이터 복사해서 시작
./copy_production_data.sh         # 데이터 복사
./start_test_server_with_data.sh  # 서버 시작
```

### 환경변수

```bash
MODE=test ENV_FILE=.env.test.yml ./bin/backend_test
```

---

## 🔄 데이터 흐름

```
1. HTTP 요청
   ↓
2. 라우터 (router/router.go)
   ↓
3. JWT 인증 (router/auth.go)
   ↓
4. 컨트롤러 초기화 (controllers/controllers.go)
   ├─ 세션 로드
   └─ DB 연결
   ↓
5. 비즈니스 로직 (controllers/api/*.go)
   ├─ 모델 조회/수정 (models/*)
   ├─ 트랜잭션 처리
   └─ 파일 처리
   ↓
6. 응답 반환 (JSON)
   ↓
7. 리소스 정리 (DB 연결 닫기)
```

---

## 🛠️ 개발 도구

### 빌드
```bash
go build -o bin/backend main.go
```

### 실행
```bash
# 개발 모드
MODE=develop ./bin/backend

# 프로덕션 모드
MODE=production ./bin/backend

# 테스트 모드
MODE=test ENV_FILE=.env.test.yml ./bin/backend
```

### 데이터베이스 초기화
```bash
# MySQL
mysql -u user -p database < projectile.sql

# SQLite (테스트)
sqlite3 test.db < index.sql
sqlite3 test.db < projectile.sql
```

---

## 📝 코딩 컨벤션

### 파일 구조
- 각 도메인별로 독립적인 패키지
- 모델, 컨트롤러, 서비스 분리

### 네이밍
- 모델: `{Domain}` (예: `Periodic`, `Repair`)
- 컨트롤러: `{Domain}Controller` (예: `PeriodicController`)
- 매니저: `New{Domain}Manager()` (예: `NewPeriodicManager()`)

### 에러 처리
- 컨트롤러에서 `Code`, `Result`로 응답
- 로그는 `log` 패키지 사용

---

## 🚦 API 응답 형식

### 성공
```json
{
  "code": 200,
  "data": { ... },
  "message": "success"
}
```

### 실패
```json
{
  "code": 400,
  "error": "에러 메시지",
  "message": "fail"
}
```

---

## 📊 주요 기능별 API

### 정기점검
- `GET /api/periodic/search` - 검색
- `POST /api/periodic/insert` - 등록
- `POST /api/periodic/update` - 수정
- `POST /api/periodic/delete` - 삭제
- `GET /api/periodic/report` - 리포트

### 보수
- `GET /api/repair/search` - 검색
- `POST /api/repair/insert` - 등록
- `POST /api/repair/update` - 수정
- `GET /api/repair/export` - Excel 내보내기

### 순회점검
- `GET /api/patrol/search` - 검색
- `POST /api/patrol/insert` - 등록
- `POST /api/patrol/image` - 이미지 업로드

---

## 🎯 다음 단계

1. **테스트 작성** - main_test.go 기반으로 통합 테스트
2. **API 문서화** - Swagger 또는 Postman 컬렉션
3. **성능 최적화** - 쿼리 최적화, 캐싱 강화
4. **모니터링** - 로그 분석, 에러 추적
5. **배포 자동화** - Docker, CI/CD

---

## 📚 참고 자료

- [Gin 웹 프레임워크](https://github.com/gin-gonic/gin)
- [GORM (ORM)](https://gorm.io/) - 현재는 직접 SQL 사용
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Socket.IO Go](https://github.com/googollee/go-socket.io)
