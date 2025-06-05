package storage

type Service struct {
	ID       uint     `gorm:"primaryKey"`
	Name     string   `gorm:"not null"`
	URL      string   `gorm:"not null"`
	Method   string   `gorm:"not null"`
	Interval int      `gorm:"not null"`
	Body     string   // optional
	Headers  []Header `gorm:"foreignKey:ServiceID"`
	UserId   uint     `gorm:"foreignKey:UserID"`
}

type Header struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"not null"`
	Value     string `gorm:"not null"`
	ServiceID uint   `gorm:"not null"`
}

func (db *DB) FetchServices() ([]Service, error) {
	var services []Service
	err := db.conn.Preload("Headers").Find(&services).Error
	return services, err
}

func (db *DB) FetchServicesByUser(userID uint) ([]Service, error) {
	var services []Service
	err := db.conn.Preload("Headers").Where(
		&Service{UserId: userID},
	).Find(&services).Error
	return services, err
}

func (db *DB) CreateService(svc *Service) (error) {
	return db.conn.Create(svc).Error
}

