package data

import (
	"github.com/google/uuid"
	"time"
)

// UserSession is a struct that is using for store info about users logged in a Ferrum authorization server
/* UserId - uuid representing unique user identifier
 * Started - time when token was Issued
 * Expired - time when session expires
 * RefreshExpired - time when refresh expires
 * JwtAccessToken and JwtRefreshToken - access and refresh tokens
 */
type UserSession struct {
	Id              uuid.UUID
	UserId          uuid.UUID
	Started         time.Time
	Expired         time.Time
	RefreshExpired  time.Time
	JwtAccessToken  string
	JwtRefreshToken string
}
