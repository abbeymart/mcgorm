// @Author: abbeymart | Abi Akindele | @Created: 2021-06-24 | @Updated: 2021-06-24
// @Company: mConnect.biz | @License: MIT
// @Description: crud - instance and methods

package mcgorm

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcresponse"
)

// Crud object / struct
type Crud struct {
	CrudParamsType
	CrudOptionsType
	CurrentRecords []interface{}
	TransLog       LogParam
	CacheKey       string // Unique for exactly the same query
}

// NewCrud constructor returns a new crud-instance
func NewCrud(params CrudParamsType, options CrudOptionsType) (crudInstance *Crud) {
	crudInstance = &Crud{}
	// compute crud params
	crudInstance.AppDb = params.AppDb
	crudInstance.GormDb = params.GormDb
	crudInstance.TableName = params.TableName
	crudInstance.UserInfo = params.UserInfo
	crudInstance.ActionParams = params.ActionParams
	crudInstance.RecordIds = params.RecordIds
	crudInstance.QueryParams = params.QueryParams
	crudInstance.SortParams = params.SortParams
	crudInstance.ProjectParams = params.ProjectParams
	crudInstance.Token = params.Token
	crudInstance.TaskName = params.TaskName
	crudInstance.Skip = params.Skip
	crudInstance.Limit = params.Limit

	// crud options
	crudInstance.MaxQueryLimit = options.MaxQueryLimit
	crudInstance.AuditTable = options.AuditTable
	crudInstance.AccessTable = options.AccessTable
	crudInstance.RoleTable = options.RoleTable
	crudInstance.UserTable = options.UserTable
	crudInstance.ProfileTable = options.ProfileTable
	crudInstance.ServiceTable = options.ServiceTable
	crudInstance.AuditDb = options.AuditDb
	crudInstance.AccessDb = options.AccessDb
	crudInstance.ServiceDb = options.ServiceDb
	crudInstance.GormAuditDb = options.GormAuditDb
	crudInstance.GormAccessDb = options.GormAccessDb
	crudInstance.GormServiceDb = options.GormServiceDb
	crudInstance.LogCrud = options.LogCrud
	crudInstance.LogRead = options.LogRead
	crudInstance.LogCreate = options.LogCreate
	crudInstance.LogUpdate = options.LogUpdate
	crudInstance.LogDelete = options.LogDelete
	crudInstance.CheckAccess = options.CheckAccess // Dec 09/2020: user to implement auth as a middleware
	crudInstance.CacheExpire = options.CacheExpire // cache expire in secs
	// Compute CacheKey from TableName, QueryParams, SortParams, ProjectParams and RecordIds
	qParam, _ := json.Marshal(params.QueryParams)
	sParam, _ := json.Marshal(params.SortParams)
	pParam, _ := json.Marshal(params.ProjectParams)
	dIds, _ := json.Marshal(params.RecordIds)
	crudInstance.CacheKey = params.TableName + string(qParam) + string(sParam) + string(pParam) + string(dIds)

	// Default values
	if crudInstance.AuditTable == "" {
		crudInstance.AuditTable = "audits"
	}
	if crudInstance.AccessTable == "" {
		crudInstance.AccessTable = "accesses"
	}
	if crudInstance.RoleTable == "" {
		crudInstance.RoleTable = "roles"
	}
	if crudInstance.UserTable == "" {
		crudInstance.UserTable = "users"
	}
	if crudInstance.ProfileTable == "" {
		crudInstance.ProfileTable = "profiles"
	}
	if crudInstance.ServiceTable == "" {
		crudInstance.ServiceTable = "services"
	}
	if crudInstance.AuditDb == nil {
		crudInstance.AuditDb = crudInstance.AppDb
	}
	if crudInstance.AccessDb == nil {
		crudInstance.AccessDb = crudInstance.AppDb
	}
	if crudInstance.GormAuditDb == nil {
		crudInstance.GormAuditDb = crudInstance.GormDb
	}
	if crudInstance.GormAccessDb == nil {
		crudInstance.GormAccessDb = crudInstance.GormDb
	}
	if crudInstance.GormServiceDb == nil {
		crudInstance.GormServiceDb = crudInstance.GormDb
	}
	if crudInstance.Skip < 0 {
		crudInstance.Skip = 0
	}
	if crudInstance.MaxQueryLimit == 0 {
		crudInstance.MaxQueryLimit = 10000
	}
	if crudInstance.Limit > crudInstance.MaxQueryLimit && crudInstance.MaxQueryLimit != 0 {
		crudInstance.Limit = crudInstance.MaxQueryLimit
	}
	if crudInstance.CacheExpire <= 0 {
		crudInstance.CacheExpire = 300 // 300 secs, 5 minutes
	}
	// Audit/TransLog instance
	crudInstance.TransLog = NewAuditLog(crudInstance.GormAuditDb, crudInstance.AuditTable)

	return crudInstance
}

// String() function implementation for crud instance/object
func (crud Crud) String() string {
	return fmt.Sprintf("CRUD Instance Information: %#v \n\n", crud)
}

// Methods

// SaveRecord function creates new record(s) or updates existing record(s)
func (crud *Crud) SaveRecord(modelRef interface{}, recs interface{}, batch int) mcresponse.ResponseMessage {
	// default value
	if batch == 0 {
		batch = 10000
	}
	// create/insert new record(s)
	if crud.TaskType == CrudTasks().Create {
		// check task-permission
		if crud.CheckAccess {
			accessRes := crud.TaskPermission(crud.TaskType)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		// save-record(s): create/insert new record(s): len(recordIds) = 0 && len(createRecs) > 0
		return crud.CreateBatch(recs, batch)
	}

	if crud.TaskType == CrudTasks().Update {
		// check task-permission
		if crud.CheckAccess {
			accessRes := crud.TaskPermission(crud.TaskType)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		// update 1 or more records by ids or queryParams
		if len(crud.ActionParams) == 1 {
			upRec := recs.([]interface{})[0]
			// update record(s) by recordIds
			if len(crud.RecordIds) > 1 {
				return crud.UpdateByIds(modelRef, upRec)
			}
			// update the record by recordId
			if len(crud.RecordIds) == 1 {
				return crud.UpdateById(modelRef, upRec, crud.RecordIds[0])
			}
			// update record(s) by queryParams
			if len(crud.QueryParams) > 0 {
				return crud.UpdateByParam(modelRef, upRec)
			}
		}
		// update multiple records
		if len(crud.ActionParams) > 1 {
			return crud.Update(modelRef, recs)
		}
	}

	// otherwise return saveError
	return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
		Message: "Save error: incomplete or invalid action/query-params provided",
		Value:   nil,
	})
}

// SaveRecord1 function creates new record(s) or updates existing record(s)
func (crud *Crud) SaveRecord1(modelRef interface{}, recs interface{}, batch int) mcresponse.ResponseMessage {
	// default value
	if batch == 0 {
		batch = 10000
	}
	//  compute taskType-records from actionParams: create or update
	var (
		createRecs ActionParamsType // records without id field-value
		updateRecs ActionParamsType // records with id field-value
		recIds     []string         // capture recordIds for separate/multiple updates
	)
	for _, record := range crud.ActionParams {
		// determine if record exists (update) or is new (create)
		if fieldValue, ok := record["id"]; ok && fieldValue != "" {
			// validate fieldValue as string
			switch fieldValue.(type) {
			case string:
				updateRecs = append(updateRecs, record)
				recIds = append(recIds, fieldValue.(string))
			default:
				// invalid fieldValue type (string)
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Invalid fieldValue type for fieldName: id, in record: %v", record),
					Value:   nil,
				})
			}
		} else if len(crud.ActionParams) == 1 && (len(crud.RecordIds) > 0 || len(crud.QueryParams) > 0) {
			updateRecs = append(updateRecs, record)
		} else {
			createRecs = append(createRecs, record)
		}
	}

	// permit only create or update, not both at the same time
	if len(createRecs) > 0 && len(updateRecs) > 0 {
		return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
			Message: "You may only create or update record(s), not both at the same time",
			Value:   nil,
		})
	}

	// create/insert new record(s)
	if len(createRecs) > 0 {
		// check task-permission - create
		crud.TaskType = CrudTasks().Create
		if crud.CheckAccess {
			accessRes := crud.TaskPermission(crud.TaskType)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		// save-record(s): create/insert new record(s): len(recordIds) = 0 && len(createRecs) > 0
		return crud.CreateBatch(recs, batch)
	}

	// check task-permission - for update task
	crud.TaskType = CrudTasks().Update
	if crud.CheckAccess {
		accessRes := crud.TaskPermission(crud.TaskType)
		if accessRes.Code != "success" {
			return accessRes
		}
	}
	// update 1 or more records by ids or queryParams
	if len(updateRecs) == 1 {
		// update the record by recordId
		upRec := recs.([]interface{})[0]
		if len(crud.RecordIds) == 1 {
			return crud.UpdateById(modelRef, upRec, crud.RecordIds[0])
		}
		// update record(s) by recordIds
		if len(crud.RecordIds) > 1 {
			return crud.UpdateByIds(modelRef, upRec)
		}
		// update record(s) by queryParams
		if len(crud.QueryParams) > 0 {
			return crud.UpdateByParam(modelRef, upRec)
		}
	}
	// update multiple records
	if len(updateRecs) > 1 {
		return crud.Update(modelRef, recs)
	}
	// otherwise return saveError
	return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
		Message: "Save error: incomplete or invalid action/query-params provided",
		Value:   nil,
	})
}

// DeleteRecord function deletes/removes record(s) by id(s) or params
func (crud *Crud) DeleteRecord(modelRef interface{}) mcresponse.ResponseMessage {
	// check task-permission - delete
	if crud.CheckAccess {
		accessRes := crud.TaskPermission(CrudTasks().Delete)
		if accessRes.Code != "success" {
			return accessRes
		}
	}
	if len(crud.RecordIds) == 1 {
		return crud.DeleteById(modelRef, crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		return crud.DeleteByIds(modelRef)
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		return crud.DeleteByParam(modelRef)
	}
	// delete-all ***RESTRICTED***
	// otherwise return error
	return mcresponse.GetResMessage("removeError", mcresponse.ResponseMessageOptions{
		Message: "Remove error: incomplete or invalid query-conditions provided",
		Value:   nil,
	})
}

// GetRecord function get records by id, params or all
func (crud *Crud) GetRecord(modelRef interface{}) mcresponse.ResponseMessage {
	// check task-permission - get/read
	if crud.CheckAccess {
		accessRes := crud.TaskPermission(CrudTasks().Read)
		if accessRes.Code != "success" {
			return accessRes
		}
	}
	if len(crud.RecordIds) == 1 {
		return crud.GetById(modelRef, crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		return crud.GetByIds(modelRef)
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		return crud.GetByParam(modelRef)
	}
	return crud.GetAll(modelRef)
}

// GetRecords function get records by id, params or all - lookup-items
func (crud *Crud) GetRecords(modelRef interface{}) mcresponse.ResponseMessage {
	if len(crud.RecordIds) == 1 {
		return crud.GetById(modelRef, crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		return crud.GetByIds(modelRef)
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		return crud.GetByParam(modelRef)
	}
	return crud.GetAll(modelRef)
}
