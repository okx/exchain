# deal

| 字段         | 数据类型    | 索引        | 说明           |
|--------------|-------------|-------------|----------------|
| timestamp    | int64       | index       | 成交记录时间戳 |
| block_height | int64       | PRIMARY_KEY | 区块高度       |
| order_id     | varchar(30) | PRIMARY_KEY | 订单id         |
| sender       | varchar(80) | index       | 订单所有者地址 |
| product      | varchar(20) | index       | 币对名称       |
| side         | varchar(10) |             | 买/卖          |
| price        | double      |             | 成交价格       |
| quantity     | double      |             | 成交数量       |

# fee_detail
| 字段      | 数据类型    | 索引  | 说明               |
|-----------|-------------|-------|--------------------|
| address   | varchar(80) | index | 用户地址           |
| fee       | varchar(40) |       | 手续费金额         |
| fee_type  | varchar(20) |       | 手续费类型         |
| timestamp | int64       | index | 手续费收取的时间戳 |

# order

| 字段             | 数据类型    | 索引        | 说明                                                                                             |
|------------------|-------------|-------------|--------------------------------------------------------------------------------------------------|
| tx_hash          | int64       |             | 订单对应的tx哈希                                                                                 |
| order_id         | varchar(30) | PRIMARY_KEY | 订单id                                                                                           |
| sender           | varchar(80) | index       | 订单所有者地址                                                                                   |
| product          | varchar(20) | index       | 币对名称                                                                                         |
| side             | varchar(10) |             | 买卖方向                                                                                         |
| price            | varchar(40) |             | 挂单价格                                                                                         |
| quantity         | varchar(40) |             | 挂单数量                                                                                         |
| status           | int64       | index       | 订单状态，(0-5):(Open, Filled, Cancelled, Expired, PartialFilledCancelled, PartialFilledExpired) |
| filled_avg_price | varchar(40) |             | 成交均价                                                                                         |
| remain_quantity  | varchar(40) |             | 剩余数量                                                                                         |
| timestamp        | int64       | index       | 订单创建时间戳                                                                                   |

# transaction
| 字段      | 数据类型    | 索引  | 说明                                            |
|-----------|-------------|-------|-------------------------------------------------|
| tx_hash   | int64       |       | 交易哈希                                        |
| type      | int64       | index | 交易类型，1:Transfer, 2:NewOrder, 3:CancelOrder |
| address   | varchar(80) | index | 交易所属用户地址                                |
| symbol    | varchar(20) |       | 币名称或币对名称                                |
| side      | varchar(10) |       | 买卖或转账方向：1:buy, 2:sell, 3:from, 4:to     |
| quantity  | varchar(40) |       | 转账数量或挂单数量                              |
| fee       | varchar(40) |       | 交易手续费                                      |
| timestamp | int64       | index | 交易发生的时间戳                                |

# Candles
1. kline 以其更新频度的差异，分别保存在不同的表下，表结构都相同
2. kline_m1中产生一条新的记录的频度（后称时间片）是一分钟（即60秒）
3. 如果一个时间片内没有成交记录（deal），那么将不产生k线保存至数据库中
4. kline_m1, kline_m3, kline_m5 预计将保留最近一个月的数据；其它 K 线数据永久保留。
5. kline_m1一个月数据量预估
    * 单记录上限字节数（20+8+8+8+8+8+8)= 68
    * 索引大小约为原始数据的 30% ~ 50%
    * 单币对一个月记录数：44640
    * 预计支持 200 币对
    * 预计总空间消耗：780MB ~ 900MB
6. 以下对 kline 表空间预估以支持 200 币对且所有币对在每分钟都存在成交为前提

## kline_m1

| 字段      | 数据类型    | 索引        | 说明                                                                |
|-----------|-------------|-------------|---------------------------------------------------------------------|
| product   | varchar(20) | primary key | 产品（即币对，如xxb_okb）                                           |
| timestamp | int64       | primary key | Unix时间戳, the number of seconds elapsed since January 1, 1970 UTC |
| open      | double      |             | 本时间片内的起始价格                                                |
| close     | double      |             | 本时间片内的结束价格                                                |
| high      | double      |             | 本时间片内的最高价格                                                |
| low       | double      |             | 本时间片内的最低价格                                                |
| volume    | double      |             | 本时间片内的成交量                                                  |

## kline表与说明
| table_name   | 更新频度 | 保留天数(预估) | 单日新增数据量(预估) | 累计数据量（预估） |
|--------------|----------|----------------|----------------------|--------------------|
| kline_m1     | 60       | 30             | 25MB~29MB            | 780MB~900MB        |
| kline_m3     | 180      | 30             | 8MB~10MB             | 240MB~300MB        |
| kline_m5     | 300      | 30             | 5MB~6MB              | 150MB~180MB        |
| kline_m15    | 900      |                | 1MB~2MB              |                    |
| kline_m30    | 1800     |                | < 1MB                |                    |
| kline_m60    | 3600     |                | < 1MB                |                    |
| kline_m120   | 7200     |                | < 1MB                |                    |
| kline_m240   | 14400    |                | < 1MB                |                    |
| kline_m360   | 21600    |                | < 1MB                |                    |
| kline_m720   | 43200    |                | < 1MB                |                    |
| kline_m1440  | 86400    |                | < 1MB                |                    |
| kline_m10080 | 604800   |                | < 1MB                |                    |

# Restful API & Persistent Map
| url                            | method | table              |
|--------------------------------|--------|--------------------|
| /transactions                  | GET    | transactions       |
| /block_tx_hashes/{blockHeight} | GET    |                    |
| /order/list/{openOrClosed}     | GET    | orders             |
| /deals                         | GET    | deals              |
| /tickers/{instrumentId}        | GET    | 未持久化，内存数据 |
| /candles/{instrumentId}        | GET    | kline_m*           |