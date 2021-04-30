CREATE TABLE Mission
(
    missionId     INT GENERATED ALWAYS AS IDENTITY,
    description   VARCHAR(255),
    name          VARCHAR(100) NOT NULL,
    boardingTime  INT          NOT NULL,
    launchTime    INT          NOT NULL,
    endTime       INT          NOT NULL,
    duration      INT          NOT NULL,
    missionType   INT          NOT NULL,
    reward        BigInt       NOT NULL,
    spaceshipCost BigInt       NOT NULL,
    missionPower  BigInt       NOT NULL,
    totalShips    BigInt       NOT NULL,
    nftContract   bytea        NOT NULL,
    nftTokenURI   VARCHAR(20)  NOT NULL,
    PRIMARY KEY (missionId)
);