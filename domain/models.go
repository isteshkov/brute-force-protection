package domain

import "time"

type GeneralTechFields struct {
	Uid       string
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
