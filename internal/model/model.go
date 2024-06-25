package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	Birthday     time.Time
	PasswordHash []byte
}

type Subscription struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	SubscribedAtUserID uuid.UUID
	NotificationTime   time.Time
}
