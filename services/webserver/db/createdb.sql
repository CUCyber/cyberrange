CREATE DATABASE IF NOT EXISTS cyberrange;
use cyberrange;

DROP TABLE IF EXISTS `machine_root_owns`;
DROP TABLE IF EXISTS `machine_user_owns`;
DROP TABLE IF EXISTS `machine_owns`;
DROP TABLE IF EXISTS `machines`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
  `points` BIGINT unsigned NOT NULL DEFAULT 0,
  `user_owns` BIGINT unsigned NOT NULL DEFAULT 0,
  `root_owns` BIGINT unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `machines` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
  `type` VARCHAR(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'linux',
  `points` BIGINT unsigned NOT NULL DEFAULT 0,
  `status` VARCHAR(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'unknown',
  `difficulty` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_flag` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `root_flag` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_owns` BIGINT unsigned NOT NULL DEFAULT 0,
  `root_owns` BIGINT unsigned NOT NULL DEFAULT 0,
  `ip_address` VARCHAR(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `machine_owns` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT unsigned NOT NULL,
  `machine_id` BIGINT unsigned NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  FOREIGN KEY (`machine_id`) REFERENCES `machines` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `machine_user_owns` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `solved_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `own_id` BIGINT unsigned NOT NULL UNIQUE,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`own_id`) REFERENCES `machine_owns` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `machine_root_owns` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `solved_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `own_id` BIGINT unsigned NOT NULL UNIQUE,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`own_id`) REFERENCES `machine_owns` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
