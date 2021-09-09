-- 발송 예약 (2021년 3월 26일 12시 20분에 발송 예약)
INSERT INTO msg(scheduledAt, payload) VALUES(
  '2021-03-26 12:20:00',
  json_object(
    'to', '01000000001',
    'from', '020000001',
    'text', '테스트 메시지'
  )
);
