package server

import (
	"database/sql"
	"errors"
	"time"

	db "github.com/duyanhitbe/go-boilerplate/internal/database/generated"
	"github.com/duyanhitbe/go-boilerplate/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *Server) register(c *gin.Context) {
	log.Trace().Msg("=========== register ==========")
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Msg("Invalid payload")

		if utils.IsValidationError(err) {
			errs := utils.MakeValidationError(err, req)
			throwValidationError(c, errs)
			return
		}

		throwBadRequest(c, err)
		return
	}

	_, err := s.store.FindOneUserByUsername(c, req.Username)
	if err != nil {
		//Create new user
		if errors.Is(err, sql.ErrNoRows) {
			pwd, err := s.h.Create(req.Password)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create password")
				throwInternalServerError(c, err)
				return
			}
			user, err := s.store.CreateUser(c, db.CreateUserParams{
				Username: req.Username,
				Password: pwd,
			})
			if err != nil {
				log.Error().Msg("Create user failed")
				throwInternalServerError(c, err)
				return
			}

			log.Trace().Msgf("Create user succeeded: %s", user.ID)
			created(c, registerResponse{
				ID:        user.ID,
				Username:  user.Username,
				CreatedAt: user.CreatedAt.Time,
				UpdatedAt: user.UpdatedAt.Time,
			})
		} else {
			log.Error().Msgf("FindOneUserByUsername fail %v", err)
			throwInternalServerError(c, err)
			return
		}
	} else {
		log.Error().Msgf("username already exists: %s", req.Username)
		throwConflict(c, errors.New("username already taken"))
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	ExpiresIn    int64  `json:"expires_in"`
	Unit         string `json:"unit"`
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) login(c *gin.Context) {
	log.Trace().Msg("=========== login ==========")
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Msg("Invalid payload")

		if utils.IsValidationError(err) {
			errs := utils.MakeValidationError(err, req)
			throwValidationError(c, errs)
			return
		}

		throwBadRequest(c, err)
		return
	}
	user, err := s.store.FindOneUserByUsername(c, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Msg("user not found")
			throwForbidden(c, errors.New("user not found"))
			return
		}

		log.Error().Msgf("FindOneUserByUsername fail %v", err)
		throwInternalServerError(c, err)
		return
	}

	if ok := s.h.Compare(req.Password, user.Password); !ok {
		log.Error().Msg("Invalid password")
		throwBadRequest(c, errors.New("invalid password"))
		return
	}

	accessToken, err := s.t.Create(user.ID.String(), 24*time.Hour)
	if err != nil {
		log.Error().Msgf("Create accessToken fail %v", err)
		throwInternalServerError(c, err)
		return
	}
	log.Trace().Msg("Create access token succeeded")

	refreshToken, err := s.t.Create(user.ID.String(), 30*24*time.Hour)
	if err != nil {
		log.Error().Msgf("Create refreshToken fail %v", err)
		throwInternalServerError(c, err)
		return
	}
	log.Trace().Msg("Create refresh token succeeded")

	ok(c, loginResponse{
		ExpiresIn:    accessToken.ExpiresIn,
		Unit:         accessToken.Unit,
		Type:         accessToken.Type,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	})
}

type meResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Server) me(c *gin.Context) {
	if u, valid := s.getUserFromContext(c); valid {
		ok(c, meResponse{
			ID:        u.ID,
			Username:  u.Username,
			CreatedAt: u.CreatedAt.Time,
			UpdatedAt: u.UpdatedAt.Time,
		})
	}
}
