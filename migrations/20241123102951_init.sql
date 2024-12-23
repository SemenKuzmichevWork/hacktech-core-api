-- +goose Up
-- +goose StatementBegin
-- Create User table

CREATE TYPE user_sex AS ENUM('female', 'male', 'other') ;

CREATE TYPE family_status AS ENUM('unknown', 'single', 'married');

CREATE TABLE company_positions (
    position_name TEXT PRIMARY KEY
);

INSERT INTO company_positions (position_name) VALUES ('headhunter');

CREATE TABLE company_roles (
    role_name TEXT PRIMARY KEY
);


CREATE TABLE departments (
    department_name TEXT PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    
    slack_id TEXT UNIQUE NOT NULL,
    CHECK (slack_id != ''),

    is_admin BOOLEAN NOT NULL DEFAULT FALSE,

    role TEXT REFERENCES company_roles(role_name) NOT NULL,
    position TEXT REFERENCES company_positions(position_name) NOT NULL,
    family_status family_status NOT NULL DEFAULT 'unknown',
    sex user_sex NOT NULL DEFAULT 'other',

    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE department_users (
    department_id TEXT REFERENCES departments(department_name) NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL,
    PRIMARY KEY (department_id, user_id)
);

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    is_sent BOOLEAN NOT NULL DEFAULT FALSE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    title TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    identifier TEXT NOT NULL
);

CREATE TYPE report_kind AS ENUM('event', 'business', 'project_participation', 'daily_checkups');

CREATE TABLE user_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reported_by UUID REFERENCES users(id) NOT NULL,

    kind report_kind NOT NULL,

    rating INT NOT NULL,
    CHECK (rating >= 0 AND rating <= 5),


    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE user_reports_events (
    report_id UUID REFERENCES user_reports(id) NOT NULL,
    event_id UUID REFERENCES events(id) NOT NULL,
    PRIMARY KEY (report_id, event_id)
);

CREATE TABLE KPI (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    ROADS DOUBLE PRECISION NOT NULL DEFAULT 0,
    emploe_engagement DOUBLE PRECISION NOT NULL DEFAULT 0
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DROP TABLE user_reports_events;
DROP TABLE user_reports;
DROP TABLE events;
DROP TABLE department_users;
DROP TABLE departments;
DROP TABLE users;


DROP TYPE report_kind;
DROP TYPE user_sex; 
DROP TYPE family_status;

DROP TABLE company_positions;
DROP TABLE company_roles;

DROP TABLE KPI;

-- +goose StatementEnd
