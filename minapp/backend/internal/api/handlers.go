package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"minapp-backend/internal/config"
	"minapp-backend/internal/telegram"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionRequest struct {
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CategoryLimitRequest struct {
	CategoryID uuid.UUID `json:"category_id" binding:"required"`
	Limit      float64   `json:"limit" binding:"required"`
	Month      int       `json:"month" binding:"required"`
	Year       int       `json:"year" binding:"required"`
}

func SetupRoutes(router *gin.Engine, bot *telegram.Bot, cfg *config.Config) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", cfg.FrontendURL)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := router.Group("/api")
	{
		// Telegram Mini App routes
		api.POST("/webhook", handleWebhook(bot, cfg))
		api.GET("/categories", getCategories(cfg))
		api.POST("/categories", createCategory(cfg))
		api.GET("/transactions", getTransactions(cfg))
		api.POST("/transactions", createTransaction(cfg))
		api.GET("/category-limits", getCategoryLimits(cfg))
		api.POST("/category-limits", createCategoryLimit(cfg))
		api.GET("/monthly-summary", getMonthlySummary(cfg))
		api.POST("/notifications/daily-reminder", sendDailyReminder(bot, cfg))

		// Planned Expenses
		api.GET("/planned-expenses", getPlannedExpenses(cfg))
		api.POST("/planned-expenses", createPlannedExpense(cfg))
		api.PUT("/planned-expenses/:id", updatePlannedExpense(cfg))
		api.DELETE("/planned-expenses/:id", deletePlannedExpense(cfg))

		// Planned Income
		api.GET("/planned-income", getPlannedIncome(cfg))
		api.POST("/planned-income", createPlannedIncome(cfg))
		api.PUT("/planned-income/:id", updatePlannedIncome(cfg))
		api.DELETE("/planned-income/:id", deletePlannedIncome(cfg))
	}
}

func handleWebhook(bot *telegram.Bot, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var update struct {
			Message struct {
				Chat struct {
					ID int64 `json:"id"`
				} `json:"chat"`
				Text string `json:"text"`
			} `json:"message"`
		}

		if err := c.ShouldBindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Handle different commands
		switch update.Message.Text {
		case "/start":
			message := "Привет! 👋\n\nДобро пожаловать в Financial Manager Platform!\n\nИспользуйте мини-приложение для управления финансами:\n• Добавление транзакций\n• Установка лимитов по категориям\n• Просмотр аналитики\n\nНажмите кнопку ниже, чтобы открыть мини-приложение:"
			if err := bot.SendMessage(update.Message.Chat.ID, message); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		case "/help":
			message := "📚 Помощь по командам:\n\n/start - Начать работу с ботом\n/help - Показать эту справку\n/stats - Показать статистику за текущий месяц\n\nДля полного функционала используйте мини-приложение!"
			if err := bot.SendMessage(update.Message.Chat.ID, message); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		case "/stats":
			// Get monthly summary for current month
			now := time.Now()
			summary, err := getMonthlySummaryFromAPI(cfg, int(now.Month()), now.Year())
			if err != nil {
				message := "❌ Не удалось получить статистику. Попробуйте позже."
				bot.SendMessage(update.Message.Chat.ID, message)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			message := formatMonthlySummary(summary)
			if err := bot.SendMessage(update.Message.Chat.ID, message); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		default:
			message := "🤔 Неизвестная команда. Используйте /help для получения справки."
			bot.SendMessage(update.Message.Chat.ID, message)
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func getCategories(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/categories", "GET", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}

func createCategory(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/categories", "POST", req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, category)
	}
}

func getTransactions(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := cfg.FMPCoreAPIURL + "/api/v1/transactions"

		// Add query parameters
		if categoryID := c.Query("category_id"); categoryID != "" {
			url += "?category_id=" + categoryID
		}
		if startDate := c.Query("start_date"); startDate != "" {
			if categoryID := c.Query("category_id"); categoryID != "" {
				url += "&start_date=" + startDate
			} else {
				url += "?start_date=" + startDate
			}
		}
		if endDate := c.Query("end_date"); endDate != "" {
			separator := "?"
			if categoryID := c.Query("category_id"); categoryID != "" || c.Query("start_date") != "" {
				separator = "&"
			}
			url += separator + "end_date=" + endDate
		}

		transactions, err := makeAPIRequest(url, "GET", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, transactions)
	}
}

func createTransaction(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set default date to today if not provided
		if req.Date.IsZero() {
			req.Date = time.Now()
		}

		transaction, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/transactions", "POST", req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, transaction)
	}
}

func getCategoryLimits(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := cfg.FMPCoreAPIURL + "/api/v1/category-limits"

		// Add query parameters
		if categoryID := c.Query("category_id"); categoryID != "" {
			url += "?category_id=" + categoryID
		}
		if month := c.Query("month"); month != "" {
			separator := "?"
			if categoryID := c.Query("category_id"); categoryID != "" {
				separator = "&"
			}
			url += separator + "month=" + month
		}
		if year := c.Query("year"); year != "" {
			separator := "?"
			if categoryID := c.Query("category_id"); categoryID != "" || c.Query("month") != "" {
				separator = "&"
			}
			url += separator + "year=" + year
		}

		limits, err := makeAPIRequest(url, "GET", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, limits)
	}
}

func createCategoryLimit(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CategoryLimitRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		limit, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/category-limits", "POST", req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, limit)
	}
}

func getMonthlySummary(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		summary, err := getMonthlySummaryFromAPI(cfg, month, year)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, summary)
	}
}

func sendDailyReminder(bot *telegram.Bot, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ChatID int64 `json:"chat_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user has any transactions today
		today := time.Now().Format("2006-01-02")
		url := cfg.FMPCoreAPIURL + "/api/v1/transactions?start_date=" + today + "&end_date=" + today

		transactions, err := makeAPIRequest(url, "GET", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Convert to slice to check length
		transactionsSlice, ok := transactions.([]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response format"})
			return
		}

		if len(transactionsSlice) == 0 {
			memes := []string{
				"😴",
				"🤔",
				"💭",
				"📝",
				"💰",
				"⏰",
				"📱",
				"🎯",
			}

			meme := memes[time.Now().Unix()%int64(len(memes))]

			messages := []string{
				fmt.Sprintf("%s Напоминание!\n\nСегодня вы еще не добавили ни одной транзакции. Не забудьте записать свои расходы! 💰\n\nИспользуйте мини-приложение для быстрого добавления трат.", meme),
				fmt.Sprintf("%s Эй! Не забывайте про учет расходов! 📊\n\nСегодня еще нет записей. Быстро добавьте траты через мини-приложение!", meme),
				fmt.Sprintf("%s Финансовая дисциплина начинается с ежедневного учета! 📈\n\nДобавьте сегодняшние расходы в мини-приложении.", meme),
				fmt.Sprintf("%s Деньги любят счет! 💸\n\nНе забудьте записать сегодняшние траты через мини-приложение.", meme),
			}

			message := messages[time.Now().Unix()%int64(len(messages))]

			if err := bot.SendMessage(req.ChatID, message); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "reminder_sent"})
	}
}

// Helper functions
func makeAPIRequest(url, method string, body interface{}) (interface{}, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	var req *http.Request
	var err error

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func getMonthlySummaryFromAPI(cfg *config.Config, month, year int) (interface{}, error) {
	url := cfg.FMPCoreAPIURL + "/api/v1/analytics/monthly-summary?month=" + strconv.Itoa(month) + "&year=" + strconv.Itoa(year)
	return makeAPIRequest(url, "GET", nil)
}

func formatMonthlySummary(summary interface{}) string {
	// This is a simplified implementation
	// In a real application, you would format the actual summary data
	return "📊 Статистика за текущий месяц:\n\n" +
		"• Общие расходы: 0 ₽\n" +
		"• Категории:\n" +
		"  - Еда: 0 ₽\n" +
		"  - Транспорт: 0 ₽\n" +
		"  - Развлечения: 0 ₽\n\n" +
		"Используйте мини-приложение для детальной аналитики! 📱"
}

// Planned Expenses handlers
func getPlannedExpenses(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := cfg.FMPCoreAPIURL + "/api/v1/planned-expenses"

		// Add query parameters
		if categoryID := c.Query("category_id"); categoryID != "" {
			url += "?category_id=" + categoryID
		}
		if startDate := c.Query("start_date"); startDate != "" {
			separator := "?"
			if categoryID := c.Query("category_id"); categoryID != "" {
				separator = "&"
			}
			url += separator + "start_date=" + startDate
		}
		if endDate := c.Query("end_date"); endDate != "" {
			separator := "?"
			if categoryID := c.Query("category_id"); categoryID != "" || c.Query("start_date") != "" {
				separator = "&"
			}
			url += separator + "end_date=" + endDate
		}
		if isCompleted := c.Query("is_completed"); isCompleted != "" {
			separator := "?"
			if categoryID := c.Query("category_id"); categoryID != "" || c.Query("start_date") != "" || c.Query("end_date") != "" {
				separator = "&"
			}
			url += separator + "is_completed=" + isCompleted
		}

		expenses, err := makeAPIRequest(url, "GET", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, expenses)
	}
}

func createPlannedExpense(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			CategoryID  uuid.UUID `json:"category_id" binding:"required"`
			Amount      float64   `json:"amount" binding:"required"`
			Description string    `json:"description"`
			PlannedDate time.Time `json:"planned_date" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		expense, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/planned-expenses", "POST", req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, expense)
	}
}

func updatePlannedExpense(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			CategoryID  uuid.UUID `json:"category_id" binding:"required"`
			Amount      float64   `json:"amount" binding:"required"`
			Description string    `json:"description"`
			PlannedDate time.Time `json:"planned_date" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		expense, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/planned-expenses/"+id, "PUT", req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, expense)
	}
}

func deletePlannedExpense(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/planned-expenses/"+id, "DELETE", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// Planned Income handlers
func getPlannedIncome(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := cfg.FMPCoreAPIURL + "/api/v1/planned-income"

		// Add query parameters
		if month := c.Query("month"); month != "" {
			url += "?month=" + month
		}
		if year := c.Query("year"); year != "" {
			separator := "?"
			if month := c.Query("month"); month != "" {
				separator = "&"
			}
			url += separator + "year=" + year
		}

		incomes, err := makeAPIRequest(url, "GET", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, incomes)
	}
}

func createPlannedIncome(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Amount      float64 `json:"amount" binding:"required"`
			Description string  `json:"description"`
			Month       int     `json:"month" binding:"required"`
			Year        int     `json:"year" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		income, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/planned-income", "POST", req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, income)
	}
}

func updatePlannedIncome(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Amount      float64 `json:"amount" binding:"required"`
			Description string  `json:"description"`
			Month       int     `json:"month" binding:"required"`
			Year        int     `json:"year" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		income, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/planned-income/"+id, "PUT", req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, income)
	}
}

func deletePlannedIncome(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := makeAPIRequest(cfg.FMPCoreAPIURL+"/api/v1/planned-income/"+id, "DELETE", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
