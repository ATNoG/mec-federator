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

func (zisc *ZonesInfoSyncController) SubscribeZoneController(c *gin.Context) {
}

func (zisc *ZonesInfoSyncController) UnsubscribeZoneController(c *gin.Context) {
}

func (zisc *ZonesInfoSyncController) GetZoneController(c *gin.Context) {
}
