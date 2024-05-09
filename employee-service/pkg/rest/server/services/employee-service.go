package services

import (
	"github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/daos"
	"github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/models"
)

type EmployeeService struct {
	employeeDao *daos.EmployeeDao
}

func NewEmployeeService() (*EmployeeService, error) {
	employeeDao, err := daos.NewEmployeeDao()
	if err != nil {
		return nil, err
	}
	return &EmployeeService{
		employeeDao: employeeDao,
	}, nil
}

func (employeeService *EmployeeService) CreateEmployee(employee *models.Employee) (*models.Employee, error) {
	return employeeService.employeeDao.CreateEmployee(employee)
}

func (employeeService *EmployeeService) GetEmployee(id int64) (*models.Employee, error) {
	return employeeService.employeeDao.GetEmployee(id)
}

func (employeeService *EmployeeService) GetEmployees(page int, limit int) ([]*models.Employee, error) {
	return employeeService.employeeDao.GetEmployees(page, limit)
}

func (employeeService *EmployeeService) UpdateEmployee(id int64, employee *models.Employee) (*models.Employee, error) {
	return employeeService.employeeDao.UpdateEmployee(id, employee)
}

func (employeeService *EmployeeService) DeleteEmployee(id int64) error {
	return employeeService.employeeDao.DeleteEmployee(id)
}

func (employeeService *EmployeeService) CreateEmployees(employees []*models.Employee) (error) {
	return employeeService.employeeDao.CreateEmployees(employees)
}