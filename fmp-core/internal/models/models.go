package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Transaction struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CategoryID  uuid.UUID `json:"category_id" db:"category_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	Date        time.Time `json:"date" db:"date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type PlannedExpense struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CategoryID  uuid.UUID `json:"category_id" db:"category_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	PlannedDate time.Time `json:"planned_date" db:"planned_date"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type PlannedIncome struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	Month       int       `json:"month" db:"month"`
	Year        int       `json:"year" db:"year"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CategoryLimit struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CategoryID uuid.UUID `json:"category_id" db:"category_id"`
	Limit      float64   `json:"limit" db:"limit"`
	Month      int       `json:"month" db:"month"`
	Year       int       `json:"year" db:"year"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type LimitExceeded struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CategoryID uuid.UUID `json:"category_id" db:"category_id"`
	Limit      float64   `json:"limit" db:"limit"`
	Actual     float64   `json:"actual" db:"actual"`
	Month      int       `json:"month" db:"month"`
	Year       int       `json:"year" db:"year"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Request/Response DTOs
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateTransactionRequest struct {
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type CreatePlannedExpenseRequest struct {
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
	Description string    `json:"description"`
	PlannedDate time.Time `json:"planned_date" binding:"required"`
}

type CreatePlannedIncomeRequest struct {
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
	Month       int     `json:"month" binding:"required"`
	Year        int     `json:"year" binding:"required"`
}

type CreateCategoryLimitRequest struct {
	CategoryID uuid.UUID `json:"category_id" binding:"required"`
	Limit      float64   `json:"limit" binding:"required"`
	Month      int       `json:"month" binding:"required"`
	Year       int       `json:"year" binding:"required"`
}

type MonthlySummary struct {
	Month      int               `json:"month"`
	Year       int               `json:"year"`
	Categories []CategorySummary `json:"categories"`
	Total      float64           `json:"total"`
}

type CategorySummary struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Amount       float64   `json:"amount"`
	Limit        *float64  `json:"limit,omitempty"`
	IsExceeded   bool      `json:"is_exceeded"`
}

// Filter types
type TransactionFilters struct {
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
}

type PlannedExpenseFilters struct {
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	IsCompleted *bool      `json:"is_completed,omitempty"`
}

type PlannedIncomeFilters struct {
	Month *int `json:"month,omitempty"`
	Year  *int `json:"year,omitempty"`
}

type CategoryLimitFilters struct {
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	Month      *int       `json:"month,omitempty"`
	Year       *int       `json:"year,omitempty"`
}

// Notification types
type NotificationType string

const (
	NotificationTypeDailyReminder  NotificationType = "daily_reminder"
	NotificationTypeLimitWarning   NotificationType = "limit_warning"
	NotificationTypeLimitExceeded  NotificationType = "limit_exceeded"
	NotificationTypeIncomeReminder NotificationType = "income_reminder"
)

type Notification struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	Type      NotificationType `json:"type" db:"type"`
	Title     string           `json:"title" db:"title"`
	Message   string           `json:"message" db:"message"`
	IsRead    bool             `json:"is_read" db:"is_read"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}

type CreateNotificationRequest struct {
	Type    NotificationType `json:"type" binding:"required"`
	Title   string           `json:"title" binding:"required"`
	Message string           `json:"message" binding:"required"`
}

type NotificationStats struct {
	UnreadCount int `json:"unread_count"`
	TotalCount  int `json:"total_count"`
}
