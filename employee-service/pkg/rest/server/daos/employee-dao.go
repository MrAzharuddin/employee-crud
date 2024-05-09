package daos

import (
	"errors"
	"github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/daos/clients/sqls"
	"github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeDao struct {
	db *gorm.DB
}

func NewEmployeeDao() (*EmployeeDao, error) {
	sqlClient, err := sqls.InitGORMSQLiteDB()
	if err != nil {
		return nil, err
	}
	err = sqlClient.DB.AutoMigrate(models.Employee{})
	if err != nil {
		return nil, err
	}
	return &EmployeeDao{
		db: sqlClient.DB,
	}, nil
}

func (employeeDao *EmployeeDao) CreateEmployee(m *models.Employee) (*models.Employee, error) {
	if err := employeeDao.db.Create(&m).Error; err != nil {
		log.Debugf("failed to create employee: %v", err)
		return nil, err
	}

	log.Debugf("employee created")
	return m, nil
}

func (employeeDao *EmployeeDao) GetEmployee(id int64) (*models.Employee, error) {
	var m *models.Employee
	if err := employeeDao.db.Where("id = ?", id).First(&m).Error; err != nil {
		log.Debugf("failed to get employee: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}
	log.Debugf("employee retrieved")
	return m, nil
}

func (employeeDao *EmployeeDao) GetEmployees(page int, limit int) ([]*models.Employee, error) {
	var m []*models.Employee
	if err := employeeDao.db.Offset((page - 1) * limit).Limit(limit).Find(&m).Error; err != nil {
		log.Debugf("failed to get employees: %v", err)
		return nil, err
	}
	log.Debugf("employees retrieved")
	return m, nil
}

func (employeeDao *EmployeeDao) UpdateEmployee(id int64, m *models.Employee) (*models.Employee, error) {
	if id == 0 {
		return nil, errors.New("invalid employee ID")
	}
	if id != int64(m.ID) {
		return nil, errors.New("id and payload don't match")
	}

	var employee *models.Employee
	if err := employeeDao.db.Where("id = ?", id).First(&employee).Error; err != nil {
		log.Debugf("failed to find employee for update: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	if err := employeeDao.db.Save(&m).Error; err != nil {
		log.Debugf("failed to update employee: %v", err)
		return nil, err
	}
	log.Debugf("employee updated")
	return m, nil
}

func (employeeDao *EmployeeDao) DeleteEmployee(id int64) error {
	var m *models.Employee
	if err := employeeDao.db.Where("id = ?", id).Delete(&m).Error; err != nil {
		log.Debugf("failed to delete employee: %v", err)
		return err
	}

	log.Debugf("employee deleted")
	return nil
}

func (employeeDao *EmployeeDao) CreateEmployees(employees []*models.Employee) error {
	if err := employeeDao.db.Create(&employees).Error; err != nil {
		log.Debugf("failed to create employees: %v", err)
		return err
	}
	log.Debugf("employees created")
	return nil
}