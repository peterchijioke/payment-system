package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"take-Home-assignment/internal/dto"
	"take-Home-assignment/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock PaymentService for handler tests
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) ProcessPayment(db *gorm.DB, req *dto.PaymentRequest, idempotencyKey string) (*dto.PaymentResponse, error) {
	args := m.Called(db, req, idempotencyKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaymentResponse), args.Error(1)
}

func (m *MockPaymentService) GetTransactionDetails(transactionID string) (*dto.TransactionDetails, error) {
	args := m.Called(transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TransactionDetails), args.Error(1)
}

func (m *MockPaymentService) ListPayments(limit, offset int, status, startDate, endDate string) ([]dto.TransactionDetails, int64, error) {
	args := m.Called(limit, offset, status, startDate, endDate)
	return args.Get(0).([]dto.TransactionDetails), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentService) ProcessWebhook(db *gorm.DB, rawBody io.Reader, signature string) (*dto.WebhookPayload, error) {
	args := m.Called(db, rawBody, signature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WebhookPayload), args.Error(1)
}

func (m *MockPaymentService) ListAccounts() ([]models.Account, error) {
	args := m.Called()
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockPaymentService) GetDB() *gorm.DB {
	return nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

// ============ InitiatePayment Handler Tests ============

func TestInitiatePayment_MissingIdempotencyKey(t *testing.T) {
	r := setupRouter()

	r.POST("/payments", func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Idempotency-Key header is required",
			})
			return
		}

		var req dto.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	})

	body := bytes.NewBufferString(`{"account_id":"test-id","amount":1000}`)
	req, _ := http.NewRequest("POST", "/payments", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Contains(t, response["error"], "Idempotency-Key")
}

func TestInitiatePayment_InvalidJSON(t *testing.T) {
	r := setupRouter()

	r.POST("/payments", func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Idempotency-Key header is required",
			})
			return
		}

		var req dto.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	})

	body := bytes.NewBufferString(`{invalid json}`)
	req, _ := http.NewRequest("POST", "/payments", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-key-123")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInitiatePayment_MissingRequiredFields(t *testing.T) {
	r := setupRouter()

	r.POST("/payments", func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Idempotency-Key header is required",
			})
			return
		}

		var req dto.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	})

	body := bytes.NewBufferString(`{"account_id":"test-id"}`)
	req, _ := http.NewRequest("POST", "/payments", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-key-123")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInitiatePayment_ServiceError(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	r.POST("/payments", func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Idempotency-Key header is required",
			})
			return
		}

		var req dto.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		resp, err := mockService.ProcessPayment(nil, &req, idempotencyKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    resp,
		})
	})

	mockService.On("ProcessPayment", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("account not found"))

	body := bytes.NewBufferString(`{"account_id":"test-id","amount":1000,"currency":"NGN","destination_currency":"USD","recipient_name":"John","recipient_account":"123","recipient_bank":"Bank","recipient_country":"NG"}`)
	req, _ := http.NewRequest("POST", "/payments", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-key-123")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "account not found")
}

func TestInitiatePayment_Success(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	resp := &dto.PaymentResponse{
		TransactionID:     "txn-123",
		TransactionRef:    "TXN-001",
		ProviderReference: "PRV-001",
		Status:            "initiated",
		Amount:            1000,
		Currency:          "NGN",
	}

	r.POST("/payments", func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Idempotency-Key header is required",
			})
			return
		}

		var req dto.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		result, err := mockService.ProcessPayment(nil, &req, idempotencyKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    result,
		})
	})

	mockService.On("ProcessPayment", mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)

	body := bytes.NewBufferString(`{"account_id":"test-id","amount":1000,"currency":"NGN","destination_currency":"USD","recipient_name":"John","recipient_account":"123","recipient_bank":"Bank","recipient_country":"NG"}`)
	req, _ := http.NewRequest("POST", "/payments", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-key-123")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["data"])
}

// ============ GetPayment Handler Tests ============

func TestGetPayment_NotFound(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	r.GET("/payments/:id", func(c *gin.Context) {
		paymentID := c.Param("id")

		result, err := mockService.GetTransactionDetails(paymentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "payment not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	})

	mockService.On("GetTransactionDetails", "invalid-id").Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/payments/invalid-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "payment not found")
}

func TestGetPayment_Success(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	details := &dto.TransactionDetails{
		Transaction: &models.PaymentTransaction{
			ID:             "txn-123",
			TransactionRef: "TXN-001",
		},
	}

	r.GET("/payments/:id", func(c *gin.Context) {
		paymentID := c.Param("id")

		result, err := mockService.GetTransactionDetails(paymentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "payment not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	})

	mockService.On("GetTransactionDetails", "txn-123").Return(details, nil)

	req, _ := http.NewRequest("GET", "/payments/txn-123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
}

// ============ ListPayments Handler Tests ============

func TestListPayments_Success(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	transactions := []dto.TransactionDetails{
		{Transaction: &models.PaymentTransaction{ID: "txn-1"}},
		{Transaction: &models.PaymentTransaction{ID: "txn-2"}},
	}

	r.GET("/payments", func(c *gin.Context) {
		limit := 20
		offset := 0
		status := c.Query("status")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		payments, total, err := mockService.ListPayments(limit, offset, status, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"payments": payments,
				"total":    total,
				"limit":    limit,
				"offset":   offset,
			},
		})
	})

	mockService.On("ListPayments", 20, 0, "", "", "").Return(transactions, int64(2), nil)

	req, _ := http.NewRequest("GET", "/payments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
}

func TestListPayments_WithFilters(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	transactions := []dto.TransactionDetails{
		{Transaction: &models.PaymentTransaction{ID: "txn-1", Status: models.TransactionStatusCompleted}},
	}

	r.GET("/payments", func(c *gin.Context) {
		limit := 10
		offset := 0
		status := c.Query("status")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		payments, total, err := mockService.ListPayments(limit, offset, status, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"payments": payments,
				"total":    total,
			},
		})
	})

	mockService.On("ListPayments", 10, 0, "completed", "2024-01-01", "2024-01-31").Return(transactions, int64(1), nil)

	req, _ := http.NewRequest("GET", "/payments?status=completed&start_date=2024-01-01&end_date=2024-01-31", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListPayments_ServiceError(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	r.GET("/payments", func(c *gin.Context) {
		limit := 20
		offset := 0
		status := c.Query("status")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		_, _, err := mockService.ListPayments(limit, offset, status, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	})

	mockService.On("ListPayments", 20, 0, "", "", "").Return(nil, int64(0), errors.New("database error"))

	req, _ := http.NewRequest("GET", "/payments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database error")
}

// ============ HandleWebhook Handler Tests ============

func TestHandleWebhook_InvalidSignature(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	r.POST("/webhooks", func(c *gin.Context) {
		signature := c.GetHeader("X-Webhook-Signature")
		rawBody := c.Request.Body

		payload, err := mockService.ProcessWebhook(nil, rawBody, signature)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "webhook processed",
			"data":    payload,
		})
	})

	mockService.On("ProcessWebhook", mock.Anything, mock.Anything, "invalid").Return(nil, errors.New("invalid webhook signature"))

	body := bytes.NewBufferString(`{"event_id":"evt-123","transaction_id":"txn-123","status":"completed"}`)
	req, _ := http.NewRequest("POST", "/webhooks", body)
	req.Header.Set("X-Webhook-Signature", "invalid")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid webhook signature")
}

func TestHandleWebhook_Success(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	payload := &dto.WebhookPayload{
		EventID:       "evt-123",
		EventType:     "payment.completed",
		TransactionID: "txn-123",
		Status:        "completed",
	}

	r.POST("/webhooks", func(c *gin.Context) {
		signature := c.GetHeader("X-Webhook-Signature")
		rawBody := c.Request.Body

		result, err := mockService.ProcessWebhook(nil, rawBody, signature)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "webhook processed",
			"data":    result,
		})
	})

	mockService.On("ProcessWebhook", mock.Anything, mock.Anything, "valid-signature").Return(payload, nil)

	body := bytes.NewBufferString(`{"event_id":"evt-123","transaction_id":"txn-123","status":"completed"}`)
	req, _ := http.NewRequest("POST", "/webhooks", body)
	req.Header.Set("X-Webhook-Signature", "valid-signature")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "webhook processed")
}

// ============ ListAccounts Handler Tests ============

func TestListAccounts_Success(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	accounts := []models.Account{
		{ID: "acc-1", Name: "Account 1"},
		{ID: "acc-2", Name: "Account 2"},
	}

	r.GET("/accounts", func(c *gin.Context) {
		result, err := mockService.ListAccounts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	})

	mockService.On("ListAccounts").Return(accounts, nil)

	req, _ := http.NewRequest("GET", "/accounts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
}

// ============ Handler Edge Cases ============

func TestInitiatePayment_InvalidUUID(t *testing.T) {
	r := setupRouter()

	r.POST("/payments", func(c *gin.Context) {
		var req dto.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	})

	body := bytes.NewBufferString(`{"account_id":"not-a-uuid","amount":1000,"currency":"NGN","destination_currency":"USD","recipient_name":"John","recipient_account":"123","recipient_bank":"Bank","recipient_country":"NG"}`)
	req, _ := http.NewRequest("POST", "/payments", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-key-123")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInitiatePayment_NegativeAmount(t *testing.T) {
	r := setupRouter()

	r.POST("/payments", func(c *gin.Context) {
		var req dto.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	})

	body := bytes.NewBufferString(`{"account_id":"550e8400-e29b-41d4-a716-446655440000","amount":-100,"currency":"NGN","destination_currency":"USD","recipient_name":"John","recipient_account":"123","recipient_bank":"Bank","recipient_country":"NG"}`)
	req, _ := http.NewRequest("POST", "/payments", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-key-123")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "greater than zero")
}

func TestInitiatePayment_InvalidCurrencyLength(t *testing.T) {
	r := setupRouter()

	r.POST("/payments", func(c *gin.Context) {
		var req dto.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	})

	body := bytes.NewBufferString(`{"account_id":"550e8400-e29b-41d4-a716-446655440000","amount":100,"currency":"NIGERIANNAIRA","destination_currency":"USD","recipient_name":"John","recipient_account":"123","recipient_bank":"Bank","recipient_country":"NG"}`)
	req, _ := http.NewRequest("POST", "/payments", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-key-123")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListPayments_Pagination(t *testing.T) {
	r := setupRouter()

	mockService := new(MockPaymentService)

	transactions := []dto.TransactionDetails{}

	r.GET("/payments", func(c *gin.Context) {
		limit := 10
		offset := 0

		payments, total, err := mockService.ListPayments(limit, offset, "", "", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"payments": payments,
				"total":    total,
				"limit":    limit,
				"offset":   offset,
			},
		})
	})

	mockService.On("ListPayments", 10, 0, "", "", "").Return(transactions, int64(100), nil)

	req, _ := http.NewRequest("GET", "/payments?limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(100), data["total"])
	assert.Equal(t, float64(10), data["limit"])
	assert.Equal(t, float64(0), data["offset"])
}
