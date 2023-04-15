package entity

import (
	"time"
)

type Port struct {
	Id           int64      `db:"id"`
	PrimaryUnloc string     `db:"primary_unloc"`
	Name         string     `db:"name"`
	Code         string     `db:"code"`
	City         string     `db:"city"`
	Country      string     `db:"country"`
	Alias        *[]Alias   `db:"alias"`
	Regions      *[]Region  `db:"regions"`
	Coordinantes [2]float32 `db:"coordinantes"`
	Province     string     `db:"province"`
	Timezone     string     `db:"timezone"`
	Unlocs       *[]Unloc   `db:"unlocs"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	// // time.Time cannot be set to nil (for records that have not been deleted)
	// // so setting DeletedAt to *time.Time
	DeletedAt *time.Time `db:"deleted_at"`
}

type Alias struct {
	PortId int64  `db:"port_id"`
	Name   string `db:"name"`
}

type Region struct {
	PortId int64  `db:"port_id"`
	Name   string `db:"name"`
}

type Unloc struct {
	PortId int64  `db:"port_id"`
	Name   string `db:"name"`
}
