package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api.poster.com/pkg/response"
	"github.com/gin-gonic/gin"
)

func TestSendSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"key": "value"}
	response.SendSuccess(c, http.StatusOK, "Success Message", testData)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res["success"] != true {
		t.Errorf("expected success to be true, got %v", res["success"])
	}

	if res["message"] != "Success Message" {
		t.Errorf("expected message to be 'Success Message', got '%v'", res["message"])
	}

	dataVal, ok := res["data"].(map[string]interface{})
	if !ok || dataVal["key"] != "value" {
		t.Errorf("expected data to contain key: value, got %v", res["data"])
	}

	statusVal, ok := res["status"].(float64)
	if !ok || int(statusVal) != http.StatusOK {
		t.Errorf("expected status to be %d, got %v", http.StatusOK, res["status"])
	}

	if _, ok := res["meta"]; ok {
		t.Errorf("expected meta to be omitted, but it exists")
	}

	if _, ok := res["errors"]; ok {
		t.Errorf("expected errors to be omitted, but it exists")
	}
}

func TestSendSuccessWithMeta(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := []string{"item1", "item2"}
	meta := response.PaginationMeta{
		Page:       1,
		Limit:      10,
		Total:      2,
		TotalPages: 1,
	}
	response.SendSuccessWithMeta(c, http.StatusOK, "Success With Meta", testData, meta)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res["success"] != true {
		t.Errorf("expected success to be true, got %v", res["success"])
	}

	if res["message"] != "Success With Meta" {
		t.Errorf("expected message to be 'Success With Meta', got '%v'", res["message"])
	}

	statusVal, ok := res["status"].(float64)
	if !ok || int(statusVal) != http.StatusOK {
		t.Errorf("expected status to be %d, got %v", http.StatusOK, res["status"])
	}

	metaVal, ok := res["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected meta to be map[string]interface{}, got %v", res["meta"])
	}

	if metaVal["page"] != float64(1) || metaVal["limit"] != float64(10) || metaVal["total"] != float64(2) || metaVal["total_pages"] != float64(1) {
		t.Errorf("unexpected meta values: %v", metaVal)
	}

	if _, ok := res["errors"]; ok {
		t.Errorf("expected errors to be omitted, but it exists")
	}
}

func TestSendError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.SendError(c, http.StatusBadRequest, "Error Occurred", "some error detail")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res["success"] != false {
		t.Errorf("expected success to be false, got %v", res["success"])
	}

	if res["message"] != "Error Occurred" {
		t.Errorf("expected message to be 'Error Occurred', got '%v'", res["message"])
	}

	if res["errors"] != "some error detail" {
		t.Errorf("expected errors to be 'some error detail', got '%v'", res["errors"])
	}

	statusVal, ok := res["status"].(float64)
	if !ok || int(statusVal) != http.StatusBadRequest {
		t.Errorf("expected status to be %d, got %v", http.StatusBadRequest, res["status"])
	}

	if _, ok := res["data"]; ok {
		t.Errorf("expected data to be omitted, but it exists")
	}

	if _, ok := res["meta"]; ok {
		t.Errorf("expected meta to be omitted, but it exists")
	}
}
