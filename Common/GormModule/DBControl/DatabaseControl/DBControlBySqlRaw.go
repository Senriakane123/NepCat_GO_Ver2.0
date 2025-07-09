package DatabaseControl

import "gorm.io/gorm"

func (DB *DatabaseControl) RawQuery(sql string, values ...interface{}) (*gorm.DB, error) {
	result := DB.db.Raw(sql, values...)
	return result, result.Error
}

func (DB *DatabaseControl) ExecSQL(sql string, values ...interface{}) error {
	return DB.db.Exec(sql, values...).Error
}
