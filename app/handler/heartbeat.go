package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// GetHeartbeats godoc
// @Summary      Get all heartbeats
// @Description  Returns all heartbeats
// @Tags         heartbeats
// @Accept       json
// @Produce      json
// @Param        service_id    query     string  false  "Filtering by service_id"
// @Param        size    query     string  false  "Size"
// @Success      200  {object}  []model.Heartbeat
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     HttpBearer
// @Router       /API/v1/heartbeats/latencies [get]
func (h *Handler) GetHeartbeatsLatencies(c echo.Context) error {
	serviceID := c.QueryParam("service_id")
	sizeStr := c.QueryParam("size")

	if sizeStr == "" {
		sizeStr = "100"
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	heartbeats, err := h.HeartbeatService.GetLatencies(serviceID, size)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: err}
	}

	return c.JSON(http.StatusOK, heartbeats)
}

// GetHeartbeatsLastLatencies godoc
// @Summary      GetHeartbeatsLastLatencies stats
// @Description  Returns last latencies
// @Tags         heartbeats
// @Accept       json
// @Produce      json
// @Param        size    query     string  false  "Size"
// @Success      200  {object}  []model.HeartbeatPoint
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     HttpBearer
// @Router       /API/v1/heartbeats/latencies/last [get]
func (h *Handler) GetHeartbeatsLastLatencies(c echo.Context) error {
	s := c.QueryParam("size")
	if s == "" {
		s = "3"
	}

	size, err := strconv.Atoi(s)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	heartbeats, err := h.HeartbeatService.GetLastLatencies(size)
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"RequestID": c.Get(echo.HeaderXRequestID),
		}).Infof("Got an error getting latest latencies %s", err)

		return &echo.HTTPError{Code: http.StatusNotFound, Message: err}
	}

	return c.JSON(http.StatusOK, heartbeats)
}

// GetHeartbeatStats godoc
// @Summary      GetHeartbeatStats stats
// @Description  Returns heartbeats stats
// @Tags         heartbeats
// @Accept       json
// @Produce      json
// @Param        days    path     int  true  "Number of days to get stats for"
// @Success      200  {object}  []model.HeartbeatStatsPoint
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     HttpBearer
// @Router       /API/v1/heartbeats/stats/{days} [get]
func (h *Handler) GetHeartbeatStats(c echo.Context) error {
	_days := c.Param("days")
	days, err := strconv.Atoi(_days)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	heartbeatStatsPoints, err := h.HeartbeatService.GetStats(days)
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"RequestID": c.Get(echo.HeaderXRequestID),
		}).Infof("Got an error getting latencies stats %s", err)

		return &echo.HTTPError{Code: http.StatusNotFound, Message: err}
	}

	return c.JSON(http.StatusOK, heartbeatStatsPoints)
}
