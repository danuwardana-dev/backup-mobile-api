package smtp

import (
	"backend-mobile-api/app/config"
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/enum"
	"context"
	"fmt"
	"net/smtp"
)

type Smtp struct {
	cfg     *config.Root
	clogger *helpers.CustomLogger
}

func NewSmtp(cfg *config.Root, clogger *helpers.CustomLogger) *Smtp {
	return &Smtp{cfg: cfg, clogger: clogger}
}

func (s *Smtp) RegisterOtpMsg(otp string) string {
	body := fmt.Sprintf("Hello,\n\nWe received a request to verify your account. Please use the following code to proceed:\n\n🔑 Your OTP Code: %s\n\nThis code is valid for %v seconds. Do not share this code with anyone to keep your account secure.\n\nIf you did not request this code, please ignore this email.\n\nBest regards,\nBeyondTech", otp, s.cfg.App.OtpExpire)
	return body
}

func (s *Smtp) SendMail(c context.Context, to []string, subjectData enum.EmailSubject, bodyMsg string) error {
	var (
		adress  = fmt.Sprintf("%s:%s", s.cfg.Smtp.Host, s.cfg.Smtp.Port)
		auth    = smtp.PlainAuth("", s.cfg.Smtp.From, s.cfg.Smtp.Password, s.cfg.Smtp.Host)
		from    = s.cfg.Smtp.From
		subject = fmt.Sprintf("Subject: %s\n", subjectData)
		mime    = "MIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n"

		msg = []byte(subject + mime + "\r\n" + bodyMsg)
	)

	err := smtp.SendMail(adress, auth, from, to, msg)

	if err != nil {
		s.clogger.ErrorLogger(c, fmt.Sprintf("SendEmail.smtp.SendMail to %s", to), err)
	}
	return err

}
