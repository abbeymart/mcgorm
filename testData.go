// @Author: abbeymart | Abi Akindele | @Created: 2020-12-28 | @Updated: 2020-12-28
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package mcgorm

import (
	"time"
)

type Group struct {
	BaseModelType
	Name string `json:"name" gorm:"unique" mcorm:"name"`
}

type Category struct {
	BaseModelType
	Name      string    `json:"name"  mcorm:"name"`
	OwnerId   string    `json:"ownerId" mcorm:"owner_id"`
	Path      string    `json:"path" mcorm:"path"`
	Priority  uint      `json:"priority" mcorm:"priority"`
	ParentId  *string   `json:"parentId" mcorm:"parent_id"`
	GroupId   string    `json:"groupId" mcorm:"group_id"`
	Group     Group     `json:"group" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" mcorm:"group"`
	Parent    *Category `json:"parent" gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" mcorm:"parent"`
	IconStyle string    `json:"iconStyle" mcorm:"icon_style"`
}

const TestTable1 = "groups"
const TestTable2 = "categories"
const TestTable3 = "audits_test1"
const DeleteAllTable = "audits_test2"
const TestAuditTable = "audits"

// TODO: records for insert and updates: audits, groups, categories

// TODO: update and delete params (ids, queryParams)

const UserId = "085f48c5-8763-4e22-a1c6-ac1a68ba07de"

var TestUserInfo = UserInfoType{
	UserId:    "085f48c5-8763-4e22-a1c6-ac1a68ba07de",
	LoginName: "abbeymart",
	Email:     "abbeya1@yahoo.com",
	Language:  "en-US",
	Firstname: "Abi",
	Lastname:  "Akindele",
	Token:     "",
	Expire:    0,
	Role:      "TBD",
}

// TODO: audit-records

type TestParam struct {
	Name     string  `json:"name" mcorm:"name"`
	Desc     string  `json:"desc" mcorm:"desc"`
	Url      string  `json:"url" mcorm:"url"`
	Priority int     `json:"priority" mcorm:"priority"`
	Cost     float64 `json:"cost" mcorm:"cost"`
}

var Recs = TestParam{Name: "Abi", Desc: "Testing only", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var TableRecords, _ = DataToValueParam(Recs)

var NewRecs = TestParam{Name: "Abi Akindele", Desc: "Testing only - updated", Url: "localhost:9900", Priority: 1, Cost: 2000.00}
var NewTableRecords, _ = DataToValueParam(NewRecs)

var TestCrudParamOptions = CrudOptionsType{
	AuditTable:    "audits",
	UserTable:     "users",
	ProfileTable:  "profiles",
	ServiceTable:  "services",
	AccessTable:   "access_keys",
	VerifyTable:   "verify_users",
	RoleTable:     "roles",
	LogCrud:       true,
	LogCreate:     true,
	LogUpdate:     true,
	LogDelete:     true,
	LogRead:       true,
	LogLogin:      true,
	LogLogout:     true,
	MaxQueryLimit: 100000,
	MsgFrom:       "support@mconnect.biz",
}

// CreateRecordA create record(s)
var CreateRecordA = Group{
	Name: "services",
}
var CreateRecordB = Category{
	Name: "services",
}
var valParam1, _ = DataToValueParam(CreateRecordA)
var valParam2, _ = DataToValueParam(CreateRecordB)
var CreateActionParams = ActionParamsType{
	valParam1,
	valParam2,
}

// UpdateRecordType update record(s)
type UpdateRecordType struct {
	Id            string
	TableName     string
	LogRecords    interface{}
	NewLogRecords interface{}
	LogBy         string
	LogType       string
	LogAt         time.Time
}

var upRecs = TestParam{Name: "Abi100", Desc: "Testing only100", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var upTableRecords, _ = DataToValueParam(upRecs)
var upRecs2 = TestParam{Name: "Abi200", Desc: "Testing only200", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var upTableRecords2, _ = DataToValueParam(upRecs2)
var UpdateRecordA = UpdateRecordType{
	Id:            "d46a29db-a9a3-47b9-9598-e17a7338e474",
	TableName:     "services",
	LogRecords:    upTableRecords,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       UpdateLog,
	LogAt:         time.Now(),
}
var UpdateRecordB = UpdateRecordType{
	Id:            "8fcdc5d5-f4e3-4f98-ba19-16e798f81070",
	TableName:     "services2",
	LogRecords:    upTableRecords2,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       UpdateLog,
	LogAt:         time.Now(),
}

var UpdateRecordById = UpdateRecordType{
	TableName:     "services2",
	LogRecords:    upTableRecords,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       UpdateLog,
	LogAt:         time.Now(),
}

var UpdateRecordByParam = Group{
	Name: "services3",
}

var UpdateIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var UpdateParams = QueryParamType{
}

var updateRec1, _ = DataToValueParam(UpdateRecordA)
var updateRec2, _ = DataToValueParam(UpdateRecordB)
var updateRecId, _ = DataToValueParam(UpdateRecordById)
var updateRecParam, _ = DataToValueParam(UpdateRecordByParam)

var UpdateActionParams = ActionParamsType{
	updateRec1,
	updateRec2,
}

var UpdateActionParamsById = ActionParamsType{
	updateRecId,
}
var UpdateActionParamsByParam = ActionParamsType{
	updateRecParam,
}

// GetRecordType get record(s)
type GetRecordType struct {
	Id            string
	TableName     string
	LogRecords    interface{}
	NewLogRecords interface{}
	LogBy         string
	LogType       string
	LogAt         time.Time
}

// GetIds get by ids & params
var GetIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var GetParams = QueryParamType{
}

// DeleteIds delete record(s) by ids & params
var DeleteIds = []string{"dba4adbb-4482-4f3d-bb05-0db80c30876b", "02f83bc1-8fa3-432a-8432-709f0df3f3b0"}
var DeleteParams = QueryParamType{

}
