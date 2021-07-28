// @Author: abbeymart | Abi Akindele | @Created: 2021-06-24 | @Updated: 2021-06-24
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect - deleteRecord

package dbcrud

import (
	"fmt"
	"github.com/abbeymart/mcresponse"
)

func (crud Crud) DeleteById(modelRef interface{}, id string) mcresponse.ResponseMessage {
	var getRes mcresponse.ResponseMessage
	if crud.LogDelete {
		// get current record
		getRes = crud.GetById(modelRef, id)
	}
	// perform crud-delete task
	result := crud.GormDb.Where("id = ?", id).Unscoped().Delete(&modelRef)
	if result.Error != nil {
		return mcresponse.GetResMessage("deleteError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// LogDelete
	if crud.LogDelete {
		_, _ = crud.TransLog.AuditLog(CrudTasks().Delete, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: getRes.Value,
			TableName:  crud.TableName,
		})
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value:   nil,
		})
}

func (crud Crud) DeleteByIds(modelRef interface{}) mcresponse.ResponseMessage {
	if len(crud.RecordIds) < 0 {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: "recordIds param is required to delete-record-by-id",
				Value:   nil,
			})
	}
	var getRes mcresponse.ResponseMessage
	if crud.LogDelete {
		// get current records
		getRes = crud.GetByIds(modelRef)
	}
	// perform crud-delete task
	result := crud.GormDb.Where("id in ?", crud.RecordIds).Unscoped().Delete(&modelRef)
	if result.Error != nil {
		return mcresponse.GetResMessage("deleteError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// LogDelete
	if crud.LogDelete {
		_, _ = crud.TransLog.AuditLog(CrudTasks().Delete, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: getRes.Value,
			TableName:  crud.TableName,
		})
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value:   nil,
		})
}

func (crud Crud) DeleteByParam(modelRef interface{}) mcresponse.ResponseMessage {
	if crud.QueryParams == nil {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: "queryParams is required to delete-record-by-param",
				Value:   nil,
			})
	}
	var getRes mcresponse.ResponseMessage
	if crud.LogDelete {
		// get current records
		getRes = crud.GetByParam(modelRef)
	}
	// perform crud-delete task
	result := crud.GormDb.Where(crud.QueryParams).Unscoped().Delete(&modelRef)
	if result.Error != nil {
		return mcresponse.GetResMessage("deleteError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// LogDelete
	if crud.LogDelete {
		_, _ = crud.TransLog.AuditLog(CrudTasks().Delete, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: getRes.Value,
			TableName:  crud.TableName,
		})
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value:   nil,
		})
}
