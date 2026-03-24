CREATE TABLE IF NOT EXISTS applications (
    id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL,
    phone VARCHAR(20),
    organisation_name VARCHAR(150) NOT NULL,
    organisation_url VARCHAR(500),
    project_name VARCHAR(150) NOT NULL,
    type_id BIGINT NOT NULL REFERENCES project_types(id),
    expected_results TEXT NOT NULL,
    is_payed BOOLEAN NOT NULL DEFAULT FALSE,
    additional_information TEXT,
    rejected_reason VARCHAR(750),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED')),
    reviewer BIGINT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_applications_updated_at ON applications(updated_at);
CREATE INDEX idx_applications_project_name ON applications(project_name);
CREATE INDEX idx_applications_full_name ON applications(full_name);
CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_applications_type_id ON applications(type_id);
CREATE INDEX idx_applications_reviewer ON applications(reviewer);
