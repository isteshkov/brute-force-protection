package models

import "time"

type GeneralTechFields struct {
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Subnet struct {
	GeneralTechFields
	Address       string
	IsBlacklisted bool
}
