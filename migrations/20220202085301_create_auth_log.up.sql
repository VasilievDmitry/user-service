CREATE TABLE IF NOT EXISTS `auth_log`
(
    `id`            integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
    `user_id`       varchar(40)            NOT NULL,
    `ip`            varchar(40)            NOT NULL,
    `user_agent`    varchar(256)           NOT NULL,
    `access_token`  varchar(256)           NOT NULL,
    `refresh_token` varchar(64)            NOT NULL,
    `is_active`     boolean                NOT NULL,
    `created_at`    datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    datetime DEFAULT CURRENT_TIMESTAMP,
    `expire_at`     datetime               NOT NULL,
    INDEX idx_access_token_is_active (`access_token`, `is_active`),
    INDEX idx_refresh_token_is_active (`refresh_token`, `is_active`),
    INDEX idx_user_id (user_id)
);