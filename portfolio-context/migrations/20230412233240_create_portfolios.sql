-- Add new schema named "portfolio"
CREATE DATABASE IF NOT EXISTS `portfolio`;
-- Create "portfolios" table
CREATE TABLE `portfolio`.`portfolios` (`id` varchar(36) NOT NULL, `name` varchar(30) NOT NULL, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
