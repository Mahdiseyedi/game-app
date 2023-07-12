
-- MYSQL 8.0 set the role value to `user` for all old records
-- Don't change the order of enum values
-- TODO find a better solution instead of keeping the order!!!
-- +migrate Up
ALTER TABLE `users` ADD COLUMN `role` ENUM('user','admin') NOT NULL;

-- +migrate Up
ALTER TABLE `users` DROP COLUMN `role`;