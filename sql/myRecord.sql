-- CREATE USER 'record'@'localhost' IDENTIFIED BY '4KzyzTL9gyQpycJ9';
-- CREATE DATABASE record;
-- GRANT ALL ON record.* TO 'record'@'localhost';
-- FLUSH PRIVILEGES;

CREATE TABLE `record`.`record_day`(
    `id` INT NOT NULL AUTO_INCREMENT,
    `content` TEXT NOT NULL,
    `time` INT NOT NULL,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`)
);

CREATE TABLE `record`.`record_week`(
    `id` INT NOT NULL AUTO_INCREMENT,
    `content` TEXT NOT NULL,
    `time` INT NOT NULL,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`)
);

CREATE TABLE `record`.`record_month`(
    `id` INT NOT NULL AUTO_INCREMENT,
    `content` TEXT NOT NULL,
    `time` INT NOT NULL,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`)
);

CREATE TABLE `record`.`flag`(
    `id` INT NOT NULL AUTO_INCREMENT,
    `content` TEXT NOT NULL,
    `status` INT NOT NULL,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`)
);
