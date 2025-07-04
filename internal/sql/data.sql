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
  `student_id` VARCHAR(50) NOT NULL COMMENT '学号',
  `type` VARCHAR(20) NOT NULL COMMENT '预约类型',
  `date` DATE NOT NULL COMMENT '预约日期',
  `room` VARCHAR(50) NOT NULL COMMENT '房间',
  `seat` VARCHAR(50) NOT NULL COMMENT '座位编号',
  `status` VARCHAR(20) NOT NULL COMMENT '预约状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `room` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '房间 ID',
  `room` VARCHAR(50) NOT NULL COMMENT '房间',
  `status` VARCHAR(20) NOT NULL COMMENT '预约状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user` (
  `student_id` CHAR(15) NOT NULL COMMENT '学号',
  `score` INT NOT NULL DEFAULT 100 COMMENT '信誉分',
  PRIMARY KEY (`student_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `feedback` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `student_id` CHAR(15) NOT NULL COMMENT '学号',
  `content` VARCHAR(100) NOT NULL COMMENT '反馈内容最多一百字',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;