# go-spreadsheet-library

<img width="675" alt="slack" src="https://user-images.githubusercontent.com/5152494/127766277-94699dc5-9f13-477b-8f96-b3c065b05f02.png">

스캐터랩 사내 Go 스터디에서 구현한 사내 도서 대출/반납 관리 시스템입니다. 

Google Spreadsheet를 기반으로 동작하며, oAuth 2.0 인증을 획득하여 표를 수정하는 방식으로 관리합니다. Clean Architecture를 적용한 서버 및 Slackbot의 형태입니다.

## 구현한 기능

* 책의 이름을 기반으로 한 검색과 대출
* 반납 기일에 맞춘 Due Date 알람 및 반납

## 개발 환경 사용 방법

* Go 1.16 이상이 필요합니다.
* 환경 변수를 설정해야 합니다. `.env.example` 파일을 참조하세요.

```bash
$ make (run)    # 개발용 서버 실행
$ make build    # 바이너리 빌드
$ make test     # 테스트 실행
```

## 배포 방법

### Docker Image Build

```bash
$ docker build -t scatterlab-library:vx.y.z .
```

### Docker Run

* `.env`를 넣을 수 있을 경우 `-e` 절을 생략하고 `/.env`에 넣으면 됩니다.
* 무엇을 넣어야 하는지는 `.env.example` 을 참조해주세요.

```bash
$ docker run --rm -it \
    -e ENV_VARIABLE=asdfsadf \ # 대체해야 함
    -p 8080:8080 \
    scatterlab-library:vx.y.z
```

### Readiness / Healthcheck

* Readiness Check / Healthcheck이 필요한 경우, `/`에서 Status 200이 떨어지는 것을 기준으로 해주세요.
