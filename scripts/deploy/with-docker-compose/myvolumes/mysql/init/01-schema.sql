CREATE DATABASE IF NOT EXISTS `myexampledb` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;

USE `myexampledb`;

CREATE TABLE IF NOT EXISTS `shorturl`
(
  `shorten` varchar(255) NOT NULL COMMENT 'shorten key',
  `url` varchar(255) NOT NULL COMMENT 'original url',
  PRIMARY KEY(`shorten`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
