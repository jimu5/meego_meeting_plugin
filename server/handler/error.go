package handler

import "errors"

var (
	ErrInvalidParam    = errors.New("error invalid param")
	ErrInvalidUserInfo = errors.New("err invalid user info")

	ErrEmptyCalendarSearch = errors.New("empty calendar search")
)
