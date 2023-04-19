-- Create "event_journal" table
CREATE TABLE `portfolio`.`event_journal` (`id` varchar(36) NOT NULL, `timestamp` datetime(6) NOT NULL, `name` varchar(256) NOT NULL, `event_data` json NOT NULL, `sent` bool NOT NULL, PRIMARY KEY (`id`), INDEX `idx_sent_x_timestamp` (`sent`, `timestamp`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
