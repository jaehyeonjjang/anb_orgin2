#!/bin/bash

echo "🔧 프로덕션 데이터 가져오기"
echo "========================================="
echo ""
echo "이 스크립트는 이미 열린 SSH 터미널에서 실행해야 합니다."
echo ""
echo "1️⃣ SSH 터미널 탭으로 이동"
echo "2️⃣ 아래 명령 실행:"
echo ""
echo "   mysqldump -h localhost -P 3306 -u repair -prepairdb \\"
echo "     --single-transaction --quick repair > /tmp/repair_dump.sql"
echo ""
echo "3️⃣ 덤프 완료 후 이 스크립트를 다시 실행하세요"
echo ""
read -p "덤프 파일이 준비되었습니까? (y/n): " ready

if [ "$ready" != "y" ]; then
    echo "❌ 취소됨"
    exit 1
fi

echo ""
echo "📥 덤프 파일 다운로드 중..."

# SCP로 파일 다운로드 시도
scp -P 2222 root@64.176.224.83:/tmp/repair_dump.sql /tmp/repair_dump.sql

if [ ! -f "/tmp/repair_dump.sql" ]; then
    echo ""
    echo "❌ 다운로드 실패. 수동으로 파일을 복사하세요:"
    echo ""
    echo "SSH 터미널에서:"
    echo "  cat /tmp/repair_dump.sql"
    echo ""
    echo "로컬에서:"
    echo "  cat > /tmp/repair_dump.sql"
    echo "  (덤프 내용 붙여넣기)"
    echo "  Ctrl+D"
    exit 1
fi

echo "📦 로컬 MySQL로 가져오는 중..."
docker exec -i repair_mysql_local mysql -uroot -plocal123 repair_local < /tmp/repair_dump.sql 2>&1 | grep -v "Using a password"

if [ $? -eq 0 ]; then
    echo "✅ 데이터 가져오기 완료!"
    
    # 테이블 개수 확인
    TABLE_COUNT=$(docker exec repair_mysql_local mysql -ulocal -plocal123 repair_local \
        -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'repair_local';" 2>/dev/null | tail -1)
    
    echo "📊 총 $TABLE_COUNT 개 테이블 가져옴"
    
    # 임시 파일 삭제
    rm -f /tmp/repair_dump.sql
    
    echo ""
    echo "🚀 이제 서버를 시작할 수 있습니다:"
    echo "   ./start_local.sh"
else
    echo "❌ 데이터 가져오기 실패"
    exit 1
fi
