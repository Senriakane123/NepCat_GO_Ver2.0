package DatabaseControl

func (DB *DatabaseControl) InsertAutoTable(data interface{}) error {
	return DB.db.Create(data).Error
}

func (DB *DatabaseControl) FindAutoTable(query interface{}, result interface{}) error {
	return DB.db.Where(query).Find(result).Error
}

func (DB *DatabaseControl) UpdateAutoTable(condition interface{}, updates interface{}) error {
	return DB.db.Model(updates).Where(condition).Updates(updates).Error
}

func (DB *DatabaseControl) DeleteAutoTable(condition interface{}, model interface{}) error {
	return DB.db.Where(condition).Delete(model).Error
}
