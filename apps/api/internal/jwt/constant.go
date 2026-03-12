package jwt

import "time"

const TokenIssuer = "auth"

const RefreshTokenCookie = "refresh"
const RefreshTokenLifeTime = time.Hour * 24
const AccessTokenLifeTime = time.Minute * 15
