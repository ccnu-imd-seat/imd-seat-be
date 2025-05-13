CREATE TABLE `seat` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '主键 ID，自增',
  `seat` VARCHAR(50) NOT NULL COMMENT '座位编号',
  `room` VARCHAR(50) NOT NULL COMMENT '房间名称',
  `date` DATE NOT NULL COMMENT '预约日期',
  `status` VARCHAR(20) NOT NULL COMMENT '状态，如 available/booked',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_seat_room_date` (`seat`, `room`, `date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `reservation` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '预约 ID',
  `type` VARCHAR(20) NOT NULL COMMENT '预约类型',
  `date` DATE NOT NULL COMMENT '预约日期',
  `room` VARCHAR(50) NOT NULL COMMENT '房间',
  `seat_id` INT NOT NULL COMMENT '座位 ID',
  `status` VARCHAR(20) NOT NULL COMMENT '预约状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
