
----------------- 以下sql支持使用mysql数据库----------------------------

CREATE DATABASE `erc20_score` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin ；

use `erc20_score`;

CREATE TABLE `t_user_balance` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_account` varchar(80) COLLATE utf8mb4_bin NOT NULL COMMENT '用户账号-钱包address',
  `balance` bigint DEFAULT NULL COMMENT '余额',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '新增时间',
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户余额表';

CREATE TABLE `t_user_score` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_account` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户账号',
  `score` decimal(20,2) DEFAULT NULL COMMENT '积分(保留两位小数)',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户积分表';

CREATE TABLE `t_transaction` (
  `id` bigint NOT NULL,
  `type` int NOT NULL COMMENT '交易类型：0-铸币，1-转移，2-销毁',
  `from_account` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '交易from账户',
  `to_account` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '交易to账户',
  `amount` bigint DEFAULT NULL COMMENT '交易金额',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户交易表（记录余额变动记录）';

CREATE TABLE `t_chain` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `chain_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '链名称',
  `chain_id` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '链id',
  `contract_address` varchar(80) COLLATE utf8mb4_bin NOT NULL COMMENT '合约部署地址',
  `syn_from_block_num` bigint NOT NULL COMMENT '开始同步区块号',
  `syned_block_num` bigint NOT NULL COMMENT '已同步区块号',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='链信息表';