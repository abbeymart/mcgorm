// @Author: abbeymart | Abi Akindele | @Created: 2020-12-04 | @Updated: 2020-12-04
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect Audit Log

package mcgorm

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mcresponse"
	"gorm.io/gorm"
	"strings"
	"time"
)

// LogParam interfaces / types
type LogParam struct {
	AuditDb    *gorm.DB
	AuditTable string
}

// AuditLogOptionsType set the audit-log optional parameters
type AuditLogOptionsType struct {
	AuditTable    string
	TableName     string
	LogRecords    interface{}
	NewLogRecords interface{}
	//QueryParams   interface{}
}

// Audit describe the data-model for Audit log
type Audit struct {
	ID            string      `json:"id" gorm:"primaryKey;default:uuid_generate_v4()" mcorm:"id"`
	TableName     string      `json:"tableName" mcorm:"table_name"`
	LogRecords    interface{} `json:"logRecords" mcorm:"log_records"`
	NewLogRecords interface{} `json:"newLogRecords" mcorm:"new_log_records"`
	LogType       string      `json:"logType" mcorm:"log_type"`
	LogBy         string      `json:"logBy" mcorm:"log_by"`
	LogAt         time.Time   `json:"logAt" mcorm:"log_at"`
}

type AuditLogger interface {
	AuditLog(logType, userId string, options AuditLogOptionsType) (mcresponse.ResponseMessage, error)
}
type CreateLogger interface {
	CreateLog(tableName string, logRecords interface{}, userId string) (mcresponse.ResponseMessage, error)
}
type UpdateLogger interface {
	UpdateLog(tableName string, logRecords interface{}, newLogRecords interface{}, userId string) (mcresponse.ResponseMessage, error)
}
type ReadLogger interface {
	ReadLog(tableName string, logRecords interface{}, userId string) (mcresponse.ResponseMessage, error)
}
type DeleteLogger interface {
	DeleteLog(tableName string, logRecords interface{}, userId string) (mcresponse.ResponseMessage, error)
}
type AccessLogger interface {
	LoginLog(logRecords interface{}, userId string, tableName string) (mcresponse.ResponseMessage, error)
	LogoutLog(logRecords interface{}, userId string, tableName string) (mcresponse.ResponseMessage, error)
}

//type AuditCrudLogger interface {
//	CreateLog(tableName string, LogRecords interface{}, userId string) (mcresponse.ResponseMessage, error)
//	UpdateLog(TableName string, LogRecords interface{}, NewLogRecords interface{}, userId string) (mcresponse.ResponseMessage, error)
//	ReadLog(TableName string, LogRecords interface{}, userId string) (mcresponse.ResponseMessage, error)
//	DeleteLog(TableName string, LogRecords interface{}, userId string) (mcresponse.ResponseMessage, error)
//	LoginLog(LogRecords interface{}, userId string, TableName string) (mcresponse.ResponseMessage, error)
//	LogoutLog(LogRecords interface{}, userId string, TableName string) (mcresponse.ResponseMessage, error)
//	AuditLog(LogType, userId string, options AuditLogOptionsType) (mcresponse.ResponseMessage, error)
//}

// constants
// LogTypes
const (
	CreateLog = "create"
	UpdateLog = "update"
	ReadLog   = "read"
	GetLog    = "get"
	DeleteLog = "delete"
	RemoveLog = "remove"
	LoginLog  = "login"
	LogoutLog = "logout"
)

// NewAuditLog is the constructor function for the AuditLog
func NewAuditLog(auditDb *gorm.DB, auditTable string) LogParam {
	result := LogParam{}
	result.AuditDb = auditDb
	result.AuditTable = auditTable
	// default value
	if result.AuditTable == "" {
		result.AuditTable = "audits"
	}
	return result
}

// String() function implementation
func (log LogParam) String() string {
	return fmt.Sprintf(`
	AuditLog DB: %v \n AudiLog Table Name: %v \n
	`,
		log.AuditDb,
		log.AuditTable)
}

// AuditLog method compose and insert new audit-log record
func (log LogParam) AuditLog(logType, userId string, options AuditLogOptionsType) (mcresponse.ResponseMessage, error) {
	// variables
	logType = strings.ToLower(logType)
	logBy := userId
	var (
		tableName     = ""
		logRecords    interface{}
		newLogRecords interface{}
		//result        *gorm.DB
	)

	audit := Audit{}

	// log-cases
	switch logType {
	case CreateLog:
		// set params
		tableName = options.TableName
		logRecords = options.LogRecords
		// validate params
		var errorMessage = ""
		if tableName == "" {
			errorMessage = "Table or Collection name is required."
		}
		if logBy == "" {
			if errorMessage != "" {
				errorMessage = errorMessage + " | userId is required."
			} else {
				errorMessage = "userId is required."
			}
		}
		if logRecords == nil {
			if errorMessage != "" {
				errorMessage = errorMessage + " | Created record(s) information is required."
			} else {
				errorMessage = "Created record(s) information is required."
			}
		}
		if errorMessage != "" {
			return mcresponse.GetResMessage("paramsError",
				mcresponse.ResponseMessageOptions{
					Message: errorMessage,
					Value:   nil,
				}), errors.New(errorMessage)
		}
		// perform crud action
		audit = Audit{
			TableName:  tableName,
			LogRecords: logRecords,
			LogType:    logType,
			LogBy:      logBy,
			LogAt:      time.Now(),
		}
	case UpdateLog:
		// set params
		tableName = options.TableName
		logRecords = options.LogRecords
		newLogRecords = options.NewLogRecords
		// validate params
		var errorMessage = ""
		if tableName == "" {
			errorMessage = "Table or Collection name is required."
		}
		if logBy == "" {
			if errorMessage != "" {
				errorMessage = errorMessage + " | userId is required."
			} else {
				errorMessage = "userId is required."
			}
		}
		if logRecords == nil {
			if errorMessage != "" {
				errorMessage = errorMessage + " | Updated record(s) information is required."
			} else {
				errorMessage = "Updated record(s) information is required."
			}
		}
		if newLogRecords == nil {
			if errorMessage != "" {
				errorMessage = errorMessage + " | New/Update record(s) information is required."
			} else {
				errorMessage = "New/Update record(s) information is required."
			}
		}
		if errorMessage != "" {
			return mcresponse.GetResMessage("paramsError",
				mcresponse.ResponseMessageOptions{
					Message: errorMessage,
					Value:   nil,
				}), errors.New(errorMessage)
		}
		// perform crud action
		audit = Audit{
			TableName:     tableName,
			LogRecords:    logRecords,
			NewLogRecords: newLogRecords,
			LogType:       logType,
			LogBy:         logBy,
			LogAt:         time.Now(),
		}
	case GetLog, ReadLog:
		// set params
		tableName = options.TableName
		logRecords = options.LogRecords
		// validate params
		var errorMessage = ""
		if tableName == "" {
			errorMessage = "Table or Collection name is required."
		}
		if logBy == "" {
			if errorMessage != "" {
				errorMessage = errorMessage + " | userId is required."
			} else {
				errorMessage = "userId is required."
			}
		}
		if logRecords == nil {
			if errorMessage != "" {
				errorMessage = errorMessage + " | Read/Get Params/Keywords information is required."
			} else {
				errorMessage = "Read/Get Params/Keywords information is required."
			}
		}
		if errorMessage != "" {
			return mcresponse.GetResMessage("paramsError",
				mcresponse.ResponseMessageOptions{
					Message: errorMessage,
					Value:   nil,
				}), errors.New(errorMessage)
		}
		// perform crud action
		audit = Audit{
			TableName:  tableName,
			LogRecords: logRecords,
			LogType:    logType,
			LogBy:      logBy,
			LogAt:      time.Now(),
		}
	case DeleteLog, RemoveLog:
		// set params
		tableName = options.TableName
		logRecords = options.LogRecords
		// validate params
		var errorMessage = ""
		if tableName == "" {
			errorMessage = "Table or Collection name is required."
		}
		if logBy == "" {
			if errorMessage != "" {
				errorMessage = errorMessage + " | userId is required."
			} else {
				errorMessage = "userId is required."
			}
		}
		if logRecords == nil {
			if errorMessage != "" {
				errorMessage = errorMessage + " | Deleted record(s) information is required."
			} else {
				errorMessage = "Deleted record(s) information is required."
			}
		}
		if errorMessage != "" {
			return mcresponse.GetResMessage("paramsError",
				mcresponse.ResponseMessageOptions{
					Message: errorMessage,
					Value:   nil,
				}), errors.New(errorMessage)
		}
		// perform crud action
		audit = Audit{
			TableName:  tableName,
			LogRecords: logRecords,
			LogType:    logType,
			LogBy:      logBy,
			LogAt:      time.Now(),
		}
	case LoginLog:
		// set params
		tableName = options.TableName
		logRecords = options.LogRecords
		// validate params
		var errorMessage = ""
		if tableName == "" {
			errorMessage = "Table or Collection name is required."
		}
		if logBy == "" {
			if errorMessage != "" {
				errorMessage = errorMessage + " | userId is required."
			} else {
				errorMessage = "userId is required."
			}
		}
		if logRecords == nil {
			if errorMessage != "" {
				errorMessage = errorMessage + " | Login record(s) information is required."
			} else {
				errorMessage = "Login record(s) information is required."
			}
		}
		if errorMessage != "" {
			return mcresponse.GetResMessage("paramsError",
				mcresponse.ResponseMessageOptions{
					Message: errorMessage,
					Value:   nil,
				}), errors.New(errorMessage)
		}
		// perform crud action
		audit = Audit{
			TableName:  tableName,
			LogRecords: logRecords,
			LogType:    logType,
			LogBy:      logBy,
			LogAt:      time.Now(),
		}
	case LogoutLog:
		// set params
		tableName = options.TableName
		logRecords = options.LogRecords
		// validate params
		var errorMessage = ""
		if tableName == "" {
			errorMessage = "Table or Collection name is required."
		}
		if logBy == "" {
			if errorMessage != "" {
				errorMessage = errorMessage + " | userId is required."
			} else {
				errorMessage = "userId is required."
			}
		}
		if logRecords == nil {
			if errorMessage != "" {
				errorMessage = errorMessage + " | Logout record(s) information is required."
			} else {
				errorMessage = "Logout record(s) information is required."
			}
		}
		if errorMessage != "" {
			return mcresponse.GetResMessage("paramsError",
				mcresponse.ResponseMessageOptions{
					Message: errorMessage,
					Value:   nil,
				}), errors.New(errorMessage)
		}
		// perform crud action
		audit = Audit{
			TableName:  tableName,
			LogRecords: logRecords,
			LogType:    logType,
			LogBy:      logBy,
			LogAt:      time.Now(),
		}
	default:
		return mcresponse.GetResMessage("logError",
			mcresponse.ResponseMessageOptions{
				Message: "Unknown log type and/or incomplete log information",
				Value:   nil,
			}), errors.New("unknown log type and/or incomplete log information")
	}

	// perform audit-log-create task
	result := log.AuditDb.Create(&audit)

	// Handle error
	if result.Error != nil {
		errMsg := fmt.Sprintf("Log-record create-error: %v", result.Error.Error())
		return mcresponse.GetResMessage("logError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			}), errors.New(errMsg)
	}

	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "successful audit-log action",
			Value:   result.RowsAffected,
		}), nil
}

func (log LogParam) CreateLog(tableName string, logRecords interface{}, userId string) (mcresponse.ResponseMessage, error) {
	// validate params
	var errorMessage = ""
	if tableName == "" {
		errorMessage = "Table or Collection name is required."
	}
	if userId == "" {
		if errorMessage != "" {
			errorMessage = errorMessage + " | userId is required."
		} else {
			errorMessage = "userId is required."
		}
	}
	if logRecords == nil {
		if errorMessage != "" {
			errorMessage = errorMessage + " | Created record(s) information is required."
		} else {
			errorMessage = "Created record(s) information is required."
		}
	}
	if errorMessage != "" {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: errorMessage,
				Value:   nil,
			}), errors.New(errorMessage)
	}
	// perform crud action
	audit := Audit{
		TableName:  tableName,
		LogRecords: logRecords,
		LogType:    CreateLog,
		LogBy:      userId,
		LogAt:      time.Now(),
	}
	// perform audit-log-insert task
	result := log.AuditDb.Create(&audit)

	// Handle error
	if result.Error != nil {
		errMsg := fmt.Sprintf("Log-record create-error: %v", result.Error.Error())
		return mcresponse.GetResMessage("logError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			}), errors.New(errMsg)
	}

	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "successful audit-log action",
			Value:   result.RowsAffected,
		}), nil
}

func (log LogParam) UpdateLog(tableName string, logRecords interface{}, newLogRecords interface{}, userId string) (mcresponse.ResponseMessage, error) {
	// validate params
	var errorMessage = ""
	if tableName == "" {
		errorMessage = "Table or Collection name is required."
	}
	if userId == "" {
		if errorMessage != "" {
			errorMessage = errorMessage + " | userId is required."
		} else {
			errorMessage = "userId is required."
		}
	}
	if logRecords == nil {
		if errorMessage != "" {
			errorMessage = errorMessage + " | Updated record(s) information is required."
		} else {
			errorMessage = "Updated record(s) information is required."
		}
	}
	if newLogRecords == nil {
		if errorMessage != "" {
			errorMessage = errorMessage + " | New/Update record(s) information is required."
		} else {
			errorMessage = "New/Update record(s) information is required."
		}
	}
	if errorMessage != "" {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: errorMessage,
				Value:   nil,
			}), errors.New(errorMessage)
	}

	audit := Audit{
		TableName:     tableName,
		LogRecords:    logRecords,
		NewLogRecords: newLogRecords,
		LogType:       CreateLog,
		LogBy:         userId,
		LogAt:         time.Now(),
	}
	// perform audit-log-insert task
	result := log.AuditDb.Create(&audit)

	// Handle error
	if result.Error != nil {
		errMsg := fmt.Sprintf("Log-record create-error: %v", result.Error.Error())
		return mcresponse.GetResMessage("logError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			}), errors.New(errMsg)
	}

	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "successful audit-log action",
			Value:   result.RowsAffected,
		}), nil
}

func (log LogParam) ReadLog(tableName string, logRecords interface{}, userId string) (mcresponse.ResponseMessage, error) {
	// validate params
	var errorMessage = ""
	if tableName == "" {
		errorMessage = "Table or Collection name is required."
	}
	if userId == "" {
		if errorMessage != "" {
			errorMessage = errorMessage + " | userId is required."
		} else {
			errorMessage = "userId is required."
		}
	}
	if logRecords == nil {
		if errorMessage != "" {
			errorMessage = errorMessage + " | Read/Get Params/Keywords information is required."
		} else {
			errorMessage = "Read/Get Params/Keywords information is required."
		}
	}
	if errorMessage != "" {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: errorMessage,
				Value:   nil,
			}), errors.New(errorMessage)
	}
	// perform crud action
	audit := Audit{
		TableName:  tableName,
		LogRecords: logRecords,
		LogType:    CreateLog,
		LogBy:      userId,
		LogAt:      time.Now(),
	}
	// perform audit-log-insert task
	result := log.AuditDb.Create(&audit)

	// Handle error
	if result.Error != nil {
		errMsg := fmt.Sprintf("Log-record create-error: %v", result.Error.Error())
		return mcresponse.GetResMessage("logError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			}), errors.New(errMsg)
	}

	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "successful audit-log action",
			Value:   result.RowsAffected,
		}), nil
}

func (log LogParam) DeleteLog(tableName string, logRecords interface{}, userId string) (mcresponse.ResponseMessage, error) {
	// validate params
	var errorMessage = ""
	if tableName == "" {
		errorMessage = "Table or Collection name is required."
	}
	if userId == "" {
		if errorMessage != "" {
			errorMessage = errorMessage + " | userId is required."
		} else {
			errorMessage = "userId is required."
		}
	}
	if logRecords == nil {
		if errorMessage != "" {
			errorMessage = errorMessage + " | Deleted record(s) information is required."
		} else {
			errorMessage = "Deleted record(s) information is required."
		}
	}
	if errorMessage != "" {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: errorMessage,
				Value:   nil,
			}), errors.New(errorMessage)
	}
	// perform crud action
	audit := Audit{
		TableName:  tableName,
		LogRecords: logRecords,
		LogType:    CreateLog,
		LogBy:      userId,
		LogAt:      time.Now(),
	}
	// perform audit-log-insert task
	result := log.AuditDb.Create(&audit)

	// Handle error
	if result.Error != nil {
		errMsg := fmt.Sprintf("Log-record create-error: %v", result.Error.Error())
		return mcresponse.GetResMessage("logError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			}), errors.New(errMsg)
	}

	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "successful audit-log action",
			Value:   result.RowsAffected,
		}), nil
}

func (log LogParam) LoginLog(logRecords interface{}, userId string, tableName string) (mcresponse.ResponseMessage, error) {
	// validate params
	var errorMessage = ""
	if tableName == "" {
		errorMessage = "Table or Collection name is required."
	}
	if userId == "" {
		if errorMessage != "" {
			errorMessage = errorMessage + " | userId is required."
		} else {
			errorMessage = "userId is required."
		}
	}
	if logRecords == nil {
		if errorMessage != "" {
			errorMessage = errorMessage + " | Login record(s) information is required."
		} else {
			errorMessage = "Login record(s) information is required."
		}
	}
	if errorMessage != "" {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: errorMessage,
				Value:   nil,
			}), errors.New(errorMessage)
	}
	// perform crud action
	audit := Audit{
		TableName:  tableName,
		LogRecords: logRecords,
		LogType:    CreateLog,
		LogBy:      userId,
		LogAt:      time.Now(),
	}
	// perform audit-log-insert task
	result := log.AuditDb.Create(&audit)

	// Handle error
	if result.Error != nil {
		errMsg := fmt.Sprintf("Log-record create-error: %v", result.Error.Error())
		return mcresponse.GetResMessage("logError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			}), errors.New(errMsg)
	}

	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "successful audit-log action",
			Value:   result.RowsAffected,
		}), nil
}

func (log LogParam) LogoutLog(logRecords interface{}, userId string, tableName string) (mcresponse.ResponseMessage, error) {
	// validate params
	var errorMessage = ""
	if tableName == "" {
		errorMessage = "Table or Collection name is required."
	}
	if userId == "" {
		if errorMessage != "" {
			errorMessage = errorMessage + " | userId is required."
		} else {
			errorMessage = "userId is required."
		}
	}
	if logRecords == nil {
		if errorMessage != "" {
			errorMessage = errorMessage + " | Logout record(s) information is required."
		} else {
			errorMessage = "Logout record(s) information is required."
		}
	}
	if errorMessage != "" {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: errorMessage,
				Value:   nil,
			}), errors.New(errorMessage)
	}
	// perform crud action
	audit := Audit{
		TableName:  tableName,
		LogRecords: logRecords,
		LogType:    CreateLog,
		LogBy:      userId,
		LogAt:      time.Now(),
	}
	// perform audit-log-insert task
	result := log.AuditDb.Create(&audit)

	// Handle error
	if result.Error != nil {
		errMsg := fmt.Sprintf("Log-record create-error: %v", result.Error.Error())
		return mcresponse.GetResMessage("logError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			}), errors.New(errMsg)
	}

	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "successful audit-log action",
			Value:   result.RowsAffected,
		}), nil
}
