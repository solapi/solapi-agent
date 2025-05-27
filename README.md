# Solapi Agent

## 🛑 주의사항 🛑
* 더 이상 솔라피 DB 연동형 에이전트는 유지보수 되지 않습니다.
* 문서를 통해 [API 개발연동](https://developers.solapi.com/references/messages/sendManyDetail)을 진행해주시거나, [SOLAPI Github](https://github.com/solapi) 내 SDK를 통해 개발연동을 진행해주시기 바랍니다.

## 개요
DB INSERT로 카카오톡 및 문자를 발송 할 수 있도록 go언어로 작성되었습니다.  
현재 Agent는 go 컴파일러(go 1.18 기준)로 새로 빌드하셔야 합니다.

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
go build ./cmd/agent/agent.go
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
  "APIKey": "발급받은 API Key 입력",
  "APISecret": "발급받은 API Secret Key 입력",
  "Protocol": "https",
  "Domain": "api.solapi.com",
  "Prefix": "",
  "AppId": "",
  "AllowDuplicates": true
}
```
Prefix, AppId는 비워두시면 됩니다.
AllowDuplicates값이 true이면 동시간대 같은 수신번호로 여러건 발송이 가능하고, false이면 한건만 발송되고 나머지는 중복건으로 발송 실패됩니다.

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
