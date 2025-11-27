WITH params AS (
    SELECT
        81              AS site_id,              -- укажи нужный site_id
        27              AS group_id,             -- нужный group_id или NULL
        102             AS filter_group_id,      -- нужный filter_group_id или NULL
        DATE '2025-11-24' AS end_date,           -- «сегодня»
        (DATE '2025-11-24' - INTERVAL '3 years')::date AS start_date,
        10::int         AS rank_value,           -- можно изменить
        'https://example.com/'::text AS url_value,
        'Синтетическая позиция'::text AS title_value,
        'google'::text  AS source_value,
        'desktop'::text AS device_value,
        'windows'::text AS os_value,
        FALSE           AS ads_value,
        'ru'::text      AS country_value,
        'ru'::text      AS lang_value,
        1::int          AS pages_valu   e,
        NULL::varchar   AS wordstat_query_type
),
scoped_keywords AS (
    SELECT k.id AS keyword_id, p.site_id
    FROM keywords k
    JOIN params p ON p.site_id = k.site_id
    WHERE p.group_id IS NULL OR k.group_id = p.group_id
),
calendar AS (
    SELECT generate_series(p.start_date, p.end_date, INTERVAL '1 day')::date AS day
    FROM params p
)
INSERT INTO positions (
    keyword_id,
    site_id,
    rank,
    url,
    title,
    source,
    device,
    os,
    ads,
    country,
    lang,
    pages,
    date,
    filter_group_id,
    wordstat_query_type
)
SELECT
    k.keyword_id,
    p.site_id,
    p.rank_value,
    p.url_value,
    p.title_value,
    p.source_value,
    p.device_value,
    p.os_value,
    p.ads_value,
    p.country_value,
    p.lang_value,
    p.pages_value,
    c.day,
    p.filter_group_id,
    p.wordstat_query_type
FROM scoped_keywords k
CROSS JOIN calendar c
CROSS JOIN params p;