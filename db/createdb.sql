DROP TABLE IF EXISTS `machine_root_owns`;
DROP TABLE IF EXISTS `machine_user_owns`;
DROP TABLE IF EXISTS `machine_owns`;
DROP TABLE IF EXISTS `machines`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
  `points` BIGINT unsigned NOT NULL DEFAULT 0,
  `user_owns` BIGINT unsigned NOT NULL DEFAULT 0,
  `root_owns` BIGINT unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `machines` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
  `user_flag` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `root_flag` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `difficulty` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `points` BIGINT unsigned NOT NULL DEFAULT 0,
  `user_owns` BIGINT unsigned NOT NULL DEFAULT 0,
  `root_owns` BIGINT unsigned NOT NULL DEFAULT 0,
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
