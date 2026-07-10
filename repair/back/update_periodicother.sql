SET NAMES utf8mb4;
-- UTF-8 인코딩으로 저장
SET NAMES utf8mb4;

UPDATE periodicother_tb 
SET po_name = CONCAT('상태양호,', po_name)
WHERE po_type = 2 
  AND po_name NOT LIKE '상태양호%'
  AND po_category IN (10, 11, 12, 13, 14, 15);