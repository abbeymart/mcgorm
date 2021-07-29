// @Author: abbeymart | Abi Akindele | @Created: 2021-07-10 | @Updated: 2021-07-10
// @Company: mConnect.biz | @License: MIT
// @Description: crud-helper methods

package mcgorm

import (
	"fmt"
	"github.com/abbeymart/mcresponse"
)

func (crud *Crud) ApiGetRecord(modelRef interface{}) mcresponse.ResponseMessage {
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

func (crud *Crud) ApiSaveRecord(modelRef interface{}, recs interface{}, batch int) mcresponse.ResponseMessage {
	return crud.SaveRecord(modelRef, recs, batch)
}

func (crud *Crud) ApiDeleteRecord(modelRef interface{}) mcresponse.ResponseMessage {
	if len(crud.RecordIds) == 1 {
		return crud.DeleteById(modelRef, crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		return crud.DeleteByIds(modelRef)
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		return crud.DeleteByParam(modelRef)
	}
	return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Records can be deleted by Id or QueryParams only."),
		Value:   nil,
	})
}


func GetParamsMessage(msgObject MessageObject, msgType string) mcresponse.ResponseMessage {
	var messages = ""

	for key, val := range msgObject {
		if messages != "" {
			messages = fmt.Sprintf("%v | %v : %v", messages, key, val)
		} else {
			messages = fmt.Sprintf("%v : %v", key, val)
		}
	}
	if msgType == "" {
		msgType = "unknown"
	}
	return mcresponse.GetResMessage(msgType, mcresponse.ResponseMessageOptions{
		Message: messages,
		Value:   nil,
	})
}

