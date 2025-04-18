// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package auth

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/skychains/chain/api"
)

// Service that serves the Auth API functionality.
type Service struct {
	auth *auth
}

type Password struct {
	Password string `json:"password"` // The authorization password
}

type NewTokenArgs struct {
	Password
	// Endpoints that may be accessed with this token e.g. if endpoints is
	// ["/ext/bc/X", "/ext/admin"] then the token holder can hit the X-Chain API
	// and the admin API. If [Endpoints] contains an element "*" then the token
	// allows access to all API endpoints. [Endpoints] must have between 1 and
	// [maxEndpoints] elements
	Endpoints []string `json:"endpoints"`
}

type Token struct {
	Token string `json:"token"` // The new token. Expires in [TokenLifespan].
}

func (s *Service) NewToken(_ *http.Request, args *NewTokenArgs, reply *Token) error {
	s.auth.log.Debug("API called",
		zap.String("service", "auth"),
		zap.String("method", "newToken"),
	)

	var err error
	reply.Token, err = s.auth.NewToken(args.Password.Password, defaultTokenLifespan, args.Endpoints)
	return err
}

type RevokeTokenArgs struct {
	Password
	Token
}

func (s *Service) RevokeToken(_ *http.Request, args *RevokeTokenArgs, _ *api.EmptyReply) error {
	s.auth.log.Debug("API called",
		zap.String("service", "auth"),
		zap.String("method", "revokeToken"),
	)

	return s.auth.RevokeToken(args.Token.Token, args.Password.Password)
}

type ChangePasswordArgs struct {
	OldPassword string `json:"oldPassword"` // Current authorization password
	NewPassword string `json:"newPassword"` // New authorization password
}

func (s *Service) ChangePassword(_ *http.Request, args *ChangePasswordArgs, _ *api.EmptyReply) error {
	s.auth.log.Debug("API called",
		zap.String("service", "auth"),
		zap.String("method", "changePassword"),
	)

	return s.auth.ChangePassword(args.OldPassword, args.NewPassword)
}
