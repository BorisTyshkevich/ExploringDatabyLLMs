SELECT
    hostName() AS hostname,
    version() AS version,
    currentDatabase() AS current_db,
    now() AS checked_at;
