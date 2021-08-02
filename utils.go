// @Author: abbeymart | Abi Akindele | @Created: 2021-06-24 | @Updated: 2021-06-24
// @Company: mConnect.biz | @License: MIT
// @Description: crud-utility-helper-functions

package mcgorm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"reflect"
	"strings"
)

type EmailUserNameType struct {
	Email    string
	Username string
}

// EmailUsername processes and returns the loginName as email or username
func EmailUsername(loginName string) EmailUserNameType {
	if govalidator.IsEmail(loginName) {
		return EmailUserNameType{
			Email:    loginName,
			Username: "",
		}
	}

	return EmailUserNameType{
		Email:    "",
		Username: loginName,
	}

}

func TypeOf(rec interface{}) reflect.Type {
	return reflect.TypeOf(rec)
}

// ParseRawValues process the raw rows/records from SQL-query
func ParseRawValues(rawValues [][]byte) ([]interface{}, error) {
	// variables
	var value interface{}
	var values []interface{}
	// parse the current-raw-values
	for _, val := range rawValues {
		if err := json.Unmarshal(val, &value); err != nil {
			return nil, errors.New(fmt.Sprintf("Error parsing raw-row-value: %v", err.Error()))
		} else {
			values = append(values, value)
		}
	}
	return values, nil
}

// ArrayStringContains check if a slice of string contains/includes a string value
func ArrayStringContains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

// ArrayIntContains check if a slice of int contains/includes an int value
func ArrayIntContains(arr []int, val int) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

// ArrayToSQLStringValues transforms a slice of string to SQL-string-formatted-values
func ArrayToSQLStringValues(arr []string) string {
	result := ""
	for ind, val := range arr {
		result += "'" + val + "'"
		if ind < len(arr)-1 {
			result += ", "
		}
	}
	return result
}

// JsonToStruct converts json inputs to equivalent struct data type specification
// rec must be a pointer to a type matching the jsonRec
func JsonToStruct(jsonRec []byte, rec interface{}) error {
	if err := json.Unmarshal(jsonRec, &rec); err == nil {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error converting json-to-record-format: %v", err.Error()))
	}
}

// DataToValueParam accepts only a struct type/model and returns the ActionParamType
// data camel/Pascal-case keys are converted to underscore-keys to match table-field/columns specs
func DataToValueParam(rec interface{}) (ActionParamType, error) {
	// validate recs as struct{} type
	recType := fmt.Sprintf("%v", reflect.TypeOf(rec).Kind())
	switch recType {
	case "struct":
		dataValue := ActionParamType{}
		v := reflect.ValueOf(rec)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			dataValue[govalidator.CamelCaseToUnderscore(typeOfS.Field(i).Name)] = v.Field(i).Interface()
			//fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		}
		return dataValue, nil
	default:
		return nil, errors.New(fmt.Sprintf("rec parameter must be of type struct{}"))
	}
}

// StructToMap function converts struct to map
func StructToMap(rec interface{}) (map[string]interface{}, error) {
	var mapData map[string]interface{}
	// json record
	jsonRec, err := json.Marshal(rec)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing struct to map: %v", err.Error()))
	}
	// json-to-map
	err = json.Unmarshal(jsonRec, &mapData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing struct to map: %v", err.Error()))
	}
	return mapData, nil
}

// TagField return the field-tag (e.g. table-column-name) for mcorm tag
func TagField(rec interface{}, fieldName string, tag string) (string, error) {
	// validate recs as struct{} type
	recType := fmt.Sprintf("%v", reflect.TypeOf(rec).Kind())
	switch recType {
	case "struct":
		break
	default:
		return "", errors.New(fmt.Sprintf("rec parameter must be of type struct{}"))
	}
	t := reflect.TypeOf(rec)
	// convert the first-letter to upper-case (public field)
	field, found := t.FieldByName(strings.Title(fieldName))
	if !found {
		// check private field
		field, found = t.FieldByName(fieldName)
		if !found {
			return "", errors.New(fmt.Sprintf("error retrieving tag-field for field-name: %v", fieldName))
		}
	}
	//tagValue := field.Tag
	return field.Tag.Get(tag), nil
}

// StructToTagMap function converts struct to map (for crud-actionParams / records)
func StructToTagMap(rec interface{}, tag string) (map[string]interface{}, error) {
	// validate recs as struct{} type
	recType := fmt.Sprintf("%v", reflect.TypeOf(rec).Kind())
	switch recType {
	case "struct":
		break
	default:
		return nil, errors.New(fmt.Sprintf("rec parameter must be of type struct{}"))
	}
	tagMapData := map[string]interface{}{}
	mapData, err := StructToMap(rec)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing struct to map: %v", err.Error()))
	}
	// compose tagMapData
	for key, val := range mapData {
		tagField, tagErr := TagField(rec, key, tag)
		if tagErr != nil {
			return nil, errors.New(fmt.Sprintf("error computing tag-field: %v", tagErr.Error()))
		}
		tagMapData[tagField] = val
	}
	return tagMapData, nil
}

// StructToCaseUnderscoreMap converts struct to map (for crud-actionParams / records)
func StructToCaseUnderscoreMap(rec interface{}) (map[string]interface{}, error) {
	// validate recs as struct{} type
	recType := fmt.Sprintf("%v", reflect.TypeOf(rec).Kind())
	switch recType {
	case "struct":
		break
	default:
		return nil, errors.New(fmt.Sprintf("rec parameter must be of type struct{}"))
	}

	caseUnderscoreMapData := map[string]interface{}{}
	mapData, err := StructToMap(rec)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing struct to map: %v", err.Error()))
	}
	// compose caseUnderscoreMapData
	for key, val := range mapData {
		caseUnderscoreMapData[govalidator.CamelCaseToUnderscore(key)] = val
	}
	return caseUnderscoreMapData, nil
}

// MapToUnderscoreMap converts map-fields to underscore
func MapToUnderscoreMap(rec interface{}) (map[string]interface{}, error) {
	// validate recs as map type
	recMap, ok := rec.(map[string]interface{})
	if !ok || recMap == nil {
		return nil, errors.New(fmt.Sprintf("rec parameter must be of type struct{}"))
	}

	underscoreMapData := map[string]interface{}{}
	// compose underscoreMapData
	for key, val := range recMap {
		underscoreMapData[govalidator.CamelCaseToUnderscore(key)] = val
	}
	return underscoreMapData, nil
}

// StructToFieldValues converts struct to record fields(underscore) and associated values (columns and values)
func StructToFieldValues(rec interface{}) ([]string, []interface{}, error) {
	// validate recs as struct{} type
	recType := fmt.Sprintf("%v", reflect.TypeOf(rec).Kind())
	switch recType {
	case "struct":
		break
	default:
		return nil, nil, errors.New(fmt.Sprintf("rec parameter must be of type struct{}"))
	}
	var tableFields []string
	var fieldValues []interface{}
	mapDataValue, err := StructToMap(rec)
	if err != nil {
		return nil, nil, errors.New("error computing struct to map")
	}
	// compose table fields/column(underscore) and values
	for key, val := range mapDataValue {
		tableFields = append(tableFields, govalidator.CamelCaseToUnderscore(key))
		fieldValues = append(fieldValues, val)
	}
	return tableFields, fieldValues, nil
}

// ArrayMapToStruct converts []map/actParams to []struct/model-type
func ArrayMapToStruct(actParams ActionParamsType, recs interface{}) (interface{}, error) {
	// validate recs as slice / []struct{} type
	recsType := fmt.Sprintf("%v", reflect.TypeOf(recs).Kind())
	switch recsType {
	case "slice":
		break
	default:
		return nil, errors.New(fmt.Sprintf("recs parameter must be of type []struct{}: %v", recsType))
	}
	switch rType := recs.(type) {
	case []interface{}:
		for i, val := range rType {
			// validate each record as struct type
			recType := fmt.Sprintf("%v", reflect.TypeOf(val).Kind())
			switch recType {
			case "struct":
				break
			default:
				return nil, errors.New(fmt.Sprintf("recs[%v] parameter must be of type struct{}: %v", i, recType))
			}
		}
	default:
		fmt.Printf("recs-type: %v : %v", rType)
		return nil, errors.New(fmt.Sprintf("rec parameter must be of type []struct{}: %v", rType))
	}
	// compute json records from actParams
	jsonRec, err := json.Marshal(actParams)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing map to struct records: %v", err.Error()))
	}
	// transform json records to []struct{} (recs)
	err = json.Unmarshal(jsonRec, &recs)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing map to struct records: %v", err.Error()))
	}
	return recs, nil
}

// MapToStruct converts map to struct
func MapToStruct(mapRecord map[string]interface{}, rec interface{}) (interface{}, error) {
	// validate rec as struct{} type
	recType := fmt.Sprintf("%v", reflect.TypeOf(rec).Kind())
	switch recType {
	case "struct":
		break
	default:
		return nil, errors.New(fmt.Sprintf("rec parameter must be of type struct{}"))
	}
	// compute json records from actParams (map-record)
	jsonRec, err := json.Marshal(mapRecord)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing map to struct records: %v", err.Error()))
	}
	// transform json record to struct{} (rec)
	err = json.Unmarshal(jsonRec, &rec)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing map to struct records: %v", err.Error()))
	}
	return rec, nil
}
