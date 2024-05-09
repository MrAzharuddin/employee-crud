package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/controllers"
	"github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/daos"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var employeeController, _ = controllers.NewEmployeeController()
var employeeDao, _ = daos.NewEmployeeDao()
var router = gin.Default()

func TestEmployeeController_CreateEmployee(t *testing.T) {
	router.POST("/employees", employeeController.CreateEmployee)

	// create buffer body
	var buff bytes.Buffer
	err := json.NewEncoder(&buff).Encode(map[string]interface{}{
		"id": 209,
		"name":     "John Doe",
		"position": "Software Developer",
		"salary":   268999.90,
	})
	assert.NoError(t, err)

	// create a create request
	req, err1 := http.NewRequest("POST", "/employees", &buff)
	assert.NoError(t, err1)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// check the status code
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestEmployeeController_FetchEmployee(t *testing.T) {

	router.GET("/employees", employeeController.FetchEmployee)

	req, err := http.NewRequest("POST", "/employees/random", nil)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	// check the status code

	req1, err1 := http.NewRequest("GET", "/employees/1", nil)
	assert.NoError(t, err1)
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusOK, rec1.Code)

}
// func TestEmployeeController_DeleteEmployee(t *testing.T) {
// 	// create a delete request
// 	router.DELETE("/employees/:id", employeeController.DeleteEmployee)
// 	// create a delete request
// 	req, err := http.NewRequest("DELETE", "/employees/1", nil)
// 	assert.NoError(t, err)
// 	rec1 := httptest.NewRecorder()
// 	router.ServeHTTP(rec1, req)

// 	// check the status code
// 	assert.Equal(t, http.StatusNoContent, rec1.Code)
// 	assert.Equal(t, "", rec1.Body.String())
// }
