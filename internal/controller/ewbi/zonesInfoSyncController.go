package ewbi

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/services"
)

type ZonesInfoSyncController struct {
	zoneService *services.ZoneService
}

func NewZonesInfoSyncController(zoneService *services.ZoneService) *ZonesInfoSyncController {
	return &ZonesInfoSyncController{
		zoneService: zoneService,
	}
}

// @Summary Subscribe to a Zone
// @Description Used by origin OP to show intent on using a partner OP's zone
// @Tags EWBI - ZonesInfoSync
func (zisc *ZonesInfoSyncController) SubscribeZoneController(c *gin.Context) {
}

// @Summary Unsubscribe from a Zone
// @Description Used by origin OP to show intent on not using a partner OP's zone anymore
// @Tags EWBI - ZonesInfoSync
func (zisc *ZonesInfoSyncController) UnsubscribeZoneController(c *gin.Context) {
}

// @Summary Get Zone Details
// @Description Used by origin OP to get details of a zone that belongs to a partner OP
// @Tags EWBI - ZonesInfoSync
func (zisc *ZonesInfoSyncController) GetZoneController(c *gin.Context) {
}
