/*
   @Time : 2020/4/22 4:56 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : apps
   @Software: GoLand
*/

package structs

import "github.com/jinzhu/gorm"

/**
 * 版本控制表
 */
type VersionControlInfo struct {
	gorm.Model
	AppID    string `gorm:"not null"` // 应用 ID
	AppName  string `gorm:"not null"` // 应用名称
	Version  string `gorm:"not null"` // 版本
	Download string `gorm:"not null"` // 下载链接
}
