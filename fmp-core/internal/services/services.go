package services

import (
	"database/sql"
	"fmt"
	"time"

	"fmp-core/internal/models"

	"github.com/google/uuid"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

// Category services
func GetCategories() ([]models.Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func CreateCategory(req models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO categories (id, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, category.ID, category.Name, category.Description, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func GetCategory(id uuid.UUID) (*models.Category, error) {
	category := &models.Category{}
	query := `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	return category, nil
}

func UpdateCategory(id uuid.UUID, req models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		UpdatedAt:   time.Now(),
	}

	query := `UPDATE categories SET name = $1, description = $2, updated_at = $3 WHERE id = $4`
	result, err := db.Exec(query, category.Name, category.Description, category.UpdatedAt, category.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("category not found")
	}

	// Get the updated category
	return GetCategory(id)
}

func DeleteCategory(id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// Transaction services
func GetTransactions(filters models.TransactionFilters) ([]models.Transaction, error) {
	query := `SELECT id, category_id, amount, description, date, created_at, updated_at FROM transactions WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, *filters.CategoryID)
		argIndex++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND date >= $%d", argIndex)
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND date <= $%d", argIndex)
		args = append(args, *filters.EndDate)
		argIndex++
	}

	query += " ORDER BY date DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.CategoryID, &transaction.Amount, &transaction.Description, &transaction.Date, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func CreateTransaction(req models.CreateTransactionRequest) (*models.Transaction, error) {
	transaction := &models.Transaction{
		ID:          uuid.New(),
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO transactions (id, category_id, amount, description, date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(query, transaction.ID, transaction.CategoryID, transaction.Amount, transaction.Description, transaction.Date, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func GetTransaction(id uuid.UUID) (*models.Transaction, error) {
	transaction := &models.Transaction{}
	query := `SELECT id, category_id, amount, description, date, created_at, updated_at FROM transactions WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&transaction.ID, &transaction.CategoryID, &transaction.Amount, &transaction.Description, &transaction.Date, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}
	return transaction, nil
}

func UpdateTransaction(id uuid.UUID, req models.CreateTransactionRequest) (*models.Transaction, error) {
	transaction := &models.Transaction{
		ID:          id,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
		UpdatedAt:   time.Now(),
	}

	query := `UPDATE transactions SET category_id = $1, amount = $2, description = $3, date = $4, updated_at = $5 WHERE id = $6`
	result, err := db.Exec(query, transaction.CategoryID, transaction.Amount, transaction.Description, transaction.Date, transaction.UpdatedAt, transaction.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("transaction not found")
	}

	return GetTransaction(id)
}

func DeleteTransaction(id uuid.UUID) error {
	query := `DELETE FROM transactions WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

// Planned Expense services
func GetPlannedExpenses(filters models.PlannedExpenseFilters) ([]models.PlannedExpense, error) {
	query := `SELECT id, category_id, amount, description, planned_date, is_completed, created_at, updated_at FROM planned_expenses WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, *filters.CategoryID)
		argIndex++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND planned_date >= $%d", argIndex)
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND planned_date <= $%d", argIndex)
		args = append(args, *filters.EndDate)
		argIndex++
	}

	if filters.IsCompleted != nil {
		query += fmt.Sprintf(" AND is_completed = $%d", argIndex)
		args = append(args, *filters.IsCompleted)
		argIndex++
	}

	query += " ORDER BY planned_date ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.PlannedExpense
	for rows.Next() {
		var expense models.PlannedExpense
		err := rows.Scan(&expense.ID, &expense.CategoryID, &expense.Amount, &expense.Description, &expense.PlannedDate, &expense.IsCompleted, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func CreatePlannedExpense(req models.CreatePlannedExpenseRequest) (*models.PlannedExpense, error) {
	expense := &models.PlannedExpense{
		ID:          uuid.New(),
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		PlannedDate: req.PlannedDate,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO planned_expenses (id, category_id, amount, description, planned_date, is_completed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(query, expense.ID, expense.CategoryID, expense.Amount, expense.Description, expense.PlannedDate, expense.IsCompleted, expense.CreatedAt, expense.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return expense, nil
}

func GetPlannedExpense(id uuid.UUID) (*models.PlannedExpense, error) {
	expense := &models.PlannedExpense{}
	query := `SELECT id, category_id, amount, description, planned_date, is_completed, created_at, updated_at FROM planned_expenses WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&expense.ID, &expense.CategoryID, &expense.Amount, &expense.Description, &expense.PlannedDate, &expense.IsCompleted, &expense.CreatedAt, &expense.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("planned expense not found")
		}
		return nil, err
	}
	return expense, nil
}

func UpdatePlannedExpense(id uuid.UUID, req models.CreatePlannedExpenseRequest) (*models.PlannedExpense, error) {
	expense := &models.PlannedExpense{
		ID:          id,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		PlannedDate: req.PlannedDate,
		UpdatedAt:   time.Now(),
	}

	query := `UPDATE planned_expenses SET category_id = $1, amount = $2, description = $3, planned_date = $4, updated_at = $5 WHERE id = $6`
	result, err := db.Exec(query, expense.CategoryID, expense.Amount, expense.Description, expense.PlannedDate, expense.UpdatedAt, expense.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("planned expense not found")
	}

	return GetPlannedExpense(id)
}

func DeletePlannedExpense(id uuid.UUID) error {
	query := `DELETE FROM planned_expenses WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("planned expense not found")
	}

	return nil
}

// Planned Income services
func GetPlannedIncome(filters models.PlannedIncomeFilters) ([]models.PlannedIncome, error) {
	query := `SELECT id, amount, description, month, year, created_at, updated_at FROM planned_incomes WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.Month != nil {
		query += fmt.Sprintf(" AND month = $%d", argIndex)
		args = append(args, *filters.Month)
		argIndex++
	}

	if filters.Year != nil {
		query += fmt.Sprintf(" AND year = $%d", argIndex)
		args = append(args, *filters.Year)
		argIndex++
	}

	query += " ORDER BY year DESC, month DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []models.PlannedIncome
	for rows.Next() {
		var income models.PlannedIncome
		err := rows.Scan(&income.ID, &income.Amount, &income.Description, &income.Month, &income.Year, &income.CreatedAt, &income.UpdatedAt)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, income)
	}

	return incomes, nil
}

func CreatePlannedIncome(req models.CreatePlannedIncomeRequest) (*models.PlannedIncome, error) {
	income := &models.PlannedIncome{
		ID:          uuid.New(),
		Amount:      req.Amount,
		Description: req.Description,
		Month:       req.Month,
		Year:        req.Year,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO planned_incomes (id, amount, description, month, year, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(query, income.ID, income.Amount, income.Description, income.Month, income.Year, income.CreatedAt, income.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return income, nil
}

func UpdatePlannedIncome(id uuid.UUID, req models.CreatePlannedIncomeRequest) (*models.PlannedIncome, error) {
	income := &models.PlannedIncome{
		ID:          id,
		Amount:      req.Amount,
		Description: req.Description,
		Month:       req.Month,
		Year:        req.Year,
		UpdatedAt:   time.Now(),
	}

	query := `UPDATE planned_incomes SET amount = $1, description = $2, month = $3, year = $4, updated_at = $5 WHERE id = $6`
	result, err := db.Exec(query, income.Amount, income.Description, income.Month, income.Year, income.UpdatedAt, income.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("planned income not found")
	}

	// Get the updated income
	query = `SELECT id, amount, description, month, year, created_at, updated_at FROM planned_incomes WHERE id = $1`
	err = db.QueryRow(query, id).Scan(&income.ID, &income.Amount, &income.Description, &income.Month, &income.Year, &income.CreatedAt, &income.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return income, nil
}

func DeletePlannedIncome(id uuid.UUID) error {
	query := `DELETE FROM planned_incomes WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("planned income not found")
	}

	return nil
}

// Category Limit services
func GetCategoryLimits(filters models.CategoryLimitFilters) ([]models.CategoryLimit, error) {
	query := `SELECT id, category_id, limit_amount, month, year, created_at, updated_at FROM category_limits WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, *filters.CategoryID)
		argIndex++
	}

	if filters.Month != nil {
		query += fmt.Sprintf(" AND month = $%d", argIndex)
		args = append(args, *filters.Month)
		argIndex++
	}

	if filters.Year != nil {
		query += fmt.Sprintf(" AND year = $%d", argIndex)
		args = append(args, *filters.Year)
		argIndex++
	}

	query += " ORDER BY year DESC, month DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var limits []models.CategoryLimit
	for rows.Next() {
		var limit models.CategoryLimit
		err := rows.Scan(&limit.ID, &limit.CategoryID, &limit.Limit, &limit.Month, &limit.Year, &limit.CreatedAt, &limit.UpdatedAt)
		if err != nil {
			return nil, err
		}
		limits = append(limits, limit)
	}

	return limits, nil
}

func CreateCategoryLimit(req models.CreateCategoryLimitRequest) (*models.CategoryLimit, error) {
	limit := &models.CategoryLimit{
		ID:         uuid.New(),
		CategoryID: req.CategoryID,
		Limit:      req.Limit,
		Month:      req.Month,
		Year:       req.Year,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	query := `INSERT INTO category_limits (id, category_id, limit_amount, month, year, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(query, limit.ID, limit.CategoryID, limit.Limit, limit.Month, limit.Year, limit.CreatedAt, limit.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return limit, nil
}

func UpdateCategoryLimit(id uuid.UUID, req models.CreateCategoryLimitRequest) (*models.CategoryLimit, error) {
	limit := &models.CategoryLimit{
		ID:         id,
		CategoryID: req.CategoryID,
		Limit:      req.Limit,
		Month:      req.Month,
		Year:       req.Year,
		UpdatedAt:  time.Now(),
	}

	query := `UPDATE category_limits SET category_id = $1, limit_amount = $2, month = $3, year = $4, updated_at = $5 WHERE id = $6`
	result, err := db.Exec(query, limit.CategoryID, limit.Limit, limit.Month, limit.Year, limit.UpdatedAt, limit.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("category limit not found")
	}

	// Get the updated limit
	query = `SELECT id, category_id, limit_amount, month, year, created_at, updated_at FROM category_limits WHERE id = $1`
	err = db.QueryRow(query, id).Scan(&limit.ID, &limit.CategoryID, &limit.Limit, &limit.Month, &limit.Year, &limit.CreatedAt, &limit.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return limit, nil
}

func DeleteCategoryLimit(id uuid.UUID) error {
	query := `DELETE FROM category_limits WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category limit not found")
	}

	return nil
}

// Analytics services
func GetMonthlySummary(month, year int) (*models.MonthlySummary, error) {
	summary := &models.MonthlySummary{
		Month: month,
		Year:  year,
	}

	// Get category summaries for the month
	query := `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COALESCE(SUM(t.amount), 0) as amount,
			cl.limit_amount as limit_amount
		FROM categories c
		LEFT JOIN transactions t ON c.id = t.category_id 
			AND EXTRACT(MONTH FROM t.date) = $1 
			AND EXTRACT(YEAR FROM t.date) = $2
		LEFT JOIN category_limits cl ON c.id = cl.category_id 
			AND cl.month = $1 
			AND cl.year = $2
		GROUP BY c.id, c.name, cl.limit_amount
		ORDER BY amount DESC
	`

	rows, err := db.Query(query, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var total float64
	for rows.Next() {
		var categorySummary models.CategorySummary
		var limitAmount sql.NullFloat64

		err := rows.Scan(&categorySummary.CategoryID, &categorySummary.CategoryName, &categorySummary.Amount, &limitAmount)
		if err != nil {
			return nil, err
		}

		if limitAmount.Valid {
			categorySummary.Limit = &limitAmount.Float64
			categorySummary.IsExceeded = categorySummary.Amount > limitAmount.Float64
		}

		summary.Categories = append(summary.Categories, categorySummary)
		total += categorySummary.Amount
	}

	summary.Total = total
	return summary, nil
}

func GetCategorySummary(filters models.TransactionFilters) ([]models.CategorySummary, error) {
	query := `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COALESCE(SUM(t.amount), 0) as amount
		FROM categories c
		LEFT JOIN transactions t ON c.id = t.category_id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND c.id = $%d", argIndex)
		args = append(args, *filters.CategoryID)
		argIndex++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND t.date >= $%d", argIndex)
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND t.date <= $%d", argIndex)
		args = append(args, *filters.EndDate)
		argIndex++
	}

	query += " GROUP BY c.id, c.name ORDER BY amount DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []models.CategorySummary
	for rows.Next() {
		var summary models.CategorySummary
		err := rows.Scan(&summary.CategoryID, &summary.CategoryName, &summary.Amount)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func GetLimitExceeded() ([]models.LimitExceeded, error) {
	query := `SELECT id, category_id, limit_amount, actual_amount, month, year, created_at FROM limit_exceeded ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.LimitExceeded
	for rows.Next() {
		var record models.LimitExceeded
		err := rows.Scan(&record.ID, &record.CategoryID, &record.Limit, &record.Actual, &record.Month, &record.Year, &record.CreatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// Notification services
func GetNotifications() ([]models.Notification, error) {
	query := `SELECT id, type, title, message, is_read, created_at FROM notifications ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(&notification.ID, &notification.Type, &notification.Title, &notification.Message, &notification.IsRead, &notification.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func CreateNotification(req models.CreateNotificationRequest) (*models.Notification, error) {
	notification := &models.Notification{
		ID:        uuid.New(),
		Type:      req.Type,
		Title:     req.Title,
		Message:   req.Message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO notifications (id, type, title, message, is_read, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, notification.ID, notification.Type, notification.Title, notification.Message, notification.IsRead, notification.CreatedAt)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func MarkNotificationAsRead(id uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func GetNotificationStats() (*models.NotificationStats, error) {
	var stats models.NotificationStats

	// Get unread count
	err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE is_read = false`).Scan(&stats.UnreadCount)
	if err != nil {
		return nil, err
	}

	// Get total count
	err = db.QueryRow(`SELECT COUNT(*) FROM notifications`).Scan(&stats.TotalCount)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func CheckDailyReminder() error {
	// Check if user has entered any transactions today
	today := time.Now().Format("2006-01-02")
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE DATE(date) = $1`, today).Scan(&count)
	if err != nil {
		return err
	}

	// If no transactions today, create a reminder
	if count == 0 {
		_, err = CreateNotification(models.CreateNotificationRequest{
			Type:    models.NotificationTypeDailyReminder,
			Title:   "Daily Reminder",
			Message: "Don't forget to log your expenses for today! 💰",
		})
		return err
	}

	return nil
}

func CheckLimitWarnings() error {
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	// Get categories with limits and their current spending
	query := `
		SELECT cl.category_id, cl.limit_amount, c.name,
		       COALESCE(SUM(t.amount), 0) as current_spending
		FROM category_limits cl
		JOIN categories c ON cl.category_id = c.id
		LEFT JOIN transactions t ON cl.category_id = t.category_id 
			AND EXTRACT(MONTH FROM t.date) = $1 
			AND EXTRACT(YEAR FROM t.date) = $2
		WHERE cl.month = $1 AND cl.year = $2
		GROUP BY cl.category_id, cl.limit_amount, c.name
		HAVING COALESCE(SUM(t.amount), 0) >= cl.limit_amount * 0.8
	`

	rows, err := db.Query(query, month, year)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryID uuid.UUID
		var limitAmount, currentSpending float64
		var categoryName string

		err := rows.Scan(&categoryID, &limitAmount, &currentSpending, &categoryName)
		if err != nil {
			continue
		}

		// Check if we're at 80% or more of the limit
		if currentSpending >= limitAmount*0.8 {
			percentage := int((currentSpending / limitAmount) * 100)

			var notificationType models.NotificationType
			if currentSpending >= limitAmount {
				notificationType = models.NotificationTypeLimitExceeded
			} else {
				notificationType = models.NotificationTypeLimitWarning
			}

			_, err = CreateNotification(models.CreateNotificationRequest{
				Type:  notificationType,
				Title: fmt.Sprintf("Limit %s", string(notificationType)),
				Message: fmt.Sprintf("Category '%s' has reached %d%% of its limit (%.2f/%.2f)",
					categoryName, percentage, currentSpending, limitAmount),
			})
			if err != nil {
				continue
			}
		}
	}

	return nil
}
