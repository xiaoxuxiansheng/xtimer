CREATE TABLE IF NOT EXISTS `timer`
(
    `id`                bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `app`               varchar(255) NOT NULL COMMENT '应用名',
    `name`              varchar(255) NOT NULL COMMENT '定时器name',
    `status`            smallint(255) NOT NULL COMMENT '定时器状态 1未激活 2激活',
    `cron`              varchar(255) NOT NULL COMMENT '定时表达式',
    `notify_http_param` json         DEFAULT NULL COMMENT 'http 参数',
    `deleted_at`        datetime     DEFAULT NULL COMMENT '删除时间',
    `created_at`        datetime     NOT NULL COMMENT '创建时间',
    `updated_at`        datetime     DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_app` (`app`,`name`) USING BTREE COMMENT 'app name 索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;