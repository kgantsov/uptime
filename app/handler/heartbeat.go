package handler

import (
	"net/http"
	"strconv"

	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/echo/v4"
)

// GetHeartbeats godoc
// @Summary      Get all heartbeats
// @Description  Returns all heartbeats
// @Tags         heartbeats
// @Accept       json
// @Produce      json
// @Param        service_id    query     string  false  "Filtering by service_id"
// @Success      200  {object}  []model.Heartbeat
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/heartbeats/latencies [get]
func (h *Handler) GetHeartbeatsLatencies(c echo.Context) error {
	serviceID := c.QueryParam("service_id")

	var heartbeats []model.Heartbeat

	var err error
	if serviceID != "" {
		err = h.DB.Order("id desc").Where("service_id = ?", serviceID).Limit(100).Find(&heartbeats).Error
	} else {
		err = h.DB.Order("id desc").Limit(100).Find(&heartbeats).Error
	}

	if err != nil {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: err}
	}

	return c.JSON(http.StatusOK, heartbeats)
}

// GetHeartbeatStats godoc
// @Summary      Getheartbeats stats
// @Description  Returns heartbeats stats
// @Tags         heartbeats
// @Accept       json
// @Produce      json
// @Param        size    query     string  false  "Size"
// @Success      200  {object}  []model.HeartbeatPoint
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/heartbeats/stats [get]
func (h *Handler) GetHeartbeatStats(c echo.Context) error {
	s := c.QueryParam("size")
	if s == "" {
		s = "3"
	}
	size, err := strconv.Atoi(s)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	var heartbeats []model.Heartbeat

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
		return &echo.HTTPError{Code: http.StatusNotFound, Message: err}
	}

	return c.JSON(http.StatusOK, heartbeats)
}
