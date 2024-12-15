package emailHelper

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/smtp"
)

// SendEmail sends an email to the specified email address.
func (emailSender *EmailSender) SendEmail(ctx context.Context, recipient string, subject, body string) error {
	ctx, emailSendingSpan := tracer.Start(ctx, "Sending email")
	defer emailSendingSpan.End()

	// Authenticate with the SMTP server.
	auth := smtp.PlainAuth("", emailSender.Username, emailSender.Password, emailSender.Host)

	// Compose the message.
	message := []byte(
		"From: " + emailSender.Username + "\r\n" + // The sender of the email.
			"To: " + recipient + "\r\n" + // The recipient of the email.
			"Subject: " + subject + "\r\n" + // The subject of the email.
			"\r\n" + // Empty line to separate the header from the body.
			body + "\r\n") // The body of the email.

	// Put the recipient in a slice because the SendMail function expects a slice of recipients.
	recipients := []string{recipient}

	// Send the email.
	err := smtp.SendMail(fmt.Sprintf("%s:%s", emailSender.Host, emailSender.Port),
		auth, emailSender.Username, recipients, message)
	if err != nil {
		err = errors.Wrap(err, "Failed to send the email")
		logger.Error(ctx, err)
		return err
	}

	return nil
}
