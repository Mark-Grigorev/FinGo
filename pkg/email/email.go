package email

import (
	"fmt"

	"github.com/resend/resend-go/v3"
)

var htmlPage = `
<div style="font-family:sans-serif;max-width:480px;margin:0 auto">
  <h2>Восстановление пароля</h2>
  <p>Мы получили запрос на сброс пароля для вашего аккаунта FinGo.</p>
  <p>
    <a href="%s" style="display:inline-block;padding:12px 24px;background:#4f46e5;color:#fff;border-radius:6px;text-decoration:none">
      Сбросить пароль
    </a>
  </p>
  <p style="color:#6b7280;font-size:14px">Ссылка действительна 15 минут. Если вы не запрашивали сброс пароля — просто проигнорируйте это письмо.</p>
</div>`

// Sender отправляет письма пользователям.
type Sender interface {
	SendPasswordReset(to, resetURL string) error
}

// ResendSender реализует Sender через Resend API.
type ResendSender struct {
	client *resend.Client
	from   string
}

func NewResend(apiKey string) *ResendSender {
	return &ResendSender{
		client: resend.NewClient(apiKey),
		from:   "onboarding@resend.dev",
	}
}

func (s *ResendSender) SendPasswordReset(to, resetURL string) error {
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: "Восстановление пароля FinGo",
		Html:    fmt.Sprintf(htmlPage, resetURL),
	}
	_, err := s.client.Emails.Send(params)
	return err
}
