-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "cars" (
	"reg_num" varchar(64) PRIMARY KEY,
	"mark" varchar(64) NOT NULL,
	"model" varchar(64) NOT NULL,
	"year" integer NOT NULL,
	"owner_name" varchar(32) not null,
	"owner_surname" varchar(32) not null,
	"owner_patronymic" varchar(32) not null
);

CREATE INDEX reg_num_index ON cars(reg_num);
CREATE INDEX mark_index ON cars(mark);
CREATE INDEX model_index ON cars(model);
CREATE INDEX year_index ON cars(year);
CREATE INDEX owner_name_index ON cars(owner_name);
CREATE INDEX owner_surname_index ON cars(owner_surname);
CREATE INDEX owner_patronymic_index ON cars(owner_patronymic);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS reg_num_index;
DROP INDEX IF EXISTS mark_index;
DROP INDEX IF EXISTS model_index;
DROP INDEX IF EXISTS year_index;
DROP INDEX IF EXISTS owner_name_index;
DROP INDEX IF EXISTS owner_surname_index;
DROP INDEX IF EXISTS owner_patronymic_index;

DROP TABLE IF EXISTS "cars";

-- +goose StatementEnd
