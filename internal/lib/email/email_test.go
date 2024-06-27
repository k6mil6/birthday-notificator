package email_test

import (
	"github.com/k6mil6/birthday-notificator/internal/config"
	"github.com/k6mil6/birthday-notificator/internal/lib/email"
	"testing"
)

type emailTestCase struct {
	email    string
	expected bool
}

var emailTestCases = []emailTestCase{
	{
		email:    "kamilion843@outlook.com",
		expected: true,
	},
	{
		email:    "kamilion843@outlook.",
		expected: false,
	},
	{
		email:    "kamilion843@outlook.com.",
		expected: false,
	},
	{
		email:    "kamilion843@",
		expected: false,
	},
}

func TestIsValidEmail(t *testing.T) {
	for _, tc := range emailTestCases {
		t.Run(tc.email, func(t *testing.T) {
			result := email.IsValidEmail(tc.email)
			if result != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, result)
			}
		})
	}
}

func TestSender_Send(t *testing.T) {
	cfg := config.MustLoadPath("../../../config/config.yaml")
	sender := email.NewSender(cfg.Email.SenderAddress, cfg.Email.SenderPassword, cfg.Email.SMTPAddress, cfg.Email.SMTPPort)

	err := sender.Send("kamil_6@vk.com", "test", "test")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
