package MailFlow

import (
	"context"
	"regexp"
)

var emailRegex = regexp.MustCompile(`[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" +
	`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+([a-z0-9](?:[a-z0-9-]*[a-z0-9])){1,}?`)

// EmailIsValid checks if the email address is valid according to RFC 5321
func EmailIsValid(ctx context.Context, email string) bool {
	ctx, emailValidationSpan := tracer.Start(ctx, "Validate email address")
	defer emailValidationSpan.End()

	emailValid := emailRegex.Match([]byte(email))
	isValid := emailValid && len(email) <= 320
	if !isValid {
		logger.Debug(ctx, "Invalid email address")
	}

	return isValid
}
