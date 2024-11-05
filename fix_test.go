package homeassistant

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"
)

func AInit() {
	dsn := "user=postgres password=xxxxxx dbname=homeassistant_db host=192.168.1.7 port=15432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
}

func TestFixDao(t *testing.T) {
	AInit()
	meta := &StateMeta{}
	entityId := "sensor.nian_zong_hao_dian_liang"
	err := DB.Debug().Preload("States", func(db *gorm.DB) *gorm.DB {
		return db.Where("last_updated_ts > ?", 1720065712.0000000).Limit(1000).Order("states.last_updated_ts asc")
	}).First(meta, &StateMeta{EntityID: entityId}).Error
	assert.Equal(t, nil, err)

	var state State
	err = DB.First(&state).Error
	assert.Equal(t, nil, err)
}
