-- +migrate Up

CREATE TABLE explorer
(
    explorer_id    BIGINT GENERATED ALWAYS AS IDENTITY,
    explorer_address TEXT NOT NULL UNIQUE,
    total_stake_tlm BIGINT NOT NULL,
    total_stake_bnb BIGINT NOT NULL,
    PRIMARY KEY (explorer_id)
);

-- +migrate Down

DROP TABLE explorer;