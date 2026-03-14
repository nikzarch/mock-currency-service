INSERT INTO currency_reports (report_date, name)
VALUES ($1, $2) RETURNING id;