package service

import "errors"

var ErrInvalidToken = errors.New("invalid token")
var ErrCreateToken = errors.New("failed to create tokens")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidAmount = errors.New("invalid amount")
