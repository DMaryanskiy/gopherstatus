package server

import (
	"errors"

	"github.com/DMaryanskiy/gopherstatus/internal/storage"
)

type MockDB struct {
	CreatedService []*storage.Service
	CreatedUser    *storage.User
	User           storage.User
}

func (m *MockDB) LatestResults(limit int) ([]storage.CheckResult, error) {
	return []storage.CheckResult{
		{ServiceName: "TestService", Online: true, ResponseMS: 100},
	}, nil
}

func (m *MockDB) HistoryForService(limit int, service string) ([]storage.CheckResult, error) {
	if service == "fail" {
		return nil, errors.New("mock db failure")
	}
	return []storage.CheckResult{
		{ServiceName: service, Online: false, ResponseMS: 321},
	}, nil
}

func (m *MockDB) CreateService(service *storage.Service) error {
	m.CreatedService = append(m.CreatedService, service)
	return nil
}

func (m *MockDB) LatestResultsByID(limit int, userID uint) ([]storage.CheckResult, error) {
	return []storage.CheckResult{
		{ServiceName: "TestService", Online: true, ResponseMS: 123},
	}, nil
}

func (m *MockDB) FetchServicesByUser(userID uint) ([]storage.Service, error) {
	return []storage.Service{
		{
			Name:     "MockService",
			URL:      "https://example.com",
			Method:   "GET",
			Interval: 60,
		},
	}, nil
}

func (m *MockDB) GetUserByEmail(email string) (storage.User, error) {
	if email == "valid@example.com" {
		return m.User, nil
	}
	return storage.User{}, errors.New("user not found")
}

func (m *MockDB) CreateUser(user *storage.User) error {
	m.CreatedUser = user
	return nil
}
