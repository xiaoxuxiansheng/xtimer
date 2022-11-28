CREATE TABLE IF NOT EXISTS `task`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `app`        varchar(255) NOT NULL COMMENT '应用名',
    `timer_id`   bigint(20) NOT NULL COMMENT '定时器ID',
    `output`     varchar(256) DEFAULT NULL COMMENT '执行结果',
    `run_timer`  datetime     NOT NULL COMMENT '执行时间',
    `cost_time`  int(8) DEFAULT NULL COMMENT '执行耗时',
    `status`     int(4) NOT NULL COMMENT '当前状态',
    `created_at` datetime     NOT NULL COMMENT '创建时间',
    `updated_at` datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at` datetime     DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    UNIQUE KEY `idx_def_timer` (`timer_id`,`run_timer`) USING BTREE COMMENT '定时器执行时间索引',
    KEY `idx_run_timer` (`run_timer`) COMMENT '执行时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;