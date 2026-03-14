# monitorchi

torchi 서비스 헬스체크 모니터링 도구. 주기적으로 상태를 확인하고 장애/복구 시 Gmail로 알림을 보낸다.

## 구조

```bash
monitorchi/
├── main.go
├── internal/
│   ├── config/      # 설정 로드 (.env 파일 + 환경변수)
│   ├── checker/     # 헬스체크 요청 및 판단
│   └── notifier/    # Gmail SMTP 알림 발송
├── .env.example
└── go.mod
```

## 설정

`.env.example`을 복사해서 `.env`로 만들고 값을 채운다. 환경변수가 이미 설정되어 있으면 `.env`보다 우선한다.

```bash
cp .env.example .env
```

> `SMTP_PASS`는 Gmail 비밀번호가 아니라 [앱 비밀번호](https://myaccount.google.com/apppasswords)를 사용해야 한다 (2FA 필요).

## 빌드 & 실행

```bash
go build -o monitorchi .
./monitorchi
```

## 장애 판단 기준

- HTTP 응답 코드가 200이 아닌 경우
- 응답 JSON의 `status`가 `"ok"`가 아닌 경우
- 요청 자체가 타임아웃/실패한 경우

## 알림 동작

- 정상 → 장애: **[torchi] Service Down** 메일 발송
- 장애 → 정상: **[torchi] Service Recovered** 메일 발송
- 장애 지속 중에는 중복 알림을 보내지 않음
