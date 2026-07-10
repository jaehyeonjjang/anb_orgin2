-- 깨진 한글 '상태양호,' 부분 제거 (첫 번째 콤마까지)
UPDATE periodicother_tb 
SET po_name = SUBSTRING(po_name, LOCATE(',', po_name) + 1)
WHERE po_type = 2 
  AND po_category IN (10, 11, 12, 13, 14, 15)
  AND po_name LIKE '%,%';
