INSERT INTO currency_reports (report_date, name)
VALUES ($1, $2)
ON CONFLICT (report_date) DO
UPDATE SET name = EXCLUDED.name
RETURNING id;