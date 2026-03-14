package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"monitorchi/internal/checker"
	"monitorchi/internal/config"
	"monitorchi/internal/notifier"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("v", false, "print version and exit")
	envFile := flag.String("config", ".env", "path to .env config file")
	flag.Parse()

	if *showVersion {
		fmt.Println("monitorchi", version)
		return
	}

	cfg, err := config.Load(*envFile)
	if err != nil {
		log.Fatal(err)
	}

	chk := checker.New(cfg.HealthURL, cfg.Timeout)
	mail := notifier.NewMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.MailFrom, cfg.MailTo)
	push := notifier.NewPusher(cfg.TorchiToken)

	log.Printf("monitorchi started: checking %s every %s", cfg.HealthURL, cfg.CheckInterval)

	wasDown := false
	lastHeartbeat := time.Time{}

	for {
		now := time.Now()
		nowStr := now.Format("2006-01-02 15:04:05")

		// heartbeat: 매일 지정 시각에 한 번
		if now.Hour() == cfg.HeartbeatHour && now.YearDay() != lastHeartbeat.YearDay() {
			msg := fmt.Sprintf("[monitorchi] heartbeat - %s", nowStr)
			if err := push.Send(msg); err != nil {
				log.Printf("[HEARTBEAT ERROR] %v", err)
			} else {
				log.Printf("[HEARTBEAT] sent")
				lastHeartbeat = now
			}
		}

		// health check
		err := chk.Check()

		if err != nil {
			log.Printf("[FAIL] %s - %v", nowStr, err)

			if !wasDown {
				// 장애: 메일 + 푸시
				subject := "[torchi] Service Down"
				body := fmt.Sprintf("Time: %s\nError: %s\nURL: %s", nowStr, err, cfg.HealthURL)

				if mailErr := mail.Send(subject, body); mailErr != nil {
					log.Printf("[MAIL ERROR] %v", mailErr)
				} else {
					log.Printf("[MAIL SENT] down alert")
				}

				pushMsg := fmt.Sprintf("[torchi] Service Down\nTime: %s\nError: %s", nowStr, err)
				if pushErr := push.Send(pushMsg); pushErr != nil {
					log.Printf("[PUSH ERROR] %v", pushErr)
				} else {
					log.Printf("[PUSH SENT] down alert")
				}

				wasDown = true
			}
		} else {
			log.Printf("[OK] %s", nowStr)

			if wasDown {
				// 복구: 푸시만
				pushMsg := fmt.Sprintf("[torchi] Service Recovered\nTime: %s\nService is back to normal.", nowStr)
				if pushErr := push.Send(pushMsg); pushErr != nil {
					log.Printf("[PUSH ERROR] %v", pushErr)
				} else {
					log.Printf("[PUSH SENT] recovery alert")
				}

				wasDown = false
			}
		}

		time.Sleep(cfg.CheckInterval)
	}
}
