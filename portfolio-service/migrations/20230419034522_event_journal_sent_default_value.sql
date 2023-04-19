-- Modify "event_journal" table
ALTER TABLE `portfolio`.`event_journal` MODIFY COLUMN `sent` bool NOT NULL DEFAULT 0;
