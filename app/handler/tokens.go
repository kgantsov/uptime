package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/model"
	"github.com/rs/zerolog/log"
)

// ── Input / Output types ──────────────────────────────────────────────────────

type CreateTokenInput struct {
	Body model.CreateToken
}

type CreateTokenOutput struct {
	Body *model.Token
}

type DeleteTokenInput struct {
	Authorization string `header:"Authorization" doc:"Bearer token used to identify the session to invalidate" required:"true"`
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// CreateToken authenticates with email/password and returns a signed JWT.
func (h *Handler) CreateToken(
	ctx context.Context,
	input *CreateTokenInput,
) (*CreateTokenOutput, error) {
	token, err := h.TokenService.CreateToken(input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "email or password is incorrect")
	}

	return &CreateTokenOutput{Body: token}, nil
}

// DeleteToken invalidates the currently authenticated token.
func (h *Handler) DeleteToken(
	ctx context.Context,
	input *DeleteTokenInput,
) (*struct{}, error) {
	tokenID, err := auth.ParseTokenIDFromHeader(input.Authorization, h.JWTKey)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("DeleteToken: could not parse token ID: %s", err)
		return nil, huma.NewError(http.StatusBadRequest, "invalid authorization token")
	}

	log.Ctx(ctx).Info().Msgf("DeleteToken: invalidating token ID %d", tokenID)

	if err := h.TokenService.DeleteToken(tokenID); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "failed to delete token", err)
	}

	return nil, nil
}
