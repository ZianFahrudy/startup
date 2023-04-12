package transaction

import "bwastartup/user"

type GetCampaignTransacstionInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
