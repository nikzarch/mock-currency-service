CREATE TABLE currency_reports
(
    id          BIGSERIAL PRIMARY KEY,
    report_date DATE        NOT NULL UNIQUE,
    name        TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);