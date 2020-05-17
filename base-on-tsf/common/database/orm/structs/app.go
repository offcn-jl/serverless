/*
   @Time : 2020/5/17 1:53 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : app
   @Software: GoLand
*/

package structs

import "github.com/jinzhu/gorm"

// 应用程序 版本控制表
type VersionControlInfo struct {
	gorm.Model
	AppID    string `gorm:"not null"` // 应用 ID
	AppName  string `gorm:"not null"` // 应用名称
	Version  string `gorm:"not null"` // 版本
	Download string `gorm:"not null"` // 下载链接
}
