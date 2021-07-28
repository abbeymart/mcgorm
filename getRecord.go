// @Author: abbeymart | Abi Akindele | @Created: 2021-06-24 | @Updated: 2021-06-24
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect - getRecord

package dbcrud

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcresponse"
	"gorm.io/gorm"
)

type GetStats struct {
	Skip              int
	Limit             int
	RecordsCount      int
	TotalRecordsCount int
}

type GetResult struct {
	Records interface{}
	Stats   GetStats
}

func (crud Crud) GetById(modelRef interface{}, id string) mcresponse.ResponseMessage {
	// for limit > 0 and skip/offset > 0 OR limit > 0
	var result *gorm.DB
	if crud.Limit > 0 && crud.Skip > 0 {
		result = crud.GormDb.Limit(crud.Limit).Offset(crud.Skip).Where("id = ?", id).Find(&modelRef)
	} else if crud.Limit > 0 {
		result = crud.GormDb.Limit(crud.Limit).Where("id = ?", id).Find(&modelRef)
	} else {
		result = crud.GormDb.Where("id = ?", id).Find(&modelRef)
	}

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
		// get snapshot value from the pointer | transform value to json-value-format
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		records = append(records, gValue)
	}
	var totalRecordsCount int64
	var _ = crud.GormDb.Find(modelRef).Count(&totalRecordsCount)
	// logRead
	if crud.LogRead {
		_, _ = crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: []string{id},
			TableName:  crud.TableName,
		})
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: GetResult{
				Records: records,
				Stats: GetStats{
					Skip:              crud.Skip,
					Limit:             crud.Limit,
					RecordsCount:      int(result.RowsAffected),
					TotalRecordsCount: int(totalRecordsCount),
				},
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
	// for limit > 0 and skip/offset > 0 OR limit > 0
	var result *gorm.DB
	if crud.Limit > 0 && crud.Skip > 0 {
		result = crud.GormDb.Limit(crud.Limit).Offset(crud.Skip).Where("id in ?", crud.RecordIds).Find(&modelRef)
	} else if crud.Limit > 0 {
		result = crud.GormDb.Limit(crud.Limit).Where("id in ?", crud.RecordIds).Find(&modelRef)
	} else {
		result = crud.GormDb.Where("id in ?", crud.RecordIds).Find(&modelRef)
	}

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
		// get snapshot value from the pointer | transform value to json-value-format
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		records = append(records, gValue)
	}
	var totalRecordsCount int64
	var _ = crud.GormDb.Find(modelRef).Count(&totalRecordsCount)
	// logRead
	if crud.LogRead {
		_, _ = crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: crud.RecordIds,
			TableName:  crud.TableName,
		})
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: GetResult{
				Records: records,
				Stats: GetStats{
					Skip:              crud.Skip,
					Limit:             crud.Limit,
					RecordsCount:      int(result.RowsAffected),
					TotalRecordsCount: int(totalRecordsCount),
				},
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
	// for limit > 0 and skip/offset > 0 OR limit > 0
	var result *gorm.DB
	if crud.Limit > 0 && crud.Skip > 0 {
		result = crud.GormDb.Limit(crud.Limit).Offset(crud.Skip).Where(crud.QueryParams).Find(&modelRef)
	} else if crud.Limit > 0 {
		result = crud.GormDb.Limit(crud.Limit).Where(crud.QueryParams).Find(&modelRef)
	} else {
		result = crud.GormDb.Where(crud.QueryParams).Find(&modelRef)
	}

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
		// get snapshot value from the pointer | transform value to json-value-format
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		records = append(records, gValue)
	}
	var totalRecordsCount int64
	var _ = crud.GormDb.Find(modelRef).Count(&totalRecordsCount)
	// logRead
	if crud.LogRead {
		_, _ = crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: crud.QueryParams,
			TableName:  crud.TableName,
		})
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: GetResult{
				Records: records,
				Stats: GetStats{
					Skip:              crud.Skip,
					Limit:             crud.Limit,
					RecordsCount:      int(result.RowsAffected),
					TotalRecordsCount: int(totalRecordsCount),
				},
			},
		})
}

func (crud Crud) GetAll(modelRef interface{}) mcresponse.ResponseMessage {
	// for limit > 0 and skip/offset > 0 OR limit > 0
	var result *gorm.DB
	if crud.Limit > 0 && crud.Skip > 0 {
		result = crud.GormDb.Limit(crud.Limit).Offset(crud.Skip).Find(&modelRef)
	} else if crud.Limit > 0 {
		result = crud.GormDb.Limit(crud.Limit).Find(&modelRef)
	} else {
		result = crud.GormDb.Find(&modelRef)
	}

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
		// get snapshot value from the pointer | transform value to json-value-format
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		records = append(records, gValue)
	}
	var totalRecordsCount int64
	var _ = crud.GormDb.Find(modelRef).Count(&totalRecordsCount)
	// logRead
	if crud.LogRead {
		_, _ = crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: map[string]interface{}{"getType": "All Records"},
			TableName:  crud.TableName,
		})
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: GetResult{
				Records: records,
				Stats: GetStats{
					Skip:              crud.Skip,
					Limit:             crud.Limit,
					RecordsCount:      int(result.RowsAffected),
					TotalRecordsCount: int(totalRecordsCount),
				},
			},
		})
}
