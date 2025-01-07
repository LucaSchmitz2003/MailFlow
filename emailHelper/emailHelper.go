package MailFlow

import (
	"context"
	"github.com/LucaSchmitz2003/FlowWatch"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"os"
	"sync"
)

var (
	tracer = otel.Tracer("MailHelperTracer")
	logger = FlowWatch.GetLogHelper()

	emailSenderInstance *EmailSender
	once                sync.Once
)

// EmailSender is a struct that contains the configuration for sending emails.
type EmailSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

// initEmailSender initializes the EmailSender by reading configuration from environment variables.
// Should be called once at application start.
func initEmailSender(ctx context.Context) (*EmailSender, error) {
	// Start a new span for the database initialization
	ctx, span := tracer.Start(ctx, "Initialize email sender")
	defer span.End()

	// Load the environment variables to make sure that the settings have already been loaded
	_ = godotenv.Load(".env")

	// Read the SMTP configuration from environment variables.
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return nil, errors.New("SMTP_HOST is not set")
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		return nil, errors.New("SMTP_PORT is not set")
	}

	username := os.Getenv("SMTP_USERNAME")
	if username == "" {
		return nil, errors.New("SMTP_USERNAME is not set")
	}

	password := os.Getenv("SMTP_PASSWORD")
	if password == "" {
		return nil, errors.New("SMTP_PASSWORD is not set")
	}

	// Set the email sender instance
	emailSender := &EmailSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}

	logger.Debug(ctx, "Initialized email sender")

	return emailSender, nil
}

// GetEmailSender creates a new email sender instance or returns an already existing instance
func GetEmailSender(ctx context.Context) *EmailSender {
	ctx, span := tracer.Start(ctx, "Get email sender")
	defer span.End()

	// Create a new email sender instance if it does not exist
	once.Do(func() {
		var err error
		emailSenderInstance, err = initEmailSender(ctx)
		if err != nil {
			err = errors.Wrap(err, "Failed to initialize the email sender instance")
			logger.Fatal(ctx, err)
		}
	})

	logger.Debug(ctx, "Got email sender instance")

	return emailSenderInstance
}
