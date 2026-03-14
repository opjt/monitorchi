package main

import (
	"fmt"
	"log"
	"time"

	"monitorchi/internal/checker"
	"monitorchi/internal/config"
	"monitorchi/internal/notifier"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	chk := checker.New(cfg.HealthURL, cfg.Timeout)
	mail := notifier.NewMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.MailFrom, cfg.MailTo)

	log.Printf("monitorchi started: checking %s every %s", cfg.HealthURL, cfg.CheckInterval)

	wasDown := false

	for {
		err := chk.Check()
		now := time.Now().Format("2006-01-02 15:04:05")

		if err != nil {
			log.Printf("[FAIL] %s - %v", now, err)

			if !wasDown {
				subject := "[torchi] Service Down"
				body := fmt.Sprintf("Time: %s\nError: %s\nURL: %s", now, err, cfg.HealthURL)

				if mailErr := mail.Send(subject, body); mailErr != nil {
					log.Printf("[MAIL ERROR] %v", mailErr)
				} else {
					log.Printf("[MAIL SENT] down alert")
				}
				wasDown = true
			}
		} else {
			log.Printf("[OK] %s", now)

			if wasDown {
				subject := "[torchi] Service Recovered"
				body := fmt.Sprintf("Time: %s\nURL: %s\nService is back to normal.", now, cfg.HealthURL)

				if mailErr := mail.Send(subject, body); mailErr != nil {
					log.Printf("[MAIL ERROR] %v", mailErr)
				} else {
					log.Printf("[MAIL SENT] recovery alert")
				}
				wasDown = false
			}
		}

		time.Sleep(cfg.CheckInterval)
	}
}
