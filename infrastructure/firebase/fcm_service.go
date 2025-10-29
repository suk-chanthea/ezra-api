package firebase

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/suk-chanthea/ezra/domain/repository"
	"google.golang.org/api/option"
)

type FCMService interface {
	SendNotification(ctx context.Context, tokens []string, title, body string, data map[string]string) error
	SendToUser(ctx context.Context, userID uint, title, body string, data map[string]string) error
	SendToBand(ctx context.Context, bandID uint, title, body string, data map[string]string) error
	SendToAllExcept(ctx context.Context, excludeUserID uint, title, body string, data map[string]string) error
}

type fcmService struct {
	client    *messaging.Client
	tokenRepo repository.DeviceTokenRepository
}

// NewFCMService creates a new FCM service
// If credentialsPath is empty, it will be skipped (for development without Firebase)
func NewFCMService(credentialsPath string, tokenRepo repository.DeviceTokenRepository) (FCMService, error) {
	// If no credentials path provided, return a dummy service
	if credentialsPath == "" {
		log.Println("⚠️  Firebase credentials not provided. Push notifications will be disabled.")
		return &dummyFCMService{}, nil
	}

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting messaging client: %v", err)
	}

	log.Println("✅ Firebase Cloud Messaging initialized successfully")
	return &fcmService{
		client:    client,
		tokenRepo: tokenRepo,
	}, nil
}

// SendNotification sends a push notification to multiple tokens
func (s *fcmService) SendNotification(ctx context.Context, tokens []string, title, body string, data map[string]string) error {
	if len(tokens) == 0 {
		log.Println("📭 No tokens to send notification to")
		return nil
	}

	// FCM has a limit of 500 tokens per request, batch if needed
	batchSize := 500
	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}
		batch := tokens[i:end]

		if err := s.sendBatch(ctx, batch, title, body, data); err != nil {
			log.Printf("❌ Error sending batch: %v", err)
		}
	}

	return nil
}

func (s *fcmService) sendBatch(ctx context.Context, tokens []string, title, body string, data map[string]string) error {
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:   data,
		Tokens: tokens,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Sound:       "default",
				ChannelID:   "default",
				Priority:    messaging.PriorityHigh,
				ClickAction: "FLUTTER_NOTIFICATION_CLICK",
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
					ContentAvailable: true,
				},
			},
		},
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: title,
				Body:  body,
				Icon:  "/icon.png", // You can customize this
			},
		},
	}

	response, err := s.client.SendMulticast(ctx, message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	// Handle invalid tokens
	if response.FailureCount > 0 {
		var invalidTokens []string
		for idx, resp := range response.Responses {
			if !resp.Success {
				log.Printf("❌ Error sending to token: %v", resp.Error)
				// If token is invalid, mark it for deletion
				if messaging.IsRegistrationTokenNotRegistered(resp.Error) ||
					messaging.IsInvalidArgument(resp.Error) {
					invalidTokens = append(invalidTokens, tokens[idx])
				}
			}
		}

		// Remove invalid tokens from database
		if len(invalidTokens) > 0 {
			log.Printf("🗑️  Removing %d invalid tokens", len(invalidTokens))
			if err := s.tokenRepo.DeleteTokens(ctx, invalidTokens); err != nil {
				log.Printf("⚠️  Error deleting invalid tokens: %v", err)
			}
		}
	}

	log.Printf("✅ FCM: Sent %d/%d messages successfully", response.SuccessCount, len(tokens))
	return nil
}

// SendToUser sends notification to all devices of a specific user
func (s *fcmService) SendToUser(ctx context.Context, userID uint, title, body string, data map[string]string) error {
	tokens, err := s.tokenRepo.GetActiveTokensByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user tokens: %v", err)
	}

	if len(tokens) == 0 {
		log.Printf("📭 No active tokens found for user %d", userID)
		return nil
	}

	log.Printf("📤 Sending notification to user %d (%d devices)", userID, len(tokens))
	return s.SendNotification(ctx, tokens, title, body, data)
}

// SendToBand sends notification to all members of a band
func (s *fcmService) SendToBand(ctx context.Context, bandID uint, title, body string, data map[string]string) error {
	tokens, err := s.tokenRepo.GetTokensByBandID(ctx, bandID)
	if err != nil {
		return fmt.Errorf("error getting band tokens: %v", err)
	}

	if len(tokens) == 0 {
		log.Printf("📭 No active tokens found for band %d", bandID)
		return nil
	}

	log.Printf("📤 Sending notification to band %d (%d devices)", bandID, len(tokens))
	return s.SendNotification(ctx, tokens, title, body, data)
}

// SendToAllExcept sends notification to all users except the specified user
func (s *fcmService) SendToAllExcept(ctx context.Context, excludeUserID uint, title, body string, data map[string]string) error {
	tokens, err := s.tokenRepo.GetAllActiveTokensExcept(ctx, excludeUserID)
	if err != nil {
		return fmt.Errorf("error getting tokens: %v", err)
	}

	if len(tokens) == 0 {
		log.Println("📭 No active tokens found")
		return nil
	}

	log.Printf("📤 Broadcasting notification to all users (%d devices, excluding user %d)", len(tokens), excludeUserID)
	return s.SendNotification(ctx, tokens, title, body, data)
}

// dummyFCMService is used when Firebase is not configured
type dummyFCMService struct{}

func (d *dummyFCMService) SendNotification(ctx context.Context, tokens []string, title, body string, data map[string]string) error {
	log.Printf("🔕 [DUMMY FCM] Would send to %d tokens: %s - %s", len(tokens), title, body)
	return nil
}

func (d *dummyFCMService) SendToUser(ctx context.Context, userID uint, title, body string, data map[string]string) error {
	log.Printf("🔕 [DUMMY FCM] Would send to user %d: %s - %s", userID, title, body)
	return nil
}

func (d *dummyFCMService) SendToBand(ctx context.Context, bandID uint, title, body string, data map[string]string) error {
	log.Printf("🔕 [DUMMY FCM] Would send to band %d: %s - %s", bandID, title, body)
	return nil
}

func (d *dummyFCMService) SendToAllExcept(ctx context.Context, excludeUserID uint, title, body string, data map[string]string) error {
	log.Printf("🔕 [DUMMY FCM] Would broadcast (except user %d): %s - %s", excludeUserID, title, body)
	return nil
}

