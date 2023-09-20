package models

type JWTClaims struct {
	Email       string      `json:"email,omitempty"`
	AppMetadata AppMetadata `json:"app_metadata,omitempty"`
	Issuer      string      `json:"iss,omitempty"`
	Subject     string      `json:"sub,omitempty"`
	Audience    string      `json:"aud,omitempty"`
	Expiration  int64       `json:"exp,omitempty"`
	IssuedAt    int64       `json:"iat,omitempty"`
	JWTID       string      `json:"jti,omitempty"`
}

type AppMetadata struct {
	Authorization Authorization `json:"authorization,omitempty"`
}

type Authorization struct {
	Roles []string `json:"roles,omitempty"`
}
