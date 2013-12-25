DROP TABLE IF EXISTS `main_user`;
CREATE TABLE `main_user` (
    `uid` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`account` varchar(64) NOT NULL DEFAULT '',
	`pwd` varchar(64) NOT NULL DEFAULT '',
	PRIMARY KEY (`uid`),
	UNIQUE KEY `account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
