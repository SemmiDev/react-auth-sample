package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}

type createUserResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  int64        `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt int64        `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (s *Server) registerUserHandler(c *fiber.Ctx) error {
	timeZone := c.Get("TimeZone")
	log.Println(timeZone)

	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	account, err := NewAccount(req.Username, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errChan := make(chan error, 1)
	go s.dataStore.CreateAccount(c.Context(), errChan, account)

	select {
	case <-time.After(5 * time.Second):
		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
			"error": "request timeout",
		})
	case err := <-errChan:
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					return c.Status(fiber.StatusConflict).JSON(fiber.Map{
						"error": "username or email already exists",
					})
				}
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		accessToken, accessPayload, err := s.tokenMaker.CreateToken(
			account.Username,
			timeZone,
			s.config.AccessTokenDuration,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
			account.Username,
			timeZone,
			s.config.RefreshTokenDuration,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		res := createUserResponse{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessPayload.ExpiredAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
			User: userResponse{
				Username:  account.Username,
				Email:     account.Email,
				CreatedAt: account.CreatedAt,
			},
		}
		return c.Status(fiber.StatusCreated).JSON(res)
	}
}

func (s *Server) loginUserHandler(c *fiber.Ctx) error {
	timeZone := c.Get("TimeZone")
	log.Println(timeZone)

	var req loginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	accountChan, errChan := make(chan *Account, 1), make(chan error, 1)
	go s.dataStore.GetAccountByUsername(c.Context(), errChan, accountChan, req.Username)

	select {
	case <-time.After(5 * time.Second):
		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
			"error": "request timeout",
		})
	case err := <-errChan:
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		account := <-accountChan
		err = account.VerifyPassword(req.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		accessToken, accessPayload, err := s.tokenMaker.CreateToken(
			account.Username,
			timeZone,
			s.config.AccessTokenDuration,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
			account.Username,
			timeZone,
			s.config.RefreshTokenDuration,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		res := createUserResponse{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessPayload.ExpiredAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
			User: userResponse{
				Username:  account.Username,
				Email:     account.Email,
				CreatedAt: account.CreatedAt,
			},
		}
		return c.Status(fiber.StatusCreated).JSON(res)
	}
}
