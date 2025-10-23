package notifications

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"minapp-backend/internal/config"
	"minapp-backend/internal/telegram"
)

type NotificationService struct {
	bot *telegram.Bot
	cfg *config.Config
}

func NewNotificationService(bot *telegram.Bot, cfg *config.Config) *NotificationService {
	return &NotificationService{
		bot: bot,
		cfg: cfg,
	}
}

// SendDailyReminder sends a daily reminder to check if user has added transactions
func (ns *NotificationService) SendDailyReminder(chatID int64) error {
	// Check if user has any transactions today
	today := time.Now().Format("2006-01-02")
	url := ns.cfg.FMPCoreAPIURL + "/api/v1/transactions?start_date=" + today + "&end_date=" + today

	// Make API request to check transactions
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to check transactions: %w", err)
	}
	defer resp.Body.Close()

	// Parse response to check if there are transactions
	var transactions []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return fmt.Errorf("failed to parse transactions response: %w", err)
	}

	if len(transactions) == 0 {
		message := ns.generateReminderMessage()
		if err := ns.bot.SendMessage(chatID, message); err != nil {
			return fmt.Errorf("failed to send reminder: %w", err)
		}
		log.Printf("Daily reminder sent to chat %d", chatID)
	}

	return nil
}

// SendLimitWarning sends a warning when category limit is approaching
func (ns *NotificationService) SendLimitWarning(chatID int64, categoryName string, currentAmount, limitAmount float64) error {
	percentage := (currentAmount / limitAmount) * 100

	var message string
	if percentage >= 100 {
		message = fmt.Sprintf("🚨 Лимит превышен!\n\nКатегория: %s\nПотрачено: %.2f ₽\nЛимит: %.2f ₽\nПревышение: %.2f ₽",
			categoryName, currentAmount, limitAmount, currentAmount-limitAmount)
	} else if percentage >= 80 {
		message = fmt.Sprintf("⚠️ Приближается к лимиту!\n\nКатегория: %s\nПотрачено: %.2f ₽ (%.1f%%)\nЛимит: %.2f ₽\nОсталось: %.2f ₽",
			categoryName, currentAmount, percentage, limitAmount, limitAmount-currentAmount)
	} else {
		return nil // No warning needed
	}

	if err := ns.bot.SendMessage(chatID, message); err != nil {
		return fmt.Errorf("failed to send limit warning: %w", err)
	}

	log.Printf("Limit warning sent to chat %d for category %s", chatID, categoryName)
	return nil
}

// SendPlannedExpenseReminder sends a reminder about upcoming planned expenses
func (ns *NotificationService) SendPlannedExpenseReminder(chatID int64, expenseName string, plannedDate time.Time, amount float64) error {
	daysUntil := int(time.Until(plannedDate).Hours() / 24)

	var message string
	if daysUntil == 0 {
		message = fmt.Sprintf("📅 Сегодня планируемый расход!\n\n%s\nСумма: %.2f ₽", expenseName, amount)
	} else if daysUntil == 1 {
		message = fmt.Sprintf("📅 Завтра планируемый расход!\n\n%s\nСумма: %.2f ₽", expenseName, amount)
	} else if daysUntil <= 3 {
		message = fmt.Sprintf("📅 Через %d дня планируемый расход!\n\n%s\nСумма: %.2f ₽", daysUntil, expenseName, amount)
	} else {
		return nil // Too far in the future
	}

	if err := ns.bot.SendMessage(chatID, message); err != nil {
		return fmt.Errorf("failed to send planned expense reminder: %w", err)
	}

	log.Printf("Planned expense reminder sent to chat %d for %s", chatID, expenseName)
	return nil
}

// SendIncomeCopyNotification sends notification when income is copied from previous month
func (ns *NotificationService) SendIncomeCopyNotification(chatID int64, month int, year int, amount float64) error {
	monthNames := []string{
		"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}

	message := fmt.Sprintf("📋 Доход скопирован!\n\nМесяц: %s %d\nСумма: %.2f ₽\n\nДоход был скопирован из предыдущего месяца. Проверьте и при необходимости отредактируйте.",
		monthNames[month-1], year, amount)

	if err := ns.bot.SendMessage(chatID, message); err != nil {
		return fmt.Errorf("failed to send income copy notification: %w", err)
	}

	log.Printf("Income copy notification sent to chat %d for %s %d", chatID, monthNames[month-1], year)
	return nil
}

func (ns *NotificationService) generateReminderMessage() string {
	memes := []string{
		"😴", "🤔", "💭", "📝", "💰", "⏰", "📱", "🎯",
	}

	messages := []string{
		"Напоминание!\n\nСегодня вы еще не добавили ни одной транзакции. Не забудьте записать свои расходы! 💰\n\nИспользуйте мини-приложение для быстрого добавления трат.",
		"Эй! Не забывайте про учет расходов! 📊\n\nСегодня еще нет записей. Быстро добавьте траты через мини-приложение!",
		"Финансовая дисциплина начинается с ежедневного учета! 📈\n\nДобавьте сегодняшние расходы в мини-приложении.",
		"Деньги любят счет! 💸\n\nНе забудьте записать сегодняшние траты через мини-приложение.",
		"Учет расходов - основа финансового благополучия! 💪\n\nДобавьте транзакции за сегодня.",
		"Маленькие траты складываются в большие суммы! 🔍\n\nНе пропустите ни одной покупки сегодня.",
	}

	meme := memes[time.Now().Unix()%int64(len(memes))]
	message := messages[time.Now().Unix()%int64(len(messages))]

	return fmt.Sprintf("%s %s", meme, message)
}
