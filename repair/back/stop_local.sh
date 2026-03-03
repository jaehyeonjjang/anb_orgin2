#!/bin/bash

echo "🛑 로컬 개발 환경 중지"
echo "========================================="

# 포트 9107에서 실행 중인 프로세스만 종료
PID=$(lsof -ti :9107)

if [ -n "$PID" ]; then
    echo "📍 포트 9107의 서버 종료 (PID: $PID)"
    kill $PID
    sleep 1
    
    # 강제 종료 필요 시
    if lsof -ti :9107 > /dev/null 2>&1; then
        echo "⚠️  강제 종료 중..."
        kill -9 $(lsof -ti :9107)
    fi
    echo "✅ 서버 종료 완료"
else
    echo "ℹ️  포트 9107에서 실행 중인 서버가 없습니다"
fi

# MySQL 컨테이너 중지 옵션
echo ""
read -p "MySQL 컨테이너도 중지하시겠습니까? (y/n): " stop_mysql

if [ "$stop_mysql" = "y" ]; then
    echo "🐳 MySQL 컨테이너 중지..."
    docker-compose -f docker-compose.local.yml stop
    echo "✅ MySQL 중지 완료"
    echo ""
    echo "💡 완전히 삭제하려면:"
    echo "   docker-compose -f docker-compose.local.yml down -v"
else
    echo "✅ 백엔드만 중지. MySQL은 계속 실행 중"
fi

echo ""
echo "✅ 완료!"
