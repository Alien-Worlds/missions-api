-- +migrate Up

CREATE TABLE explorer_mission
(
    explorer_mission_id BIGINT GENERATED ALWAYS AS IDENTITY,
    explorer          BIGINT     NOT NULL,
    mission           BIGINT     NOT NULL,
    withdrawn         Boolean NOT NULL,
    number_ships       BIGINT  NOT NULL,
    total_stake_tlm     BIGINT  NOT NULL,
    total_stake_bnb     BIGINT  NOT NULL,
    PRIMARY KEY (explorer_mission_id),
    CONSTRAINT fk_explorer
        FOREIGN KEY (explorer)
            REFERENCES explorer (explorer_id)
            ON DELETE SET NULL,
    CONSTRAINT fk_mission
        FOREIGN KEY (mission)
            REFERENCES mission (mission_id)
            ON DELETE SET NULL
);

-- +migrate Down

DROP TABLE explorer_mission;