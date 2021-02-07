CREATE TABLE msg (
  id integer  AUTO_INCREMENT primary key,
  createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  sendAttempts SMALLINT DEFAULT 0,
  reportAttempts SMALLINT DEFAULT 0,
  `to` VARCHAR(20) AS (payload->>'$.to') STORED,
  `from` VARCHAR(20) AS (payload->>'$.from') STORED,
  groupId VARCHAR(255) AS (result->>'$.groupId') STORED,
  messageId VARCHAR(255) AS (result->>'$.messageId') STORED,
  status VARCHAR(20) AS (result->>'$.status') STORED,
  statusCode VARCHAR(255) AS (result->>'$.statusCode') STORED,
  statusMessage VARCHAR(255) AS (result->>'$.statusMessage') STORED,
  payload JSON,
  result JSON default NULL,
  sent BOOLEAN default false,
  KEY (`createdAt`),
  KEY (`updatedAt`),
  KEY (`sendAttempts`),
  KEY (`reportAttempts`),
  KEY (`to`),
  KEY (`from`),
  KEY (groupId),
  KEY (messageId),
  KEY (status),
  KEY (statusCode),
  KEY (sent)
) DEFAULT CHARSET=utf8mb4;
