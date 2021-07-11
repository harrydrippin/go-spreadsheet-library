# 스캐터랩 도서관

[21/2Q Go 스터디](https://www.notion.so/mlpingpong/3-17ec6ad241b3464cac94dfa421a78741)에서 구현된 도서 대출/반납 관리 시스템입니다.

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

```bash
$ docker run --rm -it \
    -e GOOGLE_OAUTH_CLIENT_ID=... \
    -e GOOGLE_OAUTH_CLIENT_SECRET=... \
    -e GOOGLE_SPREADSHEET_ID=... \
    -e SLACK_TOKEN=... \
    -p 8080:8080 \
    scatterlab-library:vx.y.z
```

### Readiness / Health Check

* Readiness / Health Check이 필요한 경우, `/`에서 Status 200이 떨어지는 것을 기준으로 해주세요.