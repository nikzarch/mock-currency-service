CREATE TABLE valutes
(
    id         BIGSERIAL PRIMARY KEY,
    report_id  BIGINT         NOT NULL REFERENCES currency_reports (id) ON DELETE CASCADE,
    num_code   CHAR(3)        NOT NULL,
    char_code  CHAR(3)        NOT NULL,
    nominal    INTEGER        NOT NULL,
    name       TEXT           NOT NULL,
    value      NUMERIC(18, 4) NOT NULL,
    vunit_rate NUMERIC(18, 6) NOT NULL
);

CREATE INDEX idx_valutes_report_id ON valutes (report_id);