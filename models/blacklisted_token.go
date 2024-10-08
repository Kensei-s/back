// models/blacklisted_token.go
package models

import "gorm.io/gorm"

type BlacklistedToken struct {
	gorm.Model
	Token string `json:"token"`
}
