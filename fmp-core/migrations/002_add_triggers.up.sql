-- +goose Up
-- Add triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_planned_expenses_updated_at BEFORE UPDATE ON planned_expenses FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_planned_incomes_updated_at BEFORE UPDATE ON planned_incomes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_category_limits_updated_at BEFORE UPDATE ON category_limits FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_category_limits_updated_at ON category_limits;
DROP TRIGGER IF EXISTS update_planned_incomes_updated_at ON planned_incomes;
DROP TRIGGER IF EXISTS update_planned_expenses_updated_at ON planned_expenses;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP FUNCTION IF EXISTS update_updated_at_column();
