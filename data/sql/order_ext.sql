-- 订单扩展信息表（参考 carsaas order_ext 结构）
CREATE TABLE IF NOT EXISTS `order_ext` (
    `id`                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `order_no`             VARCHAR(64)     NOT NULL                COMMENT '订单编号',
    `bill_distance`        INT UNSIGNED    NOT NULL DEFAULT 0      COMMENT '计费距离(米)',
    `bill_duration`        INT UNSIGNED    NOT NULL DEFAULT 0      COMMENT '计费时长(秒)',
    `end_after_order`      TINYINT UNSIGNED NOT NULL DEFAULT 0     COMMENT '完单后是否继续听单 0否 1是',
    `estimate_price_info`  TEXT            NOT NULL                COMMENT '预估计费详情JSON',
    `price_info`           TEXT            NOT NULL                COMMENT '乘客计费详情JSON',
    `driver_price_info`    TEXT            NOT NULL                COMMENT '司机计费详情JSON',
    `split_info`           TEXT            NOT NULL                COMMENT '分账比例信息',
    `extra_info`           TEXT            NOT NULL                COMMENT '扩展信息JSON',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单扩展信息表';

-- 模拟订单扩展数据，与 order_info 中的模拟订单对应
INSERT IGNORE INTO `order_ext` (
    `order_no`, `bill_distance`, `bill_duration`, `end_after_order`,
    `estimate_price_info`, `price_info`, `driver_price_info`, `split_info`, `extra_info`
) VALUES
(
    'WO202611060001', 32000, 3000, 1,
    '{"total_price":5800,"start_price":1300,"mileage_price":3000,"duration_price":1000,"long_price":500,"bill_distance":32000,"bill_duration":3000}',
    '{"total_price":5800,"start_price":1300,"mileage_price":3000,"duration_price":1000,"long_price":500,"bill_distance":32000,"bill_duration":3000,"charge_detail":{"start_fee":{"start_fee":1300,"mileage":3000,"duration":600}},"third_discount_amount":600}',
    '{"total_price":4800,"start_price":1000,"mileage_price":2600,"duration_price":800,"long_price":400,"bill_distance":32000,"bill_duration":3000}',
    '{"sp_rate":0.80,"platform_rate":0.20}',
    '{"price_rule_version":{"1":"v2026.11"},"station_type":1,"pick_up_distance":1200,"delivery_distance":32000,"cross_city_type":0,"commission_free_order":0}'
),
(
    'WO202611060002', 18000, 1620, 0,
    '{"total_price":3200,"start_price":1300,"mileage_price":1400,"duration_price":500,"bill_distance":18000,"bill_duration":1620}',
    '{"total_price":3200,"start_price":1300,"mileage_price":1400,"duration_price":500,"bill_distance":18000,"bill_duration":1620,"charge_detail":{"start_fee":{"start_fee":1300,"mileage":3000,"duration":600}}}',
    '{"total_price":2600,"start_price":1000,"mileage_price":1200,"duration_price":400,"bill_distance":18000,"bill_duration":1620}',
    '{"sp_rate":0.80,"platform_rate":0.20}',
    '{"price_rule_version":{"1":"v2026.11"},"station_type":0,"pick_up_distance":800,"delivery_distance":18000,"cross_city_type":0,"commission_free_order":0}'
),
(
    'WO202611060003', 0, 0, 0,
    '{"total_price":7600,"start_price":1500,"mileage_price":4200,"duration_price":1200,"long_price":700,"bill_distance":52000,"bill_duration":3600}',
    '{}',
    '{}',
    '{"sp_rate":0.75,"platform_rate":0.25}',
    '{"price_rule_version":{"1":"v2026.11"},"station_type":1,"pick_up_distance":1500,"delivery_distance":52000,"cross_city_type":1}'
),
(
    'WO202611060004', 0, 0, 0,
    '{"total_price":2400,"start_price":1200,"mileage_price":900,"duration_price":300,"bill_distance":12000,"bill_duration":1200}',
    '{}',
    '{}',
    '{"sp_rate":0.80,"platform_rate":0.20}',
    '{"price_rule_version":{"1":"v2026.11"},"station_type":2,"pick_up_distance":900,"delivery_distance":12000,"cross_city_type":1}'
),
(
    'WO202611060005', 0, 0, 0,
    '{"total_price":4100,"start_price":1300,"mileage_price":2000,"duration_price":800,"bill_distance":22000,"bill_duration":2400}',
    '{}',
    '{}',
    '{"sp_rate":0.80,"platform_rate":0.20}',
    '{"price_rule_version":{"1":"v2026.11"},"station_type":2,"pick_up_distance":1100,"delivery_distance":22000,"cross_city_type":0}'
);

-- 为批量补充订单生成扩展数据（可重复执行）
INSERT IGNORE INTO `order_ext` (
    `order_no`, `bill_distance`, `bill_duration`, `end_after_order`,
    `estimate_price_info`, `price_info`, `driver_price_info`, `split_info`, `extra_info`
)
SELECT
    o.order_no,
    8000 + (o.id % 20) * 1000,
    600 + (o.id % 20) * 60,
    CASE WHEN o.order_status IN (50, 39) THEN 1 ELSE 0 END,
    JSON_OBJECT(
        'total_price', o.total_price,
        'start_price', 1300,
        'mileage_price', 2000,
        'duration_price', 800,
        'bill_distance', 8000 + (o.id % 20) * 1000,
        'bill_duration', 600 + (o.id % 20) * 60
    ),
    CASE WHEN o.order_status IN (50, 39, 10) THEN JSON_OBJECT(
        'total_price', o.total_price,
        'start_price', 1300,
        'mileage_price', 2000,
        'duration_price', 800,
        'bill_distance', 8000 + (o.id % 20) * 1000,
        'bill_duration', 600 + (o.id % 20) * 60
    ) ELSE JSON_OBJECT() END,
    CASE WHEN o.order_status IN (50, 39, 10) THEN JSON_OBJECT(
        'total_price', o.driver_total_price,
        'start_price', 1000,
        'mileage_price', 1600,
        'duration_price', 600
    ) ELSE JSON_OBJECT() END,
    JSON_OBJECT('sp_rate', 0.80, 'platform_rate', 0.20),
    JSON_OBJECT(
        'price_rule_version', JSON_OBJECT('1', 'v2026.11'),
        'station_type', o.id % 3,
        'pick_up_distance', 800 + (o.id % 10) * 100,
        'delivery_distance', 8000 + (o.id % 20) * 1000,
        'cross_city_type', 0
    )
FROM `order_info` o
WHERE o.order_no LIKE 'WO20261107%'
  AND o.status = 1
  AND NOT EXISTS (SELECT 1 FROM `order_ext` e WHERE e.order_no = o.order_no);
