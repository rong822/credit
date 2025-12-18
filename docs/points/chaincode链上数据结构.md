# 积分相关，chaincode链上数据结构


## 积分账户数据结构 Account

字段名称 | 字段类型 | 说明 | 备注
- | - | - | -
accountName | string | 积分账户的名称 | Key，由公钥SHA256得到
balances | []balance | 积分账户中各项积分的余额 | 
txCount | uint | 此账户转出积分的次数累计 | 用于验证nonce

### 单项积分余额 balance

字段名称 | 字段类型 | 说明 | 备注
- | - | - | -
pointsName | string | 积分名称 |
pointsBalance | uint32 | 积分余额 | 

### 积分种类 Points

字段名称 | 字段类型 | 说明 | 备注
- | - | - | -
name | string | 积分名称 | Key
circulation | uint32 | 积分目前发放总量 | 不能超过最大限额
ceiling | uint64 | 积分最大发放限额 | 最大限额，默认值为60,000,000，年增长5%，起点为2018年
issuer | string | 发放此积分的账户 | 为签发者公钥的SHA256 Hash
issueCount | string | 此积分发放的次数累计 | 用于验证nonce
