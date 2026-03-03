#!/bin/bash

# Docker 경로 설정
export PATH="/Applications/Docker.app/Contents/Resources/bin:$PATH"

echo "🚀 로컬 개발 서버 시작 (포트: 9107)"
echo "========================================="

# Docker 컨테이너 확인
if ! docker ps | grep -q repair_mysql_local; then
    echo "⚠️  MySQL 컨테이너가 실행 중이지 않습니다."
    echo ""
    read -p "MySQL을 시작하시겠습니까? (y/n): " start_mysql
    if [ "$start_mysql" = "y" ]; then
        echo "🐳 MySQL 시작 중..."
        docker-compose -f docker-compose.local.yml up -d
        echo "⏳ MySQL 준비 대기..."
        sleep 5
    else
        echo "❌ MySQL이 필요합니다. 먼저 초기화를 실행하세요:"
        echo "   ./init_local.sh"
        exit 1
    fi
fi

# MySQL 연결 확인
echo "🔍 MySQL 연결 확인..."
if ! docker exec repair_mysql_local mysqladmin ping -h localhost -ulocal -plocal123 --silent 2>/dev/null; then
    echo "❌ MySQL 연결 실패. 초기화를 실행하세요:"
    echo "   ./init_local.sh"
    exit 1
fi

# 디렉토리 확인
mkdir -p webdata_local
mkdir -p webdata_local/detail
mkdir -p logs

echo ""
echo "✅ 로컬 개발 서버 준비 완료!"
echo "========================================="
echo "📍 접속 주소: http://localhost:9107"
echo "🐳 데이터베이스: MySQL (localhost:3308)"
echo "📁 업로드 경로: webdata_local/"
echo "📋 로그: logs/local.log"
echo "========================================="
echo "종료: Ctrl+C"
echo ""

# develop 모드로 실행 (환경변수 불필요)
go run main.go
