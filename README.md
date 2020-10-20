# Solapi Agent

## 설치

/opt/agent 디렉토리를 만들고 아래로 에이전트 실행파일을 복사합니다.
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
  "Prefix": "",
  "AppId": ""
}
```

## 서비스 데몬 실행
서비스 데몬을 시스템에 등록 및 실행합니다.
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
