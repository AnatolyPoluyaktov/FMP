package api

import (
	"net/http"
	"strconv"
	"time"

	"fmp-core/internal/models"
	"fmp-core/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, db interface{}) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		// Categories
		api.GET("/categories", getCategories)
		api.POST("/categories", createCategory)
		api.GET("/categories/:id", getCategory)
		api.PUT("/categories/:id", updateCategory)
		api.DELETE("/categories/:id", deleteCategory)

		// Transactions
		api.GET("/transactions", getTransactions)
		api.POST("/transactions", createTransaction)
		api.GET("/transactions/:id", getTransaction)
		api.PUT("/transactions/:id", updateTransaction)
		api.DELETE("/transactions/:id", deleteTransaction)

		// Planned Expenses
		api.GET("/planned-expenses", getPlannedExpenses)
		api.POST("/planned-expenses", createPlannedExpense)
		api.GET("/planned-expenses/:id", getPlannedExpense)
		api.PUT("/planned-expenses/:id", updatePlannedExpense)
		api.DELETE("/planned-expenses/:id", deletePlannedExpense)

		// Planned Income
		api.GET("/planned-income", getPlannedIncome)
		api.POST("/planned-income", createPlannedIncome)
		api.PUT("/planned-income/:id", updatePlannedIncome)
		api.DELETE("/planned-income/:id", deletePlannedIncome)

		// Category Limits
		api.GET("/category-limits", getCategoryLimits)
		api.POST("/category-limits", createCategoryLimit)
		api.PUT("/category-limits/:id", updateCategoryLimit)
		api.DELETE("/category-limits/:id", deleteCategoryLimit)

		// Analytics
		api.GET("/analytics/monthly-summary", getMonthlySummary)
		api.GET("/analytics/category-summary", getCategorySummary)
		api.GET("/analytics/limit-exceeded", getLimitExceeded)

		// Notifications
		api.GET("/notifications", getNotifications)
		api.POST("/notifications", createNotification)
		api.PUT("/notifications/:id/read", markNotificationAsRead)
		api.GET("/notifications/stats", getNotificationStats)
		api.POST("/notifications/check-daily", checkDailyReminder)
		api.POST("/notifications/check-limits", checkLimitWarnings)
	}
}

// Categories handlers
// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Router /categories [get]
func getCategories(c *gin.Context) {
	categories, err := services.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Category data"
// @Success 201 {object} models.Category
// @Router /categories [post]
func createCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := services.CreateCategory(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// @Summary Get category by ID
// @Description Get category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Category
// @Router /categories/{id} [get]
func getCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	category, err := services.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary Update category
// @Description Update category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body models.CreateCategoryRequest true "Category data"
// @Success 200 {object} models.Category
// @Router /categories/{id} [put]
func updateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := services.UpdateCategory(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary Delete category
// @Description Delete category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 204
// @Router /categories/{id} [delete]
func deleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := services.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Transactions handlers
// @Summary Get all transactions
// @Description Get all transactions with optional filtering
// @Tags transactions
// @Accept json
// @Produce json
// @Param category_id query string false "Category ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} models.Transaction
// @Router /transactions [get]
func getTransactions(c *gin.Context) {
	var filters models.TransactionFilters

	if categoryID := c.Query("category_id"); categoryID != "" {
		id, err := uuid.Parse(categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		filters.CategoryID = &id
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			filters.StartDate = &date
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			filters.EndDate = &date
		}
	}

	transactions, err := services.GetTransactions(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// @Summary Create a new transaction
// @Description Create a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} models.Transaction
// @Router /transactions [post]
func createTransaction(c *gin.Context) {
	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default date to today if not provided
	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	transaction, err := services.CreateTransaction(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// @Summary Get transaction by ID
// @Description Get transaction by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Router /transactions/{id} [get]
func getTransaction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	transaction, err := services.GetTransaction(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// @Summary Update transaction
// @Description Update transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param transaction body models.CreateTransactionRequest true "Transaction data"
// @Success 200 {object} models.Transaction
// @Router /transactions/{id} [put]
func updateTransaction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := services.UpdateTransaction(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// @Summary Delete transaction
// @Description Delete transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 204
// @Router /transactions/{id} [delete]
func deleteTransaction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := services.DeleteTransaction(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Planned Expenses handlers
// @Summary Get all planned expenses
// @Description Get all planned expenses with optional filtering
// @Tags planned-expenses
// @Accept json
// @Produce json
// @Param category_id query string false "Category ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param is_completed query bool false "Is completed"
// @Success 200 {array} models.PlannedExpense
// @Router /planned-expenses [get]
func getPlannedExpenses(c *gin.Context) {
	var filters models.PlannedExpenseFilters

	if categoryID := c.Query("category_id"); categoryID != "" {
		id, err := uuid.Parse(categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		filters.CategoryID = &id
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			filters.StartDate = &date
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			filters.EndDate = &date
		}
	}

	if isCompleted := c.Query("is_completed"); isCompleted != "" {
		if completed, err := strconv.ParseBool(isCompleted); err == nil {
			filters.IsCompleted = &completed
		}
	}

	expenses, err := services.GetPlannedExpenses(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// @Summary Create a new planned expense
// @Description Create a new planned expense
// @Tags planned-expenses
// @Accept json
// @Produce json
// @Param expense body models.CreatePlannedExpenseRequest true "Planned expense data"
// @Success 201 {object} models.PlannedExpense
// @Router /planned-expenses [post]
func createPlannedExpense(c *gin.Context) {
	var req models.CreatePlannedExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense, err := services.CreatePlannedExpense(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, expense)
}

// @Summary Get planned expense by ID
// @Description Get planned expense by ID
// @Tags planned-expenses
// @Accept json
// @Produce json
// @Param id path string true "Planned expense ID"
// @Success 200 {object} models.PlannedExpense
// @Router /planned-expenses/{id} [get]
func getPlannedExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	expense, err := services.GetPlannedExpense(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// @Summary Update planned expense
// @Description Update planned expense
// @Tags planned-expenses
// @Accept json
// @Produce json
// @Param id path string true "Planned expense ID"
// @Param expense body models.CreatePlannedExpenseRequest true "Planned expense data"
// @Success 200 {object} models.PlannedExpense
// @Router /planned-expenses/{id} [put]
func updatePlannedExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.CreatePlannedExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense, err := services.UpdatePlannedExpense(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// @Summary Delete planned expense
// @Description Delete planned expense
// @Tags planned-expenses
// @Accept json
// @Produce json
// @Param id path string true "Planned expense ID"
// @Success 204
// @Router /planned-expenses/{id} [delete]
func deletePlannedExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := services.DeletePlannedExpense(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Planned Income handlers
// @Summary Get all planned income
// @Description Get all planned income with optional filtering
// @Tags planned-income
// @Accept json
// @Produce json
// @Param month query int false "Month"
// @Param year query int false "Year"
// @Success 200 {array} models.PlannedIncome
// @Router /planned-income [get]
func getPlannedIncome(c *gin.Context) {
	var filters models.PlannedIncomeFilters

	if month := c.Query("month"); month != "" {
		if m, err := strconv.Atoi(month); err == nil {
			filters.Month = &m
		}
	}

	if year := c.Query("year"); year != "" {
		if y, err := strconv.Atoi(year); err == nil {
			filters.Year = &y
		}
	}

	incomes, err := services.GetPlannedIncome(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incomes)
}

// @Summary Create a new planned income
// @Description Create a new planned income
// @Tags planned-income
// @Accept json
// @Produce json
// @Param income body models.CreatePlannedIncomeRequest true "Planned income data"
// @Success 201 {object} models.PlannedIncome
// @Router /planned-income [post]
func createPlannedIncome(c *gin.Context) {
	var req models.CreatePlannedIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	income, err := services.CreatePlannedIncome(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, income)
}

// @Summary Update planned income
// @Description Update planned income
// @Tags planned-income
// @Accept json
// @Produce json
// @Param id path string true "Planned income ID"
// @Param income body models.CreatePlannedIncomeRequest true "Planned income data"
// @Success 200 {object} models.PlannedIncome
// @Router /planned-income/{id} [put]
func updatePlannedIncome(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.CreatePlannedIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	income, err := services.UpdatePlannedIncome(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, income)
}

// @Summary Delete planned income
// @Description Delete planned income
// @Tags planned-income
// @Accept json
// @Produce json
// @Param id path string true "Planned income ID"
// @Success 204
// @Router /planned-income/{id} [delete]
func deletePlannedIncome(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := services.DeletePlannedIncome(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Category Limits handlers
// @Summary Get all category limits
// @Description Get all category limits with optional filtering
// @Tags category-limits
// @Accept json
// @Produce json
// @Param category_id query string false "Category ID"
// @Param month query int false "Month"
// @Param year query int false "Year"
// @Success 200 {array} models.CategoryLimit
// @Router /category-limits [get]
func getCategoryLimits(c *gin.Context) {
	var filters models.CategoryLimitFilters

	if categoryID := c.Query("category_id"); categoryID != "" {
		id, err := uuid.Parse(categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		filters.CategoryID = &id
	}

	if month := c.Query("month"); month != "" {
		if m, err := strconv.Atoi(month); err == nil {
			filters.Month = &m
		}
	}

	if year := c.Query("year"); year != "" {
		if y, err := strconv.Atoi(year); err == nil {
			filters.Year = &y
		}
	}

	limits, err := services.GetCategoryLimits(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, limits)
}

// @Summary Create a new category limit
// @Description Create a new category limit
// @Tags category-limits
// @Accept json
// @Produce json
// @Param limit body models.CreateCategoryLimitRequest true "Category limit data"
// @Success 201 {object} models.CategoryLimit
// @Router /category-limits [post]
func createCategoryLimit(c *gin.Context) {
	var req models.CreateCategoryLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := services.CreateCategoryLimit(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, limit)
}

// @Summary Update category limit
// @Description Update category limit
// @Tags category-limits
// @Accept json
// @Produce json
// @Param id path string true "Category limit ID"
// @Param limit body models.CreateCategoryLimitRequest true "Category limit data"
// @Success 200 {object} models.CategoryLimit
// @Router /category-limits/{id} [put]
func updateCategoryLimit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.CreateCategoryLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := services.UpdateCategoryLimit(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, limit)
}

// @Summary Delete category limit
// @Description Delete category limit
// @Tags category-limits
// @Accept json
// @Produce json
// @Param id path string true "Category limit ID"
// @Success 204
// @Router /category-limits/{id} [delete]
func deleteCategoryLimit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := services.DeleteCategoryLimit(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Analytics handlers
// @Summary Get monthly summary
// @Description Get monthly summary for a specific month and year
// @Tags analytics
// @Accept json
// @Produce json
// @Param month query int true "Month"
// @Param year query int true "Year"
// @Success 200 {object} models.MonthlySummary
// @Router /analytics/monthly-summary [get]
func getMonthlySummary(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
		return
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	summary, err := services.GetMonthlySummary(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// @Summary Get category summary
// @Description Get category summary for a specific period
// @Tags analytics
// @Accept json
// @Produce json
// @Param category_id query string false "Category ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} models.CategorySummary
// @Router /analytics/category-summary [get]
func getCategorySummary(c *gin.Context) {
	var filters models.TransactionFilters

	if categoryID := c.Query("category_id"); categoryID != "" {
		id, err := uuid.Parse(categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		filters.CategoryID = &id
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			filters.StartDate = &date
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			filters.EndDate = &date
		}
	}

	summary, err := services.GetCategorySummary(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// @Summary Get limit exceeded records
// @Description Get all limit exceeded records
// @Tags analytics
// @Accept json
// @Produce json
// @Success 200 {array} models.LimitExceeded
// @Router /analytics/limit-exceeded [get]
func getLimitExceeded(c *gin.Context) {
	records, err := services.GetLimitExceeded()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// Notification handlers
// @Summary Get all notifications
// @Description Get all notifications
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {array} models.Notification
// @Router /notifications [get]
func getNotifications(c *gin.Context) {
	notifications, err := services.GetNotifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// @Summary Create a new notification
// @Description Create a new notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body models.CreateNotificationRequest true "Notification data"
// @Success 201 {object} models.Notification
// @Router /notifications [post]
func createNotification(c *gin.Context) {
	var req models.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := services.CreateNotification(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// @Summary Mark notification as read
// @Description Mark notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200
// @Router /notifications/{id}/read [put]
func markNotificationAsRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := services.MarkNotificationAsRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Get notification statistics
// @Description Get notification statistics
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} models.NotificationStats
// @Router /notifications/stats [get]
func getNotificationStats(c *gin.Context) {
	stats, err := services.GetNotificationStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// @Summary Check daily reminder
// @Description Check if daily reminder should be sent
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200
// @Router /notifications/check-daily [post]
func checkDailyReminder(c *gin.Context) {
	if err := services.CheckDailyReminder(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Check limit warnings
// @Description Check if limit warnings should be sent
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200
// @Router /notifications/check-limits [post]
func checkLimitWarnings(c *gin.Context) {
	if err := services.CheckLimitWarnings(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
