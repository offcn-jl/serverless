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
type AppVersionControlInfo struct {
	gorm.Model
	AppID    string `gorm:"not null"` // 应用 ID
	AppName  string `gorm:"not null"` // 应用名称
	Version  string `gorm:"not null"` // 版本
	Download string `gorm:"not null"` // 下载链接
}

// 应用程序 授权信息
type AppAuthInfo struct {
	gorm.Model
	Authorized string `gorm:"not null"`                                 // 被授权的应用 ( 名称 )
	Token      string `gorm:"not null" json:"Token" binding:"required"` // 授权令牌
	ExpiresIn  uint   `gorm:"not null"`                                 // 授权有效期 ( 秒 )
}
