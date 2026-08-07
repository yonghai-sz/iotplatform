CREATE TABLE IF NOT EXISTS `api` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(127) DEFAULT NULL,
  `path` varchar(127) DEFAULT NULL,
  `action` varchar(16) DEFAULT NULL,
  `type` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
