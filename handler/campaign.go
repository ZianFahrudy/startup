package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {
		response := helper.APIResponse("Get Campaigns Failed", http.StatusBadRequest, "Failed", campaign.FormatCampaigns(campaigns))
		c.JSON(http.StatusBadRequest, response)
		return

	}

	response := helper.APIResponse("Get Campaigns Success", http.StatusOK, "Success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get Detail Campaign", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	campaignDetail, err := h.service.GetCampaignByID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get Detail Campaign", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to get Detail Campaign", http.StatusOK, "Error", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)

}
