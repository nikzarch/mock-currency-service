INSERT INTO valutes (report_id, num_code, char_code, nominal, name, value, vunit_rate)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;