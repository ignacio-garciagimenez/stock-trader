-- Modify "portfolios" table
ALTER TABLE `portfolio`.`portfolios` ADD UNIQUE INDEX `idx_name` (`name`);
