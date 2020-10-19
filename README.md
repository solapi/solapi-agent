# Solapi Agent

기본 설치 경로는 /opt/agent 입니다.

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

기본 설치 경로와 다르게 설치한 경우 아래와 같이 AGENT_HOME 환경변수를 설정해 주세요
```
export AGENT_HOME=/home/ubuntu/agent
```
