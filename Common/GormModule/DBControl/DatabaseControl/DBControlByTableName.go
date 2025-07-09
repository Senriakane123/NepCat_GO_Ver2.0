package DatabaseControl

// 参数名	用途（作用）	示例	类似 SQL
// query	查询条件，用于 SELECT 时过滤结果	map[string]interface{}{"id": 1}	WHERE id = 1
// condition	更新/删除条件，用于 UPDATE 或 DELETE 时选择目标	map[string]interface{}{"status": "inactive"}	WHERE status = 'inactive'
// updates	更新内容，用于更新哪些字段	map[string]interface{}{"name": "NewName"}	SET name = 'NewName'

//以下是使用实例
//var users []User
//err := DBHandele.FindWithTableName("users", map[string]interface{}{"status": "active"}, &users)
//对应的sql语句为：SELECT * FROM users WHERE status = 'active';

func (DB *DatabaseControl) InsertWithTableName(table string, data interface{}) error {
	return DB.db.Table(table).Create(data).Error
}

func (DB *DatabaseControl) FindWithTableName(table string, query interface{}, result interface{}) error {
	return DB.db.Table(table).Where(query).Find(result).Error
}

func (DB *DatabaseControl) UpdateWithTableName(table string, condition interface{}, updates interface{}) error {
	return DB.db.Table(table).Where(condition).Updates(updates).Error
}

func (DB *DatabaseControl) DeleteWithTableName(table string, condition interface{}) error {
	return DB.db.Table(table).Where(condition).Delete(nil).Error
}

func (DB *DatabaseControl) FindPageWithTableName(
	table string,
	query interface{},
	result interface{},
	page int,
	pageSize int,
) (total int64, err error) {
	// 统计总条数
	err = DB.db.Table(table).Where(query).Count(&total).Error
	if err != nil {
		return
	}

	// 查询分页数据
	offset := (page - 1) * pageSize
	err = DB.db.Table(table).
		Where(query).
		Limit(pageSize).
		Offset(offset).
		Find(result).Error

	return
}

//func BuildDBqueryData(tablelistname,data []string)  {
//	return make
//}
