# 权限模块
CREATE TABLE `admin_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '管理用户表主键',
  `hg_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '鹰角员工的工号',
  `email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '管理用户的email',
  `nick_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '管理用户昵称',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '管理用户真实姓名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '管理员用户的创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '管理员用户的更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_hg_id` (`hg_id`) COMMENT '员工工号在系统中唯一'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `admin_role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色id',
  `name` varchar(32) NOT NULL COMMENT '角色名',
  `description` varchar(256) NOT NULL COMMENT '角色描述',
  `resource_list` json NOT NULL COMMENT '角色所包含的权限点列表',
  `status` tinyint(3) unsigned NOT NULL COMMENT '角色状态，1为可用，2为已删除',
  `app_id` tinyint(3) unsigned NOT NULL COMMENT '游戏id',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '角色创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '角色更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_app_relation` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户-app_id表主键',
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `app_id` int(11) NOT NULL COMMENT '游戏id',
  `status` int(11) NOT NULL COMMENT '用户app id关系状态，1为正常，2为已封禁',
  `identity` tinyint(4) NOT NULL COMMENT '用户类型，1为一般用户，2为超级管理员',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_app_id` (`user_id`,`app_id`) COMMENT '用户和appid关系唯一'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户、appid关系表';

CREATE TABLE `user_role_relation` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户id',
  `role_id` int(10) unsigned NOT NULL COMMENT '角色id',
  `app_id` int(10) NOT NULL COMMENT '项目id，冗余字段',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`,`app_id`) USING BTREE COMMENT '用户-角色对唯一'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户、角色关系表';

CREATE TABLE `auth_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '日志主键',
  `app_id` tinyint NOT NULL COMMENT '游戏id',
  `operator_id` int unsigned NOT NULL COMMENT '操作员id',
  `operation` int unsigned NOT NULL COMMENT '操作类型，1为创建问卷，2为编辑问卷，3为编辑定时，4为取消定时，5为终止问卷,6为删除问卷',
  `content` varchar(128) NOT NULL COMMENT '权限管理日志的可读详情，每一种具体操作有一种可读的操作详情模版',
  `user_name` varchar(16) DEFAULT NULL COMMENT '用于搜索的冗余字段，被操作的用户名',
  `role_name` varchar(32) DEFAULT NULL COMMENT '用于搜索的冗余字段，被操作的角色名',
  `user_id` int DEFAULT NULL COMMENT '被操作的管理员用户的id',
  `role_id` int DEFAULT NULL COMMENT '被操作的角色的id',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '日志创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '日志更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_operator_id` (`operator_id`) USING BTREE COMMENT '用于按操作员搜索日志',
  FULLTEXT KEY `ft_user_name` (`user_name`) COMMENT '用于按被操作者名字模糊搜索日志',
  FULLTEXT KEY `ft_role_name` (`role_name`) COMMENT '用于按被操作角色名模糊搜索日志'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

