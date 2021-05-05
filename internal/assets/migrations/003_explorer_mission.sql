-- +migrate Up

CREATE TABLE ExplorerMission
(
    explorerMissionId INT GENERATED ALWAYS AS IDENTITY,
    explorer          INT     NOT NULL,
    mission           INT     NOT NULL,
    withdrawn         Boolean NOT NULL,
    numberShips       BigInt  NOT NULL,
    totalStakeTLM     BigInt  NOT NULL,
    totalStakeBNB     BigInt  NOT NULL,
    PRIMARY KEY (explorerMissionId),
    CONSTRAINT fk_explorer
        FOREIGN KEY (explorer)
            REFERENCES Explorer (explorerId)
            ON DELETE SET NULL,
    CONSTRAINT fk_mission
        FOREIGN KEY (mission)
            REFERENCES Mission (missionId)
            ON DELETE SET NULL
);

-- +migrate Down

DROP TABLE ExplorerMission;