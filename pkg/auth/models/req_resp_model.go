package models

import (
)

type RegReq struct {
	Email string `json:"email" binding:"required,email"`
	Pass  string `json:"password" binding:"required,min=8"`
	Role  string `json:"role"`
}

type LoginReq struct {
	Email string `json:"email" binding:"required,email"`
	Pass  string `json:"password" binding:"required"`
}

type AuthResp struct {
	Token string `json:"access_token"`
	Email string `json:"email"`
	Role  string `json:"role"`
}