// @Author: abbeymart | Abi Akindele | @Created: 2021-06-24 | @Updated: 2021-06-24
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect - getRecord

package mcgorm

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcresponse"
)

func (crud Crud) GetById(modelRef interface{}, id string) mcresponse.ResponseMessage {
	// perform get-query
	//var result *gorm.DB
	result := crud.GormDb.Limit(crud.Limit).Offset(crud.Skip).Where("id = ?", id).Find(&modelRef)
	if result.Error != nil {
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// rows
	var records []interface{}
	rows, err := result.Rows()
	if err != nil {
		errMsg := fmt.Sprintf("%v", result.Error.Error())
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			})
	}
	for rows.Next() {
		err = crud.GormDb.ScanRows(rows, &modelRef)
		if err != nil {
			return mcresponse.GetResMessage("readError",
				mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("%v", result.Error.Error()),
					Value:   nil,
				})
		}
		// get snapshot value from the pointer | transform value to json-value([]byte)
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming record(row-value) into json-value([]byte): %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming json-value to result-value: %v", jErr.Error()),
				Value:   nil,
			})
		}
		records = append(records, gValue)
	}
	var totalRecordsCount int64
	var _ = crud.GormDb.Find(modelRef).Count(&totalRecordsCount)
	// logRead
	var logRes mcresponse.ResponseMessage
	if crud.LogRead {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: []string{id},
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
			Value: GetResultType{
				Records: records,
				Stats: GetStatType{
					Skip:              crud.Skip,
					Limit:             crud.Limit,
					RecordsCount:      int(result.RowsAffected),
					TotalRecordsCount: int(totalRecordsCount),
				},
				LogRes: logRes,
			},
		})
}

func (crud Crud) GetByIds(modelRef interface{}) mcresponse.ResponseMessage {
	if len(crud.RecordIds) < 1 {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: "recordIds param is required to get-record-by-id",
				Value:   nil,
			})
	}
	// perform get-query
	result := crud.GormDb.Limit(crud.Limit).Offset(crud.Skip).Where("id in ?", crud.RecordIds).Find(&modelRef)
	if result.Error != nil {
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// rows
	var records []interface{}
	rows, err := result.Rows()
	if err != nil {
		errMsg := fmt.Sprintf("%v", result.Error.Error())
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			})
	}
	for rows.Next() {
		err = crud.GormDb.ScanRows(rows, &modelRef)
		if err != nil {
			return mcresponse.GetResMessage("readError",
				mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("%v", result.Error.Error()),
					Value:   nil,
				})
		}
		// get snapshot value from the pointer | transform value to json-value([]byte)
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming record(row-value) into json-value([]byte): %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming json-value to result-value: %v", jErr.Error()),
				Value:   nil,
			})
		}
		records = append(records, gValue)
	}
	var totalRecordsCount int64
	var _ = crud.GormDb.Find(modelRef).Count(&totalRecordsCount)
	// logRead
	var logRes mcresponse.ResponseMessage
	if crud.LogRead {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: crud.RecordIds,
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
			Value: GetResultType{
				Records: records,
				Stats: GetStatType{
					Skip:              crud.Skip,
					Limit:             crud.Limit,
					RecordsCount:      int(result.RowsAffected),
					TotalRecordsCount: int(totalRecordsCount),
				},
				LogRes: logRes,
			},
		})
}

func (crud Crud) GetByParam(modelRef interface{}) mcresponse.ResponseMessage {
	if crud.QueryParams == nil {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: "queryParams is required to get-record-by-param",
				Value:   nil,
			})
	}
	// perform get-query
	result := crud.GormDb.Limit(crud.Limit).Offset(crud.Skip).Where(crud.QueryParams).Find(&modelRef)
	if result.Error != nil {
		errMsg := fmt.Sprintf("%v", result.Error.Error())
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			})
	}
	// rows
	var records []interface{}
	rows, err := result.Rows()
	if err != nil {
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	for rows.Next() {
		err = crud.GormDb.ScanRows(rows, &modelRef)
		if err != nil {
			return mcresponse.GetResMessage("readError",
				mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("%v", result.Error.Error()),
					Value:   nil,
				})
		}
		// get snapshot value from the pointer | transform value to json-value([]byte)
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming record(row-value) into json-value([]byte): %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming json-value to result-value: %v", jErr.Error()),
				Value:   nil,
			})
		}
		records = append(records, gValue)
	}
	var totalRecordsCount int64
	var _ = crud.GormDb.Find(modelRef).Count(&totalRecordsCount)
	// logRead
	var logRes mcresponse.ResponseMessage
	if crud.LogRead {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: crud.QueryParams,
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
			Value: GetResultType{
				Records: records,
				Stats: GetStatType{
					Skip:              crud.Skip,
					Limit:             crud.Limit,
					RecordsCount:      int(result.RowsAffected),
					TotalRecordsCount: int(totalRecordsCount),
				},
				LogRes: logRes,
			},
		})
}

func (crud Crud) GetAll(modelRef interface{}) mcresponse.ResponseMessage {
	// perform get-query
	result := crud.GormDb.Limit(crud.Limit).Offset(crud.Skip).Find(&modelRef)
	if result.Error != nil {
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// rows
	var records []interface{}
	rows, err := result.Rows()
	if err != nil {
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	for rows.Next() {
		err = crud.GormDb.ScanRows(rows, &modelRef)
		if err != nil {
			return mcresponse.GetResMessage("readError",
				mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("%v", result.Error.Error()),
					Value:   nil,
				})
		}
		// get snapshot value from the pointer | transform value to json-value([]byte)
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming record(row-value) into json-value([]byte): %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming json-value to result-value: %v", jErr.Error()),
				Value:   nil,
			})
		}
		records = append(records, gValue)
	}
	var totalRecordsCount int64
	var _ = crud.GormDb.Find(modelRef).Count(&totalRecordsCount)
	// logRead
	var logRes mcresponse.ResponseMessage
	if crud.LogRead {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: map[string]interface{}{"getType": "All Records"},
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
			Value: GetResultType{
				Records: records,
				Stats: GetStatType{
					Skip:              crud.Skip,
					Limit:             crud.Limit,
					RecordsCount:      int(result.RowsAffected),
					TotalRecordsCount: int(totalRecordsCount),
				},
				LogRes: logRes,
			},
		})
}
