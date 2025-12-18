# 积分相关，chaincode接口（logic）调用参数的数据结构


## 积分发放时参数 paramIssuePoints

字段名称 | 字段类型 | 说明 | java字段名称 | 备注
- | - | - | - | -
issueProposal | issueProposal | 积分发放参数 |  | 包含了发放积分的对象，积分种类，数额和验证信息
pubKey | string | 验证公钥 |  | 此公钥需和链上积分Points中的[issuer](chaincode链上数据结构.md#积分种类-Points)相同
signature | string | 验证签名 |  | 

### issueProposal
字段名称 | 字段类型 | 说明 | java字段名称 | 备注
- | - | - | - | -
accountName | string | 账户名称 |  | Key，若链上不存在则新建积分账户
pointsName | string | 积分名称 |  | 需在链上已经存在
amount | uint32 | 积分数额 |  | 
nonce | uint | 待验证的nonce |  | 需为Points中的[issueCount](chaincode链上数据结构.md#积分种类-Points)+1


## 积分查询时参数 paramQueryBalance

字段名称 | 字段类型 | 说明 | java字段名称 | 备注
- | - | - | - | -
accountName | string | 积分查询对象 | 可为空，为空时查询此积分类型的详情 
pointsName | string | 查询积分名称 | 可为空，为空时查询账户内的所有的积分余额，两个字段不可同时为空


## 积分转让时参数 paramTransferPoints

字段名称 | 字段类型 | 说明 | java字段名称 | 备注
- | - | - | - | -
transferProposal | transferProposal | 积分转让参数 |  | 包含转入账户，转出账户，积分种类，数额和校验nonce值
pubKey | string | 验证公钥 |  | 
signature | string | 验证签名 |  | 

### transferProposal
字段名称 | 字段类型 | 说明 | java字段名称 | 备注
- | - | - | - | -
payerName | string | 积分转出账户 |  | 必须链上存在
payeeName | string | 积分转入账户 |  | 若链上不存在，则会自动创建账户
pointsName | string | 转让积分名称 |  | 必须链上存在
amount | uint32 | 积分转让数额 |  | 不能为0，且要小于转出方余额
nonce | uint | 待验证的nonce |  | 需为转出账户Account中的[TxCount](chaincode链上数据结构.md#积分账户数据结构-Account)+1


## 积分历史查询时参数 paramQueryHistory

字段名称 | 字段类型 | 说明 | java字段名称 | 备注
- | - | - | - | -
accoutName | string | 积分历史查询对象 |  | 对应Account中的accountName



## 注册新积分种类 points
与链上[Points](chaincode链上数据结构.md#积分种类-Points)结构相同

字段名称 | 字段类型 | 说明 | java字段名称 | 备注
- | - | - | - | -
name | string | 积分名称 |  | Key
circulation | uint32 | 积分目前发放总量 |  | 不能超过最大限额，默认为0
ceiling | uint64 | 积分最大发放限额 |  | 最大限额，默认值为60,000,000，年增长5%，起点为2018年
issuer | string | 发放此积分的账户 |  | 为签发者公钥的SHA256哈希
issueCount | string | 此积分发放的次数累计 |  | 用于验证nonce，默认为0

