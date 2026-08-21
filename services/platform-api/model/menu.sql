CREATE TABLE IF NOT EXISTS `menu` (
  `id` bigint unsigned NOT NULL,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `menu_key` varchar(127) NOT NULL,
  `title` varchar(127) NOT NULL,
  `has_child` enum('N','Y') NOT NULL DEFAULT 'N',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
