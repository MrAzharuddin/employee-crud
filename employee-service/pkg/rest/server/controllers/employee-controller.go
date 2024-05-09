package controllers

import (
	"errors"
	"github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/daos/clients/sqls"
	"github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/models"
	"github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	"strconv"
)

type EmployeeController struct {
	employeeService *services.EmployeeService
}

func NewEmployeeController() (*EmployeeController, error) {
	employeeService, err := services.NewEmployeeService()
	if err != nil {
		return nil, err
	}
	return &EmployeeController{
		employeeService: employeeService,
	}, nil
}

// CreateEmployee creates a new employee for the employee service
// @Summary Creates a new employee
// @Description Creates a new employee
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Create employee"
// @Success 201 {object} models.Employee
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /employees [post]
func (employeeController *EmployeeController) CreateEmployee(context *gin.Context) {
	// validate input
	var input models.Employee
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// trigger employee creation
	employeeCreated, err := employeeController.employeeService.CreateEmployee(&input)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, employeeCreated)
}

// FetchEmployee fetches a single employee for the employee service
// @Summary Fetches a single employee
// @Description Fetches a single employee
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.Employee
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /employees/{id} [get]
func (employeeController *EmployeeController) FetchEmployee(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger employee fetching
	employee, err := employeeController.employeeService.GetEmployee(id)
	if err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(context.Request.Context())
		currentSpan.SetAttributes(attribute.String("employee.id", strconv.FormatInt(int64(employee.ID), 10)))
	}

	context.JSON(http.StatusOK, employee)
}

// FetchEmployees fetches all employees for the employee service
// @Summary Fetches all employees
// @Description Fetches all employees
// @Tags employees
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param page_size query int false "page_size"
// @Success 200 {array} models.Employee
// @Failure 500 {object} ErrorResponse
// @Router /employees [get]
func (employeeController *EmployeeController) FetchEmployees(context *gin.Context) {
	// trigger employee fetching
	query := context.Request.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(query.Get("page_size"))
	if err != nil {
		limit = 10
	}
	employees, err := employeeController.employeeService.GetEmployees(page, limit)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, employees)
}

// UpdateEmployee updates a single employee for the employee service
// @Summary Updates a single employee
// @Description Updates a single employee
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param employee body models.Employee true "Update employee"
// @Success 204 {object} interface{}
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /employees/{id} [put]
func (employeeController *EmployeeController) UpdateEmployee(context *gin.Context) {
	// validate input
	var input models.Employee
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger employee update
	if _, err := employeeController.employeeService.UpdateEmployee(id, &input); err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

// DeleteEmployee deletes a single employee for the employee service
// @Summary Deletes a single employee
// @Description Deletes a single employee
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 204 {object} interface{}
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /employees/{id} [delete]
func (employeeController *EmployeeController) DeleteEmployee(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger employee deletion
	if err := employeeController.employeeService.DeleteEmployee(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

// PushEmployee pushes multiple random employees for the employee service
// @Summary Pushes multiple random employees
// @Description Pushes multiple random employees
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body map[string]interface{} true "Push employee"
// @Success 204 {object} interface{}
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /employees/random [post]
func (employeeController *EmployeeController) PushEmployee(context *gin.Context) {
	employees := []*models.Employee{
		{
			Name:     "Rahul Gupta",
			Position: "Software Developer",
			Salary:   76000.25,
		},
		{
			Name:     "Deepika Patel",
			Position: "Marketing Specialist",
			Salary:   59000.50,
		},
		{
			Name:     "Amit Kumar",
			Position: "Accountant",
			Salary:   53000.00,
		},
		{
			Name:     "Ananya Singh",
			Position: "Customer Service Representative",
			Salary:   40500.75,
		},
		{
			Name:     "Prakash Sharma",
			Position: "Human Resources Manager",
			Salary:   82000.00,
		},
		{
			Name:     "Neha Reddy",
			Position: "Software Developer",
			Salary:   74000.25,
		},
		{
			Name:     "Sanjay Mishra",
			Position: "Marketing Specialist",
			Salary:   60000.50,
		},
		{
			Name:     "Pooja Shah",
			Position: "Accountant",
			Salary:   56000.00,
		},
		{
			Name:     "Rajesh Patel",
			Position: "Customer Service Representative",
			Salary:   43000.75,
		},
		{
			Name:     "Nisha Desai",
			Position: "Human Resources Manager",
			Salary:   83000.00,
		},
		{
			Name:     "Ajay Verma",
			Position: "Software Developer",
			Salary:   77000.25,
		},
		{
			Name:     "Shreya Singh",
			Position: "Marketing Specialist",
			Salary:   61000.50,
		},
		{
			Name:     "Manoj Kumar",
			Position: "Accountant",
			Salary:   54000.00,
		},
		{
			Name:     "Sneha Sharma",
			Position: "Customer Service Representative",
			Salary:   41000.75,
		},
		{
			Name:     "Vivek Joshi",
			Position: "Human Resources Manager",
			Salary:   84000.00,
		},
		{
			Name:     "Meera Patel",
			Position: "Software Developer",
			Salary:   78000.25,
		},
		{
			Name:     "Raj Gupta",
			Position: "Marketing Specialist",
			Salary:   62000.50,
		},
		{
			Name:     "Aarti Sharma",
			Position: "Accountant",
			Salary:   55000.00,
		},
		{
			Name:     "Anand Singh",
			Position: "Customer Service Representative",
			Salary:   42000.75,
		},
		{
			Name:     "Divya Desai",
			Position: "Human Resources Manager",
			Salary:   85000.00,
		},
		{
			Name:     "Vikas Reddy",
			Position: "Software Developer",
			Salary:   79000.25,
		},
		{
			Name:     "Kavita Mishra",
			Position: "Marketing Specialist",
			Salary:   63000.50,
		},
		{
			Name:     "Amita Shah",
			Position: "Accountant",
			Salary:   56000.00,
		},
		{
			Name:     "Vinod Patel",
			Position: "Customer Service Representative",
			Salary:   43000.75,
		},
		{
			Name:     "Neeraj Desai",
			Position: "Human Resources Manager",
			Salary:   86000.00,
		},
		{
			Name:     "Priya Gupta",
			Position: "Software Developer",
			Salary:   80000.25,
		},
		{
			Name:     "Rahul Singh",
			Position: "Marketing Specialist",
			Salary:   64000.50,
		},
		{
			Name:     "Anjali Kumar",
			Position: "Accountant",
			Salary:   57000.00,
		},
		{
			Name:     "Arjun Sharma",
			Position: "Customer Service Representative",
			Salary:   44000.75,
		},
		{
			Name:     "Kiran Desai",
			Position: "Human Resources Manager",
			Salary:   87000.00,
		},
		{
			Name:     "Sarita Patel",
			Position: "Software Developer",
			Salary:   81000.25,
		},
		{
			Name:     "Pradeep Mishra",
			Position: "Marketing Specialist",
			Salary:   65000.50,
		},
		{
			Name:     "Nandini Shah",
			Position: "Accountant",
			Salary:   58000.00,
		},
		{
			Name:     "Rajeev Singh",
			Position: "Customer Service Representative",
			Salary:   45000.75,
		},
		{
			Name:     "Gaurav Desai",
			Position: "Human Resources Manager",
			Salary:   88000.00,
		},
		{
			Name:     "Sunita Gupta",
			Position: "Software Developer",
			Salary:   82000.25,
		},
		{
			Name:     "Pawan Singh",
			Position: "Marketing Specialist",
			Salary:   66000.50,
		},
		{
			Name:     "Anita Kumar",
			Position: "Accountant",
			Salary:   59000.00,
		},
		{
			Name:     "Suresh Sharma",
			Position: "Customer Service Representative",
			Salary:   46000.75,
		},
		{
			Name:     "Deepak Desai",
			Position: "Human Resources Manager",
			Salary:   89000.00,
		},
	}

	err := employeeController.employeeService.CreateEmployees(employees)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Random employees created"})
}
