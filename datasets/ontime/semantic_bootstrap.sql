CREATE DATABASE IF NOT EXISTS ontime_semantic;

CREATE TABLE IF NOT EXISTS ontime_semantic.entities
(
    `entity_name` LowCardinality(String),
    `database_name` LowCardinality(String),
    `table_name` LowCardinality(String),
    `role` LowCardinality(String),
    `description` String,
    `is_active` UInt8
)
ENGINE = MergeTree
ORDER BY (entity_name);

CREATE TABLE IF NOT EXISTS ontime_semantic.joins
(
    `join_name` LowCardinality(String),
    `left_entity` LowCardinality(String),
    `right_entity` LowCardinality(String),
    `left_expression` String,
    `right_expression` String,
    `join_type` LowCardinality(String),
    `priority` UInt8,
    `is_fallback` UInt8,
    `use_case` String,
    `notes` String,
    `is_active` UInt8
)
ENGINE = MergeTree
ORDER BY (join_name, priority);

CREATE TABLE IF NOT EXISTS ontime_semantic.enrichment_rules
(
    `rule_name` LowCardinality(String),
    `source_entity` LowCardinality(String),
    `target_entity` LowCardinality(String),
    `trigger_condition` String,
    `preferred_join_name` LowCardinality(String),
    `fallback_join_name` LowCardinality(String),
    `required_fields` String,
    `allowed_in_dynamic_mode` UInt8,
    `notes` String,
    `is_active` UInt8
)
ENGINE = MergeTree
ORDER BY (rule_name);

TRUNCATE TABLE ontime_semantic.entities;
TRUNCATE TABLE ontime_semantic.joins;
TRUNCATE TABLE ontime_semantic.enrichment_rules;

INSERT INTO ontime_semantic.entities
    (entity_name, database_name, table_name, role, description, is_active)
VALUES
    ('ontime_fact', 'ontime', 'ontime', 'fact', 'Primary fact table for BTS On-Time flight records.', 1),
    ('airports_latest', 'ontime', 'airports_latest', 'current_dimension', 'Current airport dimension with airport names, coordinates, and UTC offsets.', 1);

INSERT INTO ontime_semantic.joins
    (join_name, left_entity, right_entity, left_expression, right_expression, join_type, priority, is_fallback, use_case, notes, is_active)
VALUES
    (
        'origin_airport_id_to_airports_latest',
        'ontime_fact',
        'airports_latest',
        'OriginAirportID',
        'airport_id',
        'LEFT JOIN',
        1,
        0,
        'Origin airport enrichment when airport IDs are available in the fact rows.',
        'Preferred path for origin airport lookups.',
        1
    ),
    (
        'dest_airport_id_to_airports_latest',
        'ontime_fact',
        'airports_latest',
        'DestAirportID',
        'airport_id',
        'LEFT JOIN',
        1,
        0,
        'Destination airport enrichment when airport IDs are available in the fact rows.',
        'Preferred path for destination airport lookups.',
        1
    ),
    (
        'airport_code_to_airports_latest',
        'route_or_code_result',
        'airports_latest',
        'parsed_airport_code',
        'code',
        'LEFT JOIN',
        2,
        1,
        'Fallback airport enrichment when the result only exposes airport codes or route strings.',
        'Use only when airport IDs are unavailable in the result.',
        1
    );

INSERT INTO ontime_semantic.enrichment_rules
    (rule_name, source_entity, target_entity, trigger_condition, preferred_join_name, fallback_join_name, required_fields, allowed_in_dynamic_mode, notes, is_active)
VALUES
    (
        'airport_coordinate_enrichment',
        'ontime_fact',
        'airports_latest',
        'Use when a visualization needs airport names, coordinates, or UTC offsets that are not already present in the primary result.',
        'origin_airport_id_to_airports_latest',
        'airport_code_to_airports_latest',
        'name, latitude, longitude, utc_local_time_variation',
        1,
        'Prefer airport_id joins. Use code-based enrichment only when the result exposes route strings or airport codes without airport IDs.',
        1
    );

CREATE OR REPLACE VIEW ontime_semantic.active_joins AS
SELECT
    join_name,
    left_entity,
    right_entity,
    left_expression,
    right_expression,
    join_type,
    priority,
    is_fallback,
    use_case,
    notes
FROM ontime_semantic.joins
WHERE is_active = 1
ORDER BY priority, join_name;

CREATE OR REPLACE VIEW ontime_semantic.active_enrichment_rules AS
SELECT
    rule_name,
    source_entity,
    target_entity,
    trigger_condition,
    preferred_join_name,
    fallback_join_name,
    required_fields,
    allowed_in_dynamic_mode,
    notes
FROM ontime_semantic.enrichment_rules
WHERE is_active = 1
ORDER BY rule_name;
