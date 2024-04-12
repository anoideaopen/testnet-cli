CREATE TABLE IF NOT EXISTS request
(
    tx_id                VARCHAR(100) PRIMARY KEY,
    block_number         integer,
    started              integer,
    finished             integer,
    duration             numeric,
    validation_code      integer,
    validation_code_text VARCHAR(100),
    chaincode_status     integer
);

CREATE TABLE IF NOT EXISTS batch
(
    tx_id                 VARCHAR(100) PRIMARY KEY,
    block_number          integer,
    started               integer,
    batch_validation_code integer,
    request_tx_id         VARCHAR(100) REFERENCES request(tx_id),
    request_timestamp     integer
);
