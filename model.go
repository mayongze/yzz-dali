package homeassistant

import (
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

type StateAttribute struct {
	AttributesID int    `gorm:"primaryKey;autoIncrement"`
	Hash         int64  `gorm:"type:bigint"`
	SharedAttrs  string `gorm:"type:text"`

	// Relationships
	States []State `gorm:"foreignKey:AttributesID"`
}

func (StateAttribute) TableName() string {
	return "state_attributes"
}

type StateMeta struct {
	MetadataID int    `gorm:"primaryKey;autoIncrement"`
	EntityID   string `gorm:"type:varchar(255)"`

	// Relationships
	States []State1 `gorm:"foreignKey:MetadataID"`
}

func (StateMeta) TableName() string {
	return "states_meta"
}

type State1 struct {
	StateID            int       `gorm:"primaryKey;autoIncrement"`
	EntityID           string    `gorm:"type:char(1)"`
	State              string    `gorm:"type:varchar(255)"`
	Attributes         string    `gorm:"type:char(1)"`
	EventID            int16     `gorm:"type:smallint"`
	LastChanged        time.Time `gorm:"type:timestamp with time zone"`
	LastChangedTS      float64   `gorm:"type:double precision"`
	LastUpdated        time.Time `gorm:"type:timestamp with time zone"`
	LastUpdatedTS      float64   `gorm:"type:double precision"`
	OldStateID         int       `gorm:"type:int"`
	AttributesID       int       `gorm:"type:int"`
	ContextID          string    `gorm:"type:char(1)"`
	ContextUserID      string    `gorm:"type:char(1)"`
	ContextParentID    string    `gorm:"type:char(1)"`
	OriginIdx          int16     `gorm:"type:smallint"`
	ContextIDBin       []byte    `gorm:"type:bytea"`
	ContextUserIDBin   []byte    `gorm:"type:bytea"`
	ContextParentIDBin []byte    `gorm:"type:bytea"`
	MetadataID         int       `gorm:"type:int"`
	LastReportedTS     float64   `gorm:"type:double precision"`

	// Relationships
	StateAttribute `gorm:"foreignKey:AttributesID"`
	StateMeta      `gorm:"foreignKey:MetadataID"`
}

func (State1) TableName() string {
	return "states"
}
