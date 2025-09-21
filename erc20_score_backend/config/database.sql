
----------------- 以下sql支持使用mysql数据库----------------------------

CREATE DATABASE `erc20_score` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin ；

use `erc20_score`;


CREATE TABLE `t_chain` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `chain_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '链名称',
  `chain_id` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '链id',
  `contract_address` varchar(80) COLLATE utf8mb4_bin NOT NULL COMMENT '合约部署地址',
  `syn_from_block_num` bigint NOT NULL COMMENT '开始同步区块号',
  `syned_block_num` bigint NOT NULL COMMENT '已同步区块号',
  `syned_score_time` datetime DEFAULT NULL COMMENT '已计算分数计算时间',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='链信息表';

CREATE TABLE `t_transaction` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `chain_id` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '链id',
  `type` int NOT NULL COMMENT '交易类型：0-铸币，1-转移，2-销毁',
  `from_account` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '交易from账户',
  `to_account` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '交易to账户',
  `tx_hash` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `amount` bigint DEFAULT NULL COMMENT '交易金额',
  `block_num` bigint DEFAULT NULL COMMENT '区块号',
  `block_time` datetime DEFAULT NULL COMMENT '区块时间',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1022 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户交易表（记录余额变动记录）';

CREATE TABLE `t_user_balance` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `chain_id` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '链id',
  `user_account` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户账号-钱包address',
  `balance` bigint DEFAULT NULL COMMENT '余额',
  `block_num` bigint DEFAULT NULL COMMENT '区块编号',
  `block_time` datetime DEFAULT NULL COMMENT '区块时间/事件时间',
  `start_block_time` datetime DEFAULT NULL COMMENT '开始区块时间',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '新增时间',
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=210 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户余额表';

CREATE TABLE `t_user_balance_his` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `chain_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '链id',
  `user_account` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户账号-钱包address',
  `balance` bigint DEFAULT NULL COMMENT '余额',
  `block_num` bigint DEFAULT NULL COMMENT '区块编号',
  `block_time` datetime DEFAULT NULL COMMENT '区块时间、事件时间',
  `start_block_time` datetime DEFAULT NULL COMMENT '开始区块时间',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '新增时间',
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1877 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户余额表';

CREATE TABLE `t_user_score` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `chain_id` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '链id',
  `user_account` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户账号',
  `score` decimal(20,2) DEFAULT NULL COMMENT '积分(保留两位小数)',
  `score_time` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '得分时间(yyyy-MM-dd HH:00:00)',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=79 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户积分表';

CREATE TABLE `t_user_score_his` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `chain_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '链id',
  `user_account` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户账号',
  `score` decimal(20,2) DEFAULT NULL COMMENT '积分(保留两位小数)',
  `score_time` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '得分时间(yyyy-MM-dd HH:00:00)',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=119 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户积分表';


INSERT INTO `erc20_score`.`t_chain` (`id`, `chain_name`, `chain_id`, `contract_address`, `syn_from_block_num`, `syned_block_num`, `syned_score_time`, `create_time`, `update_time`) VALUES (1, '本地hardhat链1', '31337_1', '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512', 0, 0, '2000-01-01 00:00:00', '2025-09-21 20:37:01', '2025-09-21 20:37:01');
INSERT INTO `erc20_score`.`t_chain` (`id`, `chain_name`, `chain_id`, `contract_address`, `syn_from_block_num`, `syned_block_num`, `syned_score_time`, `create_time`, `update_time`) VALUES (2, '本地hardhat链2', '31337_2', '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512', 0, 0, '2000-01-01 00:00:00', '2025-09-21 20:37:02', '2025-09-21 20:37:02');
