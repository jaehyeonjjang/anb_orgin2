-- 부착물 (category 16) 데이터를 기존 모든 periodic에 삽입
-- 이미 category 16이 있는 periodic은 제외

SET NAMES utf8mb4;

INSERT INTO periodicother_tb (po_name, po_type, po_result, po_status, po_position, po_filename, po_offlinefilename, po_change, po_category, po_order, po_periodic, po_date)
SELECT 'a,b,c,d,e', 1, 0, '', '부착물 등', '', '', 0, 16, 160, p.po_periodic, NOW()
FROM (SELECT DISTINCT po_periodic FROM periodicother_tb) p
WHERE p.po_periodic NOT IN (SELECT po_periodic FROM periodicother_tb WHERE po_category = 16);

INSERT INTO periodicother_tb (po_name, po_type, po_result, po_status, po_position, po_filename, po_offlinefilename, po_change, po_category, po_order, po_periodic, po_date)
SELECT '상태양호,콘크리트 균열,콘크리트 박락,앵커 시공,앵커 풀림,앵커 탈락,앵커 부식', 2, 0, '', '앵커 및 브라켓', '', '', 0, 16, 161, p.po_periodic, NOW()
FROM (SELECT DISTINCT po_periodic FROM periodicother_tb) p
WHERE p.po_periodic NOT IN (SELECT po_periodic FROM periodicother_tb WHERE po_category = 16 AND po_position = '앵커 및 브라켓');

INSERT INTO periodicother_tb (po_name, po_type, po_result, po_status, po_position, po_filename, po_offlinefilename, po_change, po_category, po_order, po_periodic, po_date)
SELECT '상태양호,면적 적정성,균열,부식,손상', 2, 0, '', '용접(정착부)', '', '', 0, 16, 162, p.po_periodic, NOW()
FROM (SELECT DISTINCT po_periodic FROM periodicother_tb) p
WHERE p.po_periodic NOT IN (SELECT po_periodic FROM periodicother_tb WHERE po_category = 16 AND po_position = '용접(정착부)');

INSERT INTO periodicother_tb (po_name, po_type, po_result, po_status, po_position, po_filename, po_offlinefilename, po_change, po_category, po_order, po_periodic, po_date)
SELECT '상태양호,철물매립 길이,철물 여장(노출) 길이', 2, 0, '', '매립', '', '', 0, 16, 163, p.po_periodic, NOW()
FROM (SELECT DISTINCT po_periodic FROM periodicother_tb) p
WHERE p.po_periodic NOT IN (SELECT po_periodic FROM periodicother_tb WHERE po_category = 16 AND po_position = '매립');

INSERT INTO periodicother_tb (po_name, po_type, po_result, po_status, po_position, po_filename, po_offlinefilename, po_change, po_category, po_order, po_periodic, po_date)
SELECT '상태양호,풀림,탈락,부재 변형,부식', 2, 0, '', '볼트', '', '', 0, 16, 164, p.po_periodic, NOW()
FROM (SELECT DISTINCT po_periodic FROM periodicother_tb) p
WHERE p.po_periodic NOT IN (SELECT po_periodic FROM periodicother_tb WHERE po_category = 16 AND po_position = '볼트');

INSERT INTO periodicother_tb (po_name, po_type, po_result, po_status, po_position, po_filename, po_offlinefilename, po_change, po_category, po_order, po_periodic, po_date)
SELECT '상태양호,면적 적정성,균열,부식,손상', 2, 0, '', '용접(연결부)', '', '', 0, 16, 165, p.po_periodic, NOW()
FROM (SELECT DISTINCT po_periodic FROM periodicother_tb) p
WHERE p.po_periodic NOT IN (SELECT po_periodic FROM periodicother_tb WHERE po_category = 16 AND po_position = '용접(연결부)');

INSERT INTO periodicother_tb (po_name, po_type, po_result, po_status, po_position, po_filename, po_offlinefilename, po_change, po_category, po_order, po_periodic, po_date)
SELECT '상태양호,손상,꼬임 및 뒤틀림,변형,고정클립 손상', 2, 0, '', '와이어 로프', '', '', 0, 16, 166, p.po_periodic, NOW()
FROM (SELECT DISTINCT po_periodic FROM periodicother_tb) p
WHERE p.po_periodic NOT IN (SELECT po_periodic FROM periodicother_tb WHERE po_category = 16 AND po_position = '와이어 로프');

INSERT INTO periodicother_tb (po_name, po_type, po_result, po_status, po_position, po_filename, po_offlinefilename, po_change, po_category, po_order, po_periodic, po_date)
SELECT '상태양호,면외 방향 기울기,배부름 발생', 2, 0, '', '기울기 및 배부름', '', '', 0, 16, 167, p.po_periodic, NOW()
FROM (SELECT DISTINCT po_periodic FROM periodicother_tb) p
WHERE p.po_periodic NOT IN (SELECT po_periodic FROM periodicother_tb WHERE po_category = 16 AND po_position = '기울기 및 배부름');
