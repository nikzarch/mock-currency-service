SELECT id, report_date, name
FROM currency_reports
WHERE report_date = $1;