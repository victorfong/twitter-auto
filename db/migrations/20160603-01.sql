-- +migrate Up
CREATE TABLE `follower`(
  `id` BIGINT(20) AUTO_INCREMENT,
  `twitter_id` BIGINT(20),
  PRIMARY KEY (`id`)
);

CREATE TABLE `following` (
  `id` BIGINT(20) AUTO_INCREMENT,
  `twitter_id` BIGINT(20),
  `since` TIMESTAMP,
  PRIMARY KEY (`id`)
);

CREATE TABLE `temp_following` (
  `twitter_id` BIGINT(20),
  PRIMARY KEY (`twitter_id`)
);

-- +migrate Down
drop table `follower`;
drop table `following`;
drop table `temp_following`;
