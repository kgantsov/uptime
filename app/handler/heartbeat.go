package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kgantsov/uptime/app/model"
	"github.com/rs/zerolog/log"
)

// ── Input / Output types ──────────────────────────────────────────────────────

type GetHeartbeatsLatenciesInput struct {
	ServiceID string `query:"service_id" doc:"Filter by service ID"`
	Size      int    `query:"size" default:"100" min:"1" doc:"Maximum number of results to return"`
}

type GetHeartbeatsLatenciesOutput struct {
	Body []model.Heartbeat
}

type GetHeartbeatsLastLatenciesInput struct {
	Size int `query:"size" default:"3" min:"1" doc:"Number of latest entries per service to return"`
}

type GetHeartbeatsLastLatenciesOutput struct {
	Body []model.Heartbeat
}

type GetHeartbeatStatsInput struct {
	Days int `path:"days" doc:"Number of days to include in the stats window"`
}

type GetHeartbeatStatsOutput struct {
	Body []model.HeartbeatStatsPoint
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// GetHeartbeatsLatencies returns heartbeat latency records, optionally filtered
// by service ID and limited to a maximum size.
func (h *Handler) GetHeartbeatsLatencies(
	ctx context.Context,
	input *GetHeartbeatsLatenciesInput,
) (*GetHeartbeatsLatenciesOutput, error) {
	size := input.Size

	heartbeats, err := h.HeartbeatService.GetLatencies(input.ServiceID, size)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "failed to retrieve latencies", err)
	}

	if heartbeats == nil {
		heartbeats = []model.Heartbeat{}
	}

	return &GetHeartbeatsLatenciesOutput{Body: heartbeats}, nil
}

// GetHeartbeatsLastLatencies returns the most-recent latency point per service.
func (h *Handler) GetHeartbeatsLastLatencies(
	ctx context.Context,
	input *GetHeartbeatsLastLatenciesInput,
) (*GetHeartbeatsLastLatenciesOutput, error) {
	size := input.Size

	heartbeats, err := h.HeartbeatService.GetLastLatencies(size)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("error getting latest latencies: %s", err)
		return nil, huma.NewError(http.StatusNotFound, "failed to retrieve last latencies", err)
	}

	if heartbeats == nil {
		heartbeats = []model.Heartbeat{}
	}

	return &GetHeartbeatsLastLatenciesOutput{Body: heartbeats}, nil
}

// GetHeartbeatStats returns aggregated heartbeat statistics for the given
// number of days.
func (h *Handler) GetHeartbeatStats(
	ctx context.Context,
	input *GetHeartbeatStatsInput,
) (*GetHeartbeatStatsOutput, error) {
	stats, err := h.HeartbeatService.GetStats(input.Days)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("error getting latency stats: %s", err)
		return nil, huma.NewError(http.StatusNotFound, "failed to retrieve heartbeat stats", err)
	}

	if stats == nil {
		stats = []model.HeartbeatStatsPoint{}
	}

	return &GetHeartbeatStatsOutput{Body: stats}, nil
}
