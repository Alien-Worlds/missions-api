-- +migrate Up

CREATE TABLE join_event
(
    transaction_id     TEXT UNIQUE,
    -- explorer_address   TEXT   NOT NULL,
    -- mission            BIGINT NOT NULL,
    -- number_ships       BIGINT NOT NULL,
    -- stake_tlm          BIGINT NOT NULL,
    -- stake_bnb          BIGINT NOT NULL,
    PRIMARY KEY (transaction_id)
  
);

-- +migrate Down

DROP TABLE join_event;
