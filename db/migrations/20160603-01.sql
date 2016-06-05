-- +migrate Up
CREATE TABLE `follower`(
  `twitter_id` BIGINT(20),
  PRIMARY KEY (`twitter_id`)
);

CREATE TABLE `following` (
  `twitter_id` BIGINT(20),
  `since` TIMESTAMP,
  `unfollowed` BOOLEAN DEFAULT false,
  PRIMARY KEY (`twitter_id`)
);

CREATE TABLE `temp_following` (
  `twitter_id` BIGINT(20),
  PRIMARY KEY (`twitter_id`)
);

-- +migrate Down
drop table `follower`;
drop table `following`;
drop table `temp_following`;
