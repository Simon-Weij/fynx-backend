/*
 * Wayland recorder is a way to easily make clips and share them.
 * Copyright (C) 2026 Simon-Weij
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Middleware(ctx fiber.Ctx) error {
	tokenString, err := extractToken(ctx)
	if err != nil {
		return err
	}

	token, err := parseToken(tokenString)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	return validateClaims(ctx, token)
}

func extractToken(ctx fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return "", fiber.ErrUnauthorized
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return "", fiber.ErrUnauthorized
	}
	return tokenString, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return jwtSecret, nil
	})
}

func validateClaims(ctx fiber.Ctx, token *jwt.Token) error {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return fiber.ErrUnauthorized
			}
		}

		ctx.Locals("userID", claims["sub"])
		return ctx.Next()
	}

	return fiber.ErrUnauthorized
}

func GenerateToken(userID int) (string, error) {
	claims := createClaims(userID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func createClaims(userID int) jwt.MapClaims {
	return jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
}

func GetUID(ctx fiber.Ctx) (int, error) {
	userID := ctx.Locals("userID")

	switch v := userID.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		log.Error(fmt.Sprintf("Couldn't get user id for %v", userID))
		return 0, fiber.ErrInternalServerError
	}
}
