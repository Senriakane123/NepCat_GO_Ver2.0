package DatabaseControl

// 开启事务，返回一个新的 DatabaseControl 实例
func (DB *DatabaseControl) BeginTx() *DatabaseControl {
	tx := DB.db.Begin()
	return &DatabaseControl{db: tx}
}

// 提交事务
func (DB *DatabaseControl) Commit() error {
	return DB.db.Commit().Error
}

// 回滚事务
func (DB *DatabaseControl) Rollback() error {
	return DB.db.Rollback().Error
}

func (DB *DatabaseControl) GetLastInsertID(dest interface{}) error {
	return DB.db.Raw("SELECT LAST_INSERT_ID()").Scan(dest).Error
}
