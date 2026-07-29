-- Trivial migration so the db-init image has something real to carry.
CREATE TABLE IF NOT EXISTS fib_cache (
    n     INTEGER PRIMARY KEY,
    value BIGINT NOT NULL
);
