-- +migrate Up
CREATE TABLE subnet_list
(
    id             BIGSERIAL,
    uid            uuid PRIMARY KEY,
    version        INTEGER                     NOT NULL DEFAULT 1,
    created_at     TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_at     TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMP WITHOUT TIME ZONE,
    subnet_address cidr                        NOT NULL,
    is_blacklisted BOOLEAN
);

CREATE UNIQUE INDEX unique_address_subnet_list ON subnet_list (subnet_address);

-- +migrate Down
DROP TABLE subnet_list;
