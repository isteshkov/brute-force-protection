-- +migrate Up
CREATE TABLE subnets
(
    id             BIGSERIAL PRIMARY KEY,
    version        INTEGER                     NOT NULL DEFAULT 1,
    created_at     TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_at     TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMP WITHOUT TIME ZONE,
    address        cidr UNIQUE                 NOT NULL,
    is_blacklisted BOOLEAN
);

-- +migrate Down
DROP TABLE subnets;
