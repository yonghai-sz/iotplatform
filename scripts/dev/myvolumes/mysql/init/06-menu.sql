USE `myexampledb`;

CREATE TABLE IF NOT EXISTS `menu` (
  `id` bigint unsigned NOT NULL,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `menu_key` varchar(127) NOT NULL,
  `title` varchar(127) NOT NULL,
  `has_child` enum('N','Y') NOT NULL DEFAULT 'N',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `role_menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint unsigned NOT NULL,
  `menu_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`, `menu_id`),
  KEY `idx_role_menu_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
