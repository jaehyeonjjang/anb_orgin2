#!/bin/bash

echo "🔧 로컬 개발 환경 초기화 (MySQL)"
echo "========================================="

# Docker 확인
if ! command -v docker &> /dev/null; then
    echo "❌ Docker가 설치되어 있지 않습니다."
    echo "   https://www.docker.com/products/docker-desktop"
    exit 1
fi

# 기존 컨테이너 중지 및 제거
echo "🧹 기존 컨테이너 정리..."
docker-compose -f docker-compose.local.yml down -v

# MySQL 컨테이너 시작
echo "🐳 MySQL 컨테이너 시작..."
docker-compose -f docker-compose.local.yml up -d

# MySQL이 준비될 때까지 대기
echo "⏳ MySQL 초기화 대기 중..."
for i in {1..30}; do
    if docker exec repair_mysql_local mysqladmin ping -h localhost -uroot -plocal123 --silent &> /dev/null; then
        echo "✅ MySQL 준비 완료!"
        break
    fi
    echo -n "."
    sleep 1
done

# local 사용자 생성 및 권한 부여
echo "👤 사용자 권한 설정..."
docker exec repair_mysql_local mysql -uroot -plocal123 -e "
    GRANT ALL PRIVILEGES ON repair_local.* TO 'local'@'%';
    FLUSH PRIVILEGES;
" 2>/dev/null

# 초기 데이터 선택
echo ""
echo "초기 데이터를 어떻게 설정하시겠습니까?"
echo "1) 최소 테스트 데이터만 생성 (빠름)"
echo "2) 프로덕션 데이터 복사 (느림)"
echo ""
read -p "선택 (1 or 2): " choice

if [ "$choice" = "2" ]; then
    echo ""
    echo "📥 프로덕션 데이터 복사 중..."
    
    # 프로덕션 MySQL 연결 정보 (Vultr 서버 - 외부 접속용)
    PROD_HOST="64.176.224.83"
    PROD_PORT="3306"
    PROD_DB="repair"
    PROD_USER="repair"
    PROD_PASS="repairdb"
    
    # 로컬 MySQL 연결 정보
    LOCAL_HOST="localhost"
    LOCAL_PORT="3308"
    LOCAL_DB="repair_local"
    LOCAL_USER="local"
    LOCAL_PASS="local123"
    
    echo "📦 전체 테이블 복사 중... (약 75개 테이블)"
    echo "⏳ 시간이 걸릴 수 있습니다 (5~15분)"
    
    # SSH를 통해 프로덕션 서버에서 덤프 (MySQL 포트가 외부에서 막혀있음)
    echo "🔐 SSH로 프로덕션 서버 접속 중..."
    ssh -p 2222 root@64.176.224.83 \
        "mysqldump -h localhost -P 3306 -u repair -prepairdb --single-transaction --quick repair 2>/dev/null" | \
    docker exec -i repair_mysql_local mysql -uroot -plocal123 "$LOCAL_DB" 2>&1 | grep -v "Using a password"
    
    if [ $? -eq 0 ]; then
        echo "✅ 전체 데이터 복사 완료!"
    else
        echo "❌ 데이터 복사 실패"
        exit 1
    fi
    
else
    echo ""
    echo "📝 최소 테스트 데이터 생성..."
    
    # 최소 테스트 데이터
    docker exec -i repair_mysql_local mysql -ulocal -plocal123 repair_local <<EOF
-- 개발용 회사
INSERT INTO company_tb (c_id, c_name, c_hp, c_status, c_date) 
VALUES (1, '로컬개발회사', '010-0000-0000', 1, NOW());

-- 개발용 사용자 (비밀번호: dev123)
-- SHA256 해시: ecd71870d1963316a97e3ac3408c9835ad8cf0f3c1bc703527c30265534f75ae
INSERT INTO user_tb (u_id, u_loginid, u_passwd, u_name, u_level, u_hp, u_email, u_status, u_company, u_date) 
VALUES (1, 'dev', 'ecd71870d1963316a97e3ac3408c9835ad8cf0f3c1bc703527c30265534f75ae', '로컬개발자', 1, '010-0000-0000', 'dev@local.dev', 1, 1, NOW());

-- 샘플 아파트 그룹
INSERT INTO aptgroup_tb (ag_id, ag_company, ag_name, ag_facility, ag_type, ag_status, ag_user, ag_updateuser, ag_date, ag_imagecategory) 
VALUES (1, 1, '샘플 아파트', '샘플 시설', 1, 1, 1, 1, NOW(), '');

-- 샘플 정기점검
INSERT INTO periodic_tb (p_id, p_aptgroup, p_name, p_status, p_date, p_user, p_company) 
VALUES (1, 1, '2024년 샘플 정기점검', 1, NOW(), 1, 1);
EOF
fi

# webdata_local 디렉토리 생성
mkdir -p webdata_local
mkdir -p webdata_local/detail
mkdir -p logs

echo ""
echo "✅ 로컬 개발 환경 준비 완료!"
echo "========================================="
echo "📊 데이터베이스: MySQL (localhost:3308)"
echo "📁 업로드 경로: webdata_local/"
echo "📋 로그 경로: logs/local.log"
echo ""

if [ "$choice" = "1" ]; then
    echo "👤 개발용 계정:"
    echo "   ID: dev"
    echo "   PW: dev123"
fi

echo ""
echo "🔧 MySQL 직접 접속:"
echo "   mysql -h127.0.0.1 -P3308 -ulocal -plocal123 repair_local"
echo ""
echo "🚀 로컬 개발 서버를 시작하려면:"
echo "   ./start_local.sh"
echo ""
