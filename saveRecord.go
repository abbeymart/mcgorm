// @Author: abbeymart | Abi Akindele | @Created: 2021-06-24 | @Updated: 2021-06-24
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect - saveRecord

package mcgorm

import (
	"fmt"
	"github.com/abbeymart/mcresponse"
	"gorm.io/gorm"
	"strings"
)

func (crud Crud) Create(rec interface{}) mcresponse.ResponseMessage {
	result := crud.GormDb.Create(&rec)
	if result.Error != nil {
		return mcresponse.GetResMessage("insertError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// LogCreate
	var logRes mcresponse.ResponseMessage
	var err error
	if crud.LogCreate {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Create, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: rec,
			TableName:  crud.TableName,
		})
		if err != nil {
			logRes = mcresponse.ResponseMessage{
				Code:    "logError",
				Message: fmt.Sprintf("Audit-log error: %v", err.Error()),
				Value:   nil,
			}
		}
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: CrudResultType{
				RecordCount: int(result.RowsAffected),
				LogRes:       logRes,
				TaskType:     crud.TaskType,
			},
		})
}

func (crud Crud) CreateBatch(recs interface{}, batch int) mcresponse.ResponseMessage {
	// default value
	if batch == 0 {
		batch = 10000
	}
	result := crud.GormDb.CreateInBatches(&recs, batch)
	if result.Error != nil {
		return mcresponse.GetResMessage("insertError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// LogCreate
	var logRes mcresponse.ResponseMessage
	var err error
	if crud.LogCreate {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Create, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: recs,
			TableName:  crud.TableName,
		})
		if err != nil {
			logRes = mcresponse.ResponseMessage{
				Code:    "logError",
				Message: fmt.Sprintf("Audit-log error: %v", err.Error()),
				Value:   nil,
			}
		}
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: CrudResultType{
				RecordCount: int(result.RowsAffected),
				LogRes:       logRes,
				TaskType:     crud.TaskType,
			},
		})
}

func (crud Crud) UpdateById(model interface{}, rec interface{}, id string) mcresponse.ResponseMessage {
	var getRes mcresponse.ResponseMessage
	if crud.LogUpdate {
		// get current records
		getRes = crud.GetById(model, id)
	}
	// convert struct to map to save all fields (including zero-value fields)
	mapRec, err := StructToCaseUnderscoreMap(rec)
	if err != nil {
		return mcresponse.GetResMessage("updateError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", err.Error()),
				Value:   nil,
			})
	}
	// destruct id and other-fields from update-record
	upRec := map[string]interface{}{}
	for k, v := range mapRec {
		if k == "id" {
			continue
		}
		upRec[k] = v
	}
	result := crud.GormDb.Model(&model).Where("id = ?", id).Updates(upRec)
	if result.Error != nil {
		return mcresponse.GetResMessage("updateError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// LogUpdate
	var logRes mcresponse.ResponseMessage
	if crud.LogUpdate {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Update, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords:    getRes.Value,
			NewLogRecords: map[string]interface{}{"id": []string{id}, "record": rec},
			TableName:     crud.TableName,
		})
		if err != nil {
			logRes = mcresponse.ResponseMessage{
				Code:    "logError",
				Message: fmt.Sprintf("Audit-log error: %v", err.Error()),
				Value:   nil,
			}
		}
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: CrudResultType{
				RecordCount: int(result.RowsAffected),
				LogRes:      logRes,
				TaskType:    crud.TaskType,
			},
		})
}

func (crud Crud) UpdateByIds(model interface{}, rec interface{}) mcresponse.ResponseMessage {
	if len(crud.RecordIds) < 1 {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: "records-Ids param is required to get-record-by-ids",
				Value:   nil,
			})
	}
	var getRes mcresponse.ResponseMessage
	if crud.LogUpdate {
		// get current records
		getRes = crud.GetByIds(model)
	}
	// convert struct to map to save all fields (including zero-value fields)
	mapRec, err := StructToCaseUnderscoreMap(rec)
	if err != nil {
		return mcresponse.GetResMessage("updateError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", err.Error()),
				Value:   nil,
			})
	}
	// destruct id and other-fields from update-record (mapRec)
	upRec := map[string]interface{}{}
	for k, v := range mapRec {
		if k == "id" {
			continue
		}
		upRec[k] = v
	}
	result := crud.GormDb.Model(&model).Where("id in ?", crud.RecordIds).Updates(upRec)
	if result.Error != nil {
		return mcresponse.GetResMessage("updateError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// LogUpdate
	var logRes mcresponse.ResponseMessage
	if crud.LogUpdate {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Update, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords:    getRes.Value,
			NewLogRecords: map[string]interface{}{"id": crud.RecordIds, "record": rec},
			TableName:     crud.TableName,
		})
		if err != nil {
			logRes = mcresponse.ResponseMessage{
				Code:    "logError",
				Message: fmt.Sprintf("Audit-log error: %v", err.Error()),
				Value:   nil,
			}
		}
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: CrudResultType{
				RecordCount: int(result.RowsAffected),
				LogRes:       logRes,
				TaskType:     crud.TaskType,
			},
		})
}

func (crud Crud) UpdateByParam(model interface{}, rec interface{}) mcresponse.ResponseMessage {
	if crud.QueryParams == nil {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: "queryParams is required to get-record-by-param",
				Value:   nil,
			})
	}
	var getRes mcresponse.ResponseMessage
	if crud.LogUpdate {
		// get current records
		getRes = crud.GetByParam(model)
	}
	// convert struct to map to save all fields (including zero-value fields)
	mapRec, err := StructToCaseUnderscoreMap(rec)
	if err != nil {
		return mcresponse.GetResMessage("updateError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", err.Error()),
				Value:   nil,
			})
	}
	// destruct id and other-fields from update-record (mapRec)
	upRec := map[string]interface{}{}
	for k, v := range mapRec {
		if k == "id" {
			continue
		}
		upRec[k] = v
	}
	result := crud.GormDb.Model(&model).Where(crud.QueryParams).Updates(upRec)
	if result.Error != nil {
		return mcresponse.GetResMessage("updateError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// LogUpdate
	var logRes mcresponse.ResponseMessage
	if crud.LogUpdate {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Update, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords:    getRes.Value,
			NewLogRecords: map[string]interface{}{"queryParams": crud.QueryParams, "record": rec},
			TableName:     crud.TableName,
		})
		if err != nil {
			logRes = mcresponse.ResponseMessage{
				Code:    "logError",
				Message: fmt.Sprintf("Audit-log error: %v", err.Error()),
				Value:   nil,
			}
		}
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: CrudResultType{
				RecordCount: int(result.RowsAffected),
				LogRes:       logRes,
				TaskType:     crud.TaskType,
			},
		})
}

func (crud Crud) Update(model interface{}, recs interface{}) mcresponse.ResponseMessage {
	// TODO: validate recs as slice of interface/records(struct/map)

	var getRes mcresponse.ResponseMessage
	if crud.LogUpdate {
		// get current records
		// destruct ids from update-record
		var recIds []string
		for _, mapVal := range crud.ActionParams {
			for k, v := range mapVal {
				if strings.ToLower(k) == "id" {
					recIds = append(recIds, v.(string))
				}
				continue
			}
		}
		crud.RecordIds = recIds
		getRes = crud.GetByIds(model)
	}

	// TODO: perform batch updates | transactional
	var result *gorm.DB
	resultCount := 0
	for _, record := range recs.([]interface{}) {
		// convert struct to map to save all fields (including zero-value fields)
		mapRec, err := StructToCaseUnderscoreMap(record)
		if err != nil {
			return mcresponse.GetResMessage("updateError",
				mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("%v", err.Error()),
					Value:   nil,
				})
		}
		// TODO: destruct id and other-fields from update-record (mapRec)
		var id string
		upRec := map[string]interface{}{}
		for k, v := range mapRec {
			if k == "id" {
				id = k
				continue
			}
			upRec[k] = v
		}
		result = crud.GormDb.Model(&model).Where("id = ?", id).Updates(upRec)
		if result.Error != nil {
			return mcresponse.GetResMessage("updateError",
				mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("%v", result.Error.Error()),
					Value:   nil,
				})
		}
		resultCount++
	}
	// LogUpdate
	var logRes mcresponse.ResponseMessage
	var err error
	if crud.LogUpdate {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Update, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords:    getRes.Value,
			NewLogRecords: recs,
			TableName:     crud.TableName,
		})
		if err != nil {
			logRes = mcresponse.ResponseMessage{
				Code:    "logError",
				Message: fmt.Sprintf("Audit-log error: %v", err.Error()),
				Value:   nil,
			}
		}
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: CrudResultType{
				RecordCount: resultCount,
				LogRes:       logRes,
				TaskType:     crud.TaskType,
			},
		})
}
