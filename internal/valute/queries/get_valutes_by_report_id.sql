SELECT num_code,
       char_code,
       nominal,
       name,
       value,
       vunit_rate
FROM valutes
WHERE report_id = $1
ORDER BY char_code;