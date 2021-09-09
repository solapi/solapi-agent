-- 카카오 알림톡 발송 (변수값 입력 방식)
INSERT INTO msg(payload) VALUES(json_object(
  'to', '01000000001',
  'from', '0200000001',
  'kakaoOptions', json_object(
    'pfId', 'KA01PF1903260033550428GGGGGGGGGG', -- 카카오톡 채널의 아이디
    'templateId', 'KA01TP1903260033550428BBBBBBBBBB', -- 알림톡 템플릿의 아이디
    -- 변수명 / 변수값
    'variables', json_object(
      '#{변수1}', '변수값 1'
      '#{변수2}', '변수값 2'
      '#{버튼변수}', '변수값'
    )
  )
));

-- 카카오 알림톡 발송 (직접 입력 방식 ** 주의** 승인된 템플릿 내용과 100% 일치해야 합니다)
INSERT INTO msg(payload) VALUES(json_object(
  'to', '01000000001',
  'from', '020000001',
  'text', '홍길동님 가입을 환영합니다.',
  'subject', '대체 발송시 LMS 제목',
  'kakaoOptions', json_object(
    'pfId', 'KA01PF1903260033550428GGGGGGGGGG',
    'templateId', 'KA01TP1903260033550428BBBBBBBBBB',
    'buttons', json_array(json_object(
      'buttonName', '홈페이지',
      'buttonType', 'WL',
      'linkPc', 'https://www.example.com',
      'linkMo', 'https://m.example.com'
    ), json_object(
      'buttonName', '앱 링크',
      'buttonType', 'AL',
      'linkIos', 'iosscheme://',
      'linkAnd', 'androidscheme://'
    ))
  )
));

-- 해외 카카오 알림톡 발송
INSERT INTO msg(payload) VALUES(json_object(
  'country', '1', -- 발송할 국가의 국가번호 입력 (1: 미국/캐나다, 81: 일본, 61: 호주)
  'to', '01000000001',
  'from', '020000001',
  'text', '홍길동님 가입을 환영합니다.',
  'subject', '대체 발송시 LMS 제목',
  'kakaoOptions', json_object(
    'pfId', 'KA01PF1903260033550428GGGGGGGGGG', -- 카카오톡 채널의 아이디
    'templateId', 'KA01TP1903260033550428BBBBBBBBBB', -- 알림톡 템플릿의 아이디
    -- 변수명 / 변수값
    'variables', json_object(
      '#{변수1}', '변수값 1'
      '#{변수2}', '변수값 2'
      '#{버튼변수}', '변수값'
    )
  )
));
