# Solapi Agent

## 개요
DB INSERT로 카카오톡 및 문자를 발송 할 수 있도록 go언어로 작성되었습니다.
빌드되어 올려져 있는 agent파일은 Ubuntu 16.04 환경에서 빌드되었으며 다른 버전의 OS에서는 정상적으로 작동되지 않으로 go컴파일러로 새로 빌드하셔야 합니다.

## DB 준비
> MySQL 버전 5.7.14 이상을 준비해주세요.
> 혹은 MariaDB 버전 10.2 이상을 준비해주세요.

아래 내용으로 DB 및 계정을 만들어 주세요.
```
CREATE DATABASE msg;
CREATE USER 'msg'@'localhost' IDENTIFIED BY 'msg';
GRANT ALL PRIVILEGES ON msg.* TO 'msg'@'localhost';
```
아래 스키마로 테이블을 만들어 주세요. (MariaDB는 create_table_maria.sql 참고)
```
CREATE TABLE msg (
  id integer  AUTO_INCREMENT primary key,
  createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  scheduledAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  sendAttempts SMALLINT NOT NULL DEFAULT 0,
  reportAttempts SMALLINT NOT NULL DEFAULT 0,
  `to` VARCHAR(20) AS (payload->>'$.to') STORED,
  `from` VARCHAR(20) AS (payload->>'$.from') STORED,
  groupId VARCHAR(255) AS (result->>'$.groupId') STORED,
  messageId VARCHAR(255) AS (result->>'$.messageId') STORED,
  status VARCHAR(20) AS (result->>'$.status') STORED,
  statusCode VARCHAR(255) AS (result->>'$.statusCode') STORED,
  statusMessage VARCHAR(255) AS (result->>'$.statusMessage') STORED,
  payload JSON,
  result JSON default NULL,
  sent BOOLEAN NOT NULL default false,
  KEY (`createdAt`),
  KEY (`updatedAt`),
  KEY (`scheduledAt`),
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
```

## 소스 코드 빌드
아래 명령으로 빌드하면 agent 실행파일이 생성됩니다.
```
go build agent.go
```

## 서비스 데몬 설치

빌드된 agent 파일을 /opt/agent 디렉토리를 만들고 아래로 복사합니다.
```
mkdir -p /opt/agent
cp ./agent /opt/agent/agent
```

/opt/agent/db.json 파일을 만들고 DB접속 정보를 입력합니다.
```
vi /opt/agent/db.json
```
db.json 예시
```
{
  "provider": "mysql",
  "dbname": "msg",
  "table": "msg",
  "user": "root",
  "password": "root-password",
  "host": "localhost",
  "port": 3306
}
```

/opt/agent/config.json 파일을 만들고 API Key정보를 입력합니다.
```
vi /opt/agent/config.json
```
config.json 예시
```
{
  "APIKey": "NCSVYGF1IK5PUKDA",
  "APISecret": "FSD4ER2WYPZQVDBPKMLOZVAWTGYBDTRW",
  "Protocol": "https",
  "Domain": "api.solapi.com",
  "Prefix": "", // 사용안함
  "AppId": "", // 사용안함
  "AllowDuplicates": true // 동시간대 중복 발송 허용
}
```

## 서비스 데몬 실행
서비스 데몬을 시스템에 등록 및 실행합니다. (모든 명령은 root 권한으로 실행해야 합니다.)
```
./agent install
./agent start
```

기본 설치 경로(/opt/agent)와 다르게 설치한 경우 아래와 같이 AGENT_HOME 환경변수를 설정해 주세요
```
export AGENT_HOME=/home/ubuntu/agent
```

## 서비스 데몬 명령
시스템에 등록
```
./agent install
```

데몬 실행
```
./agent start
```

데몬 상태
```
./agent status
```

데몬 정지
```
./agent stop
```

시스템에서 제거
```
./agent remove
```
