// @Author: abbeymart | Abi Akindele | @Created: 2021-06-24 | @Updated: 2021-06-24
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package dbcrud

import (
	"gorm.io/gorm"
	"time"
)

type BaseModelType struct {
	ID        string         `json:"id" gorm:"primaryKey;default:uuid_generate_v4()" mcorm:"id"`
	Language  string         `json:"language" gorm:"not null;default:en-US" mcorm:"language"`
	Desc      string         `json:"desc" mcorm:"desc"`
	AppId     string         `json:"appId" mcorm:"app_id"`                           // application-id in a multi-hosted apps environment (e.g. cloud-env)
	IsActive  bool           `json:"isActive" gorm:"default:true" mcorm:"is_active"` // => activate by modelOptionsType settings...
	CreatedBy string         `json:"createdBy" mcorm:"created_by"`
	CreatedAt time.Time      `json:"createdAt" mcorm:"created_at"`
	UpdatedBy string         `json:"updatedBy" mcorm:"updated_by"`
	UpdatedAt time.Time      `json:"updatedAt" mcorm:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index" mcorm:"deleted_at"`
}

type AuditStampType struct {
	IsActive  bool      `json:"isActive"` // => activate by modelOptionsType settings...
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedBy string    `json:"updatedBy"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AppParamsType struct {
	AppId     string `json:"appId" mcorm:"app_id"`
	AccessKey string `json:"accessKey" mcorm:"access_key"`
	AppName   string `json:"appName" mcorm:"app_name"`
}

type AppType struct {
	BaseModelType
	AppName   string `json:"appName"`
	AccessKey string `json:"accessKey"`
	OwnerId   string `json:"ownerId"`
}

type UserInfoType struct {
	UserId    string `json:"userId" form:"userId" mcorm:"user_id"`
	Firstname string `json:"firstname" mcorm:"firstname"`
	Lastname  string `json:"lastname" mcorm:"lastname"`
	Language  string `json:"language" mcorm:"language"`
	LoginName string `json:"loginName" form:"loginName" mcorm:"login_name"`
	Token     string `json:"token" mcorm:"token"`
	Expire    int64  `json:"expire" mcorm:"expire"`
	Email     string `json:"email" form:"email" mcorm:"email"`
	Role      string `json:"role" mcorm:"role"`
}

type ValueType struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}
