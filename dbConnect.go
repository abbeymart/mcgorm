// @Author: abbeymart | Abi Akindele | @Created: 2021-07-06 | @Updated: 2021-07-06
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package mcgorm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GormDb(dbConfig DbConfig) (db *gorm.DB, err error) {
	sslMode := dbConfig.SecureOption.SslMode
	sslCert := dbConfig.SecureOption.SecureCert
	if sslMode == "" {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf(`port=%d host=%s user=%s password=%s dbname=%s sslmode=%v sslrootcert=%v TimeZone=%v`,
		dbConfig.Port, dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.DbName, sslMode, sslCert, dbConfig.Timezone)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
