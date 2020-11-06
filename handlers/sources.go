package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ksensehq/eventnative/logging"
	"github.com/ksensehq/eventnative/middleware"
	"github.com/ksensehq/eventnative/sources"
	"net/http"
)

type SourceSyncStatusResponse struct {
	Statuses []SourceSyncStatus `json:"statuses"`
}

type SourceSyncStatus struct {
	Collection string `json:"collection"`
	Status     string `json:"status"`
	Logs       string `json:"logs"`
}

type SourcesHandler struct {
	sourcesService *sources.Service
}

func NewSourcesHandler(sourcesService *sources.Service) *SourcesHandler {
	return &SourcesHandler{sourcesService: sourcesService}
}

func (sh *SourcesHandler) SyncHandler(c *gin.Context) {
	sourceId := c.Param("id")
	if sourceId == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{Message: "id is required path parameter"})
		return
	}

	err := sh.sourcesService.Sync(sourceId)
	if err != nil {
		logging.Error(err)
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{Message: "Sync failed: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (sh *SourcesHandler) StatusHandler(c *gin.Context) {
	sourceId := c.Param("id")
	if sourceId == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{Message: "id is required path parameter"})
		return
	}

	statusesMap, err := sh.sourcesService.GetStatus(sourceId)
	if err != nil {
		logging.Error(err)
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{Message: "Getting statuses failed:" + err.Error()})
		return
	}

	logsMap, err := sh.sourcesService.GetLogs(sourceId)
	if err != nil {
		logging.Error(err)
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{Message: "Getting statuses failed:" + err.Error()})
		return
	}

	var statuses []SourceSyncStatus
	for collection, status := range statusesMap {
		if status == "" {
			status = "NEW"
		}
		logs := ""
		actualLogs, ok := logsMap[collection]
		if ok {
			logs = actualLogs
		}
		statuses = append(statuses, SourceSyncStatus{
			Collection: collection,
			Status:     status,
			Logs:       logs,
		})
	}

	c.JSON(http.StatusOK, SourceSyncStatusResponse{Statuses: statuses})
}