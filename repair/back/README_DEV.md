# 개발 환경 가이드

## 🎯 개발 환경 종류

### 1. **Local (권장) - 포트 9107**
```bash
./init_local.sh      # 최초 1회 실행 (Docker + MySQL)
./start_local.sh     # 서버 시작
./stop_local.sh      # 서버 중지
```
- **용도**: 개인 로컬 개발
- **DB**: MySQL (Docker, localhost:3307)
- **장점**: 
  - 독립적인 환경, 다른 사람에게 영향 없음
  - DB 초기화 자유롭게 가능
  - **프로덕션과 동일한 MySQL 사용**
  - Docker로 간편한 관리
- **계정**: dev / dev123

### 2. **Test - 포트 9109**
```bash
./start_test_server.sh                  # 테스트 데이터
./start_test_server_with_data.sh        # 프로덕션 데이터 복사
```
- **용도**: 통합 테스트, QA
- **DB**: SQLite (test.db)
- **장점**: 프로덕션과 유사한 데이터로 테스트

### 3. **Develop - 포트 9108 (공유 서버, 비추천)**
```bash
MODE=develop ./bin/backend
```
- **용도**: 전임자가 사용하던 공유 개발 서버
- **DB**: MySQL (dev.netb.co.kr)
- **⚠️ 주의**: 
  - 다른 개발자와 DB 공유
  - 데이터 충돌 가능
  - 함부로 초기화하면 안 됨

### 4. **Production - 포트 8080**
- **절대 로컬에서 실행하지 말 것**
- 운영 서버 전용

---

## 사전 준비: Docker 설치

Docker Desktop을 설치하세요:
- Mac: https://www.docker.com/products/docker-desktop
- 설치 후 Docker를 실행해두세요

### 1단계: 초기화 (최초 1회)

```bash
cd /Users/seongjaehyeon/anb_origin2/repair/back
chmod +x init_local.sh start_local.sh stop_local.sh
./init_local.sh
```

**선택 옵션:**
- `1) 최소 테스트 데이터` - 빠르게 시작 (권장)
- `2) 프로덕션 데이터 복사` - 실제 데이터로 테스트 (MySQL 필요)
- `1) 최소 테스트 데이터` - 빠르게 시작 (권장)
- `2) 프로덕션 데이터 복사` - 실제 데이터로 테스트

### 2단계: 서버 시작

```bash
./start_local.sh
```

### 3단계: 접속

```
http://localhost:9107
```

**개발 계정:**
- ID: `dev`
- PW: `dev123`

---    # 로컬 설정 (MySQL)
├── .env.test.yml               # 테스트 설정 (SQLite)
├── .env.yml                    # develop/production 설정
├── docker-compose.local.yml    # 로컬 MySQL 컨테이너
├── test.db                     # 테스트 DB (SQLite)
├── webdata_local/              # 로컬 업로드 경로
├── webdata_test/               # 테스트 업로드 경로
├── logs/                       # 로그 파일
├── init_local.sh               # 로컬 환경 초기화 (Docker + MySQL)
├── start_local.sh              # 로컬 서버 시작
├── stop_local.sh               # 로컬 서버 중지
├── start_test_server.sh        # 테스트 서버 시작
└── copy_production_data.sh     # 로컬 업로드 경로
├── webdata_test/           # 테스트 업로드 경로
├── logs/                   # 로그 파일
├── init_local.sh           # 로컬 환경 초기화
├── start_local.sh          # 로컬 서버 시작
├── start_test_server.sh    # 테스트 서버 시작
└── copy_production_data.sh # 프로덕션 데이터 복사
```

---

## 🔄 DB 초기화 (데이터 리셋)

### 로컬 DB 초기화
# MySQL 컨테이너 삭제 및 재생성
docker-compose -f docker-compose.local.yml down -v
rm local.db
./init_local.sh
```

### 테스트 DB 초기화
```bash
rm test.db
./start_test_server.sh
```

---

## 🐛 디버깅

### 로그 확인
```bash
# 실시간 로그
tail -f logs/local.log

# 전체 로그
cat logs/local.log
```

### DB  (MySQL)
mysql -h127.0.0.1 -P3307 -ulocal -plocal123 repair_local

# 또는 Docker exec
docker exec -it repair_mysql_local mysql -ulocal -plocal123 repair_local

# 테스트 DB (SQLite)
sqlite3 test.db
```

**유용한 SQL:**
```sql
-- MySQL
SHOW TABLES;
SELECT * FROM user_tb;
SELECT * FROM periodic_tb;
exit

-- SQLite
.tables
SELECT * FROM user_tb;T * FROM periodic_tb;

-- 종료
.quit
```

---

## 🔧 문제 해결

### 1. 포트가 이미 사용중
```bash
# 프로세스 확인
lsof -i :9107

# 프로세스 종료
kill -9 <PID>
```
MySQL 컨테이너 상태 확인
docker ps | grep repair_mysql_local

# MySQL 로그 확인
docker logs repair_mysql_local

# MySQL 재시작
docker-compose -f docker-compose.local.yml restart

# 완전히 재생성
docker-compose -f docker-compose.local.yml down -vl.db

# DB 재생성
rm local.db
./init_local.sh
```

### 3. 컴파일 오류
```bash
# 의존성 다운로드
go mod download

# 캐시 정리
go clean -cache
```

### 4. 파일 업로드 오류
```bash
# 디렉토리 권한 확인
ls -la webdata_local/

# 디렉토리 재생성
rm -rf webdata_local
mkdir -p webdata_local/detail
```

---

## 💡 개발 팁

### 1. Hot Reload
```bash
# air 설치 (Go Hot Reload)
go install github.com/cosmtrek/air@latest

# air 실행
MODE=local ENV_FILE=.env.local.yml air
```

### 2. API 테스트
```bash
# JWT 토큰 발급
curl "http://localhost:9107/api/jwt?loginid=dev&passwd=dev123"

# API 호출 (토큰 필요)
curl -H "Authorization: Bearer <token>" \
     http://localhost:9107/api/periodic/search
```

### 3. 여러 환경 동시 실행
```bash
# 터미널 1: 로컬 개발 (9107)
./start_local.sh

# 터미널 2: 테스트 (9109)
./start_test_server.sh

# 각각 독립적으로 작동
```MySQL (Docker) | SQLite | MySQL (공유) | MySQL |
| DB 포트 | 3307 | - | 3306 | 3306 |
| 데이터 | 로컬 전용 | 테스트 전용 | 공유 | 운영 |
| 초기화 | 자유 | 자유 | 불가 | 절대 불가 |
| 용도 | 일상 개발 | 간단 테스트 | 공유 개발 | 운영 |
| 추천도 | ⭐⭐⭐⭐⭐ | 

| 구분 | Local | Test | Develop | Production |
|------|-------|------|---------|------------|
| 포트 | 9107 | 9109 | 9108 | 8080 |
| DB | SQLite | SQLite | MySQL | MySQL |
| 데이터 | 로컬 전용 | 테스트 전용 | 공유 | 운영 |
| 초기화 | 자유 | 자유 | 불가 | 절대 불가 |
| 용도 | 일상 개발 | 통합 테스트 | 공유 개발 | 운영 |
| 추천도 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐ | - |

---

## 🎓 권장 워크플로우

1. **로컬 개발** (Local)
   - 새 기능 개발
   - 버그 수정
   - 실험적 코드

2. **로컬 테스트** (Local)
   - 단위 테스트
   - API 테스트

3. **통합 테스트** (Test)
   - 전체 시나리오 테스트
   - 프로덕션 데이터로 검증
Docker Desktop 설치 및 실행
- [ ] 로컬 환경 초기화 (`./init_local.sh`)
- [ ] 로컬 서버 시작 (`./start_local.sh`)
- [ ] MySQL 접속 확인 (`mysql -h127.0.0.1 -P3308 -ulocal -plocal123 repair_local`)
- [ ] API 테스트
- [ ] 개발 시작!

## 🐳 Docker 명령어 참고

```bash
# 컨테이너 상태 확인
docker ps

# 로그 확인
docker logs repair_mysql_local
docker logs -f repair_mysql_local  # 실시간

# 컨테이너 재시작
docker-compose -f docker-compose.local.yml restart

# 컨테이너 중지
docker-compose -f docker-compose.local.yml stop

# 컨테이너 삭제 (데이터 유지)
docker-compose -f docker-compose.local.yml down

# 컨테이너 + 데이터 완전 삭제
docker-compose -f docker-compose.local.yml down -v
```roduction)

---

## 📝 다음 할 일

- [ ] 로컬 환경 초기화 (`./init_local.sh`)
- [ ] 로컬 서버 시작 (`./start_local.sh`)
- [ ] API 테스트
- [ ] 개발 시작!

---

## 🆘 도움이 필요하면

1. 백엔드 구조 문서: [doc/BACKEND_STRUCTURE.md](doc/BACKEND_STRUCTURE.md)
2. 로그 확인: `tail -f logs/local.log`
3. DB 확인: `sqlite3 local.db`
