package storage

import (
	"avatarDetail/connect"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)
import "time"
import "fmt"

type AvatarDetail struct {
	Avatar     string `gorm:"column:avatar"`
	AvatarTime uint   `gorm:"column:avatar_time"`
	Height     int    `gorm:"column:height"`
	Id         uint   `gorm:"column:id;primary_key"`
	Uid        uint   `gorm:"column:uid"`
	Width      int    `gorm:"column:width"`
}

var dbConf *viper.Viper

func init() {
	dbConf, _ = connect.ConnectConfig("database")
}

func (a *AvatarDetail) BeforeSave(scope *gorm.Scope) (err error) {
	scope.SetColumn("update_time", time.Now().Unix())
	return nil
}
func (a *AvatarDetail) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("add_time", time.Now().Unix())
	return nil
}

// TableName sets the insert table name for this struct type
func (a *AvatarDetail) TableName(shardKey uint32) string {

	tn, stc := dbConf.GetString("AvatarDetail.table_name"), dbConf.GetUint32("AvatarDetail.spli_table_count")

	return fmt.Sprintf("%s%d", tn, shardKey%stc)
}

// // TableName sets the insert table name for this struct type
// func (a *AvatarDetail) TableName(shardKey uint32) string {
// 	return "kk_avatar_detail_" + fmt.Sprintf("%d", shardKey%256)
// }
func (a *AvatarDetail) GetShardKey() string {
	return dbConf.GetString("AvatarDetail.shard_key")
}
