CREATE TABLE IF NOT EXISTS `user`
(
    `id`              varchar(40)  NOT NULL PRIMARY KEY,
    `login`           varchar(255) DEFAULT '',
    `password`        varchar(128) DEFAULT '',
    `username`        varchar(128) DEFAULT '',
    `email_code`      varchar(64)  DEFAULT '',
    `email_confirmed` boolean      DEFAULT 0,
    `is_active`       boolean      DEFAULT 1,
    `created_at`      datetime     DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      datetime     DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX udx_id (`id`),
    UNIQUE INDEX udx_login (`login`)
);