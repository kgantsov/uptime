package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/kgantsov/uptime/app/model"
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
	var err error

	serviceID := c.QueryParam("service_id")
	sizeStr := c.QueryParam("size")

	if sizeStr == "" {
		sizeStr = "100"
	}

	size, err := strconv.Atoi(sizeStr)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	heartbeats := []model.Heartbeat{}
	if serviceID != "" {
		err = h.DB.Order("id desc").Where("service_id = ?", serviceID).Limit(size).Find(&heartbeats).Error
	} else {
		err = h.DB.Order("id desc").Limit(size).Find(&heartbeats).Error
	}

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

	heartbeats := []model.Heartbeat{}

	err = h.DB.Raw(
		`
		SELECT * FROM
		(
			SELECT id, service_id, status, created_at, response_time, status_code,
			ROW_NUMBER() OVER (PARTITION BY service_id Order by created_at DESC) AS size
			FROM heartbeats
		) RNK
		WHERE size <= ?
		`,
		size,
	).Scan(&heartbeats).Error

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

	heartbeatStatsPoints := []model.HeartbeatStatsPoint{}

	err = h.DB.Raw(
		`
		SELECT
			service_id,
			status,
			count(1) as counter,
			avg(response_time) as average_response_time
		FROM heartbeats
		WHERE created_at > DATE('now', ?)
		GROUP BY service_id, status;
		`,
		fmt.Sprintf("-%d day", days),
	).Scan(&heartbeatStatsPoints).Error

	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"RequestID": c.Get(echo.HeaderXRequestID),
		}).Infof("Got an error getting latencies stats %s", err)
		return &echo.HTTPError{Code: http.StatusNotFound, Message: err}
	}

	return c.JSON(http.StatusOK, heartbeatStatsPoints)
}
