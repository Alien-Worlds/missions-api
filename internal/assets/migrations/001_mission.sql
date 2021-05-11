-- +migrate Up

CREATE TABLE mission
(
    mission_id     BIGINT UNIQUE,
    description   VARCHAR(255),
    name          VARCHAR(100) NOT NULL,
    boarding_time  BIGINT          NOT NULL,
    launch_time    BIGINT          NOT NULL,
    end_time       BIGINT          NOT NULL,
    duration      BIGINT          NOT NULL,
    mission_type   BIGINT          NOT NULL,
    reward        BIGINT       NOT NULL,
    spaceship_cost BIGINT       NOT NULL,
    mission_power  BIGINT       NOT NULL,
    total_ships    BIGINT       NOT NULL,
    nft_contract   bytea        NOT NULL,
    nft_token_uri   VARCHAR(20)  NOT NULL,
    PRIMARY KEY (mission_id)
);

-- +migrate Down

DROP TABLE mission;