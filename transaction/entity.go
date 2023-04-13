package transaction

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	PaymentURL string
	Campaign   campaign.Campaign
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       user.User
}
