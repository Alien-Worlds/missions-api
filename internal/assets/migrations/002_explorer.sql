CREATE TABLE Explorer
(
    explorerId    INT GENERATED ALWAYS AS IDENTITY,
    totalStakeTLM BigInt NOT NULL,
    totalStakeBNB BigInt NOT NULL,
    PRIMARY KEY (explorerId)
);