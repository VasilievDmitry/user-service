CREATE TABLE IF NOT EXISTS `auth_provider`
(
    `id`         integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
    `user_id`    varchar(40)            NOT NULL,
    `provider`   varchar(40)            NOT NULL,
    `token`      varchar(256)           NOT NULL,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX udx_provider_token (`provider`, `token`),
    INDEX idx_user_id (user_id)
);