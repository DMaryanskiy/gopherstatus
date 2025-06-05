package storage

import "time"

type CheckResult struct {
	ID          uint `gorm:"primaryKey"`
	ServiceName string
	URL         string
	Method      string
	Online      bool
	ResponseMS  int64
	CheckedAt   time.Time
	Error       string
}

func (db *DB) SaveResult(r *CheckResult) error {
	return db.conn.Create(r).Error
}

func (db *DB) LatestResults(limit int) ([]CheckResult, error) {
	var results []CheckResult
	err := db.conn.Limit(limit).Order("checked_at desc").Find(&results).Error
	return results, err
}

func (db *DB) LatestResultsByID(limit int, userID uint) ([]CheckResult, error) {
	var services []string
	err := db.conn.Model(&Service{}).Select("name").Where(
		&Service{UserId: userID},
	).Pluck("name", &services).Error
	if err != nil {
		return nil, err
	}

	var results []CheckResult
	err = db.conn.Limit(limit).Where("service_name IN ?", services).Order("checked_at desc").Find(&results).Error
	return results, err
}

func (db *DB) HistoryForService(limit int, service string) ([]CheckResult, error) {
	var results []CheckResult
	err := db.conn.Where(
		"service_name = ?",
		service,
	).Limit(limit).Order("checked_at desc").Find(&results).Error
	return results, err
}
