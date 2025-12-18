# Credit Server API DOC  
Credit Server api document  
  
# Info  
  
|  |        |  
|------------|---------------|  
| host       | 47.92.239.173 |  
| api-server | 8787          |  
  
  
# API  
  此处所给示例为测试用例，实际使用中的accountName应为Hash值
  
| Method | Route                           | Describe         | Body                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         | Header                                                                                                                                                                                                                                      | Response                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |  
|--------|---------------------------------|------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|  
| Post   | /api/credit/issuePoints         | 发放积分 | {   "issueProposal": {     "accountName": "A1",     "pointsName": "SBCpoints",     "amount": 100,     "nonce": 1   },   "pubKey": "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSX3l6zvgrvYFgyu4NmMjff26M\nFlhWpkdsJP+qPmlceMujbhhI8ArEayWBJGUjVL3pYevgMe2fexlynvF2a93ZX6Iu\nTzDhk/HpwqASEPD0abPsZH11uN8ApDRtnjhpmK9dz1k2hGAaGq6+Ep7mmMukhszm\n/nEjBedIK1Q5dCjhuQIDAQAB\n-----END PUBLIC KEY-----",   "signature": "3787698d16c729134b84165f0103b1332dccbbbe13aaa3fd86c028b15940cf685d991850ea4f03a9910f331acdff7c5548a65665235187e81656a133abac0e32fd506f67caed7ec0f7eb108294c52f19f5a4f8e6d66e2cda7cc040f7267e7680981da9d9bd71a511ce6d7cd9681ba7b644ef5f9f92fe61a0751ccbdee65a70e8" }                | Content-Type: application/json     Bearer: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDA5MDE5MDEsInVzZXJuYW1lIjoiYXNkZmFzZGZzZGFmYXMiLCJwYXNzd29yZCI6IjIxMzEyMyIsImlhdCI6MTU0MDg2NTkwMX0.52ARsD1cRjdHA8-EF1nS_PzHdZskF7b939ard6D_yQo | {   "code": "0",   "results": {     "status": "SUCCESS",     "info": "",     "txid": "4345a5ceba9b40dfd72dccec60e33e4dd32e6f8153c29b3999e34cd58e513477",     "payload": "{\"accountName\":\"A1\",\"balances\":[{\"pointsName\":\"SBCpoints\",\"pointsBalance\":100}]}"   },   "msg": "issuePoints succeed " }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |  
| Post   | /api/credit/queryBalance        | 查询积分详情 | {   "accountName": "A1",   "pointsName": "SBCpoints" }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       | Content-Type: application/json     Bearer: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDA5MDE5MDEsInVzZXJuYW1lIjoiYXNkZmFzZGZzZGFmYXMiLCJwYXNzd29yZCI6IjIxMzEyMyIsImlhdCI6MTU0MDg2NTkwMX0.52ARsD1cRjdHA8-EF1nS_PzHdZskF7b939ard6D_yQo | {   "code": "0",   "results": {     "status": "SUCCESS",     "info": "",     "txid": "d4d69c93fbc35f366e49c508334c95d0bda9c424f0e20da360d62982c285500d",     "payload": "{\"pointsName\":\"SBCpoints\",\"pointsBalance\":100}"   },   "msg": "queryBalance succeed " }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |  
| Post   | /api/credit/transferPoints      | 积分转让 | {   "transferProposal": {     "payerName": "A1",     "payeeName": "a1",     "pointsName": "SBCpoints",     "amount": 10,     "nonce": 1   },   "pubKey": "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSX3l6zvgrvYFgyu4NmMjff26M\nFlhWpkdsJP+qPmlceMujbhhI8ArEayWBJGUjVL3pYevgMe2fexlynvF2a93ZX6Iu\nTzDhk/HpwqASEPD0abPsZH11uN8ApDRtnjhpmK9dz1k2hGAaGq6+Ep7mmMukhszm\n/nEjBedIK1Q5dCjhuQIDAQAB\n-----END PUBLIC KEY-----",   "signature": "75049fe0288cba8d35a7898dbe70e591db9e1802bfbd8e909d462e18f350ab42aabff0ac41c629477a7f1e46e188fb13c620471573c4b8e68dde2d8afc812b41e755c841db00cfc1b591f298535ec4a7c3a0a78caf95283f846528a25e2751a179ea11164b1e580b63e45d70666d62adf213fd1dc77976f56b1f8809d58fd2fb" } | Content-Type: application/json     Bearer: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDA5MDE5MDEsInVzZXJuYW1lIjoiYXNkZmFzZGZzZGFmYXMiLCJwYXNzd29yZCI6IjIxMzEyMyIsImlhdCI6MTU0MDg2NTkwMX0.52ARsD1cRjdHA8-EF1nS_PzHdZskF7b939ard6D_yQo | {   "code": "0",   "results": {     "status": "SUCCESS",     "info": "",     "txid": "a8cfe222032a739351c07fdc29073dfa65ca4a4af83666c3b7cb9841a21b64b8",     "payload": "[Payer: {\"accountName\":\"A1\",\"balances\":[{\"pointsName\":\"SBCpoints\",\"pointsBalance\":90}]}, Payee: {\"accountName\":\"a1\",\"balances\":[{\"pointsBalance\":3190},{\"pointsName\":\"ICBCpoints\",\"pointsBalance\":100},{\"pointsName\":\"ABCpoints\",\"pointsBalance\":100},{\"pointsName\":\"BOCpoints\",\"pointsBalance\":400},{\"pointsName\":\"SBCpoints\",\"pointsBalance\":10}],\"txCount\":1}]"   },   "msg": "transferPoints succeed " }                                                                                                                                                                        |  
| Post   | /api/credit/queryAccountHistory | 查询积分账户记录 | {   "accoutName": "A1" }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | Content-Type: application/json     Bearer: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDA5MDE5MDEsInVzZXJuYW1lIjoiYXNkZmFzZGZzZGFmYXMiLCJwYXNzd29yZCI6IjIxMzEyMyIsImlhdCI6MTU0MDg2NTkwMX0.52ARsD1cRjdHA8-EF1nS_PzHdZskF7b939ard6D_yQo | {   "code": "0",   "results": {     "status": "SUCCESS",     "info": "",     "txid": "ec223df933fcbcbb12580833de2c5e9b7b384a14c64f25e1e43c6562fdace804",     "payload": "[{\"tx_id\":\"4345a5ceba9b40dfd72dccec60e33e4dd32e6f8153c29b3999e34cd58e513477\",\"value\":\"{\\\"accountName\\\":\\\"A1\\\",\\\"balances\\\":[{\\\"pointsName\\\":\\\"SBCpoints\\\",\\\"pointsBalance\\\":100}]}\",\"timestamp\":{\"seconds\":1540867039,\"nanos\":723000000}},{\"tx_id\":\"a8cfe222032a739351c07fdc29073dfa65ca4a4af83666c3b7cb9841a21b64b8\",\"value\":\"{\\\"accountName\\\":\\\"A1\\\",\\\"balances\\\":[{\\\"pointsName\\\":\\\"SBCpoints\\\",\\\"pointsBalance\\\":100}],\\\"txCount\\\":1}\",\"timestamp\":{\"seconds\":1540868149,\"nanos\":324000000}}]"   },   "msg": "queryAccountHistory succeed " } |  
| Post   | /api/credit/registerPoints      | 注册新积分种类 | {   "name": "SBCpoints",   "ceiling": 100000000,   "issuer": "ce20de88c17dd6a6f046b9434e78dc45954d4a2103bf96a9fc46189614889815",   "issueCount": 0 }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | Content-Type: application/json     Bearer: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDA5MDE5MDEsInVzZXJuYW1lIjoiYXNkZmFzZGZzZGFmYXMiLCJwYXNzd29yZCI6IjIxMzEyMyIsImlhdCI6MTU0MDg2NTkwMX0.52ARsD1cRjdHA8-EF1nS_PzHdZskF7b939ard6D_yQo | {   "code": "0",   "results": {     "status": "SUCCESS",     "info": "",     "txid": "0d8ba1bad485d41d5aacc636e8e8ec8bc7637dcd4c14f98632088e9a94cafe7c",     "payload": "{\"name\":\"SBCpoints\",\"circulation\":0,\"issuer\":\"ce20de88c17dd6a6f046b9434e78dc45954d4a2103bf96a9fc46189614889815\",\"issueCount\":0}"   },   "msg": "registerPoints succeed " }                                                                                                                                                                                                                                                                                                                                                                                                                                            |Post |  |  |  |  |  |  
  
# 字段说明  
  
## issuePoints  
  
字段名称 | 字段类型 | 说明 | 备注  
- | - | - | -   
issueProposal | issueProposal | 积分发放参数 | 包含了发放积分的对象，积分种类，数额和验证信息  
pubKey | string | 验证公钥 |  ECDSA公钥，下同
signature | string | 验证签名 | 16进制转码后的结果  
  
### issueProposal  
字段名称 | 字段类型 | 说明 | 备注  
- | - | - | -   
accountName | string | 账户名称 | Key，若链上不存在则新建积分账户, accountName是其公钥的SHA256,表示为16进制字符串，下同
pointsName | string | 积分名称 | 需在链上已经存在  
amount | uint32 | 积分数额 |   
nonce | uint | 待验证的nonce | 需为Points中的issueCount+1  
  
  
## queryBalance  
  
字段名称 | 字段类型 | 说明 | 备注  
- | - | - | -   
accountName | string | 积分查询对象 | 可为空，为空时查询链上积分类型的详情 
pointsName | string | 查询积分名称 | 可为空，为空时查询账户内的所有的积分余额，两个字段不可同时为空
  
  
## transferPoints  
  
字段名称 | 字段类型 | 说明 | 备注  
- | - | - | -   
transferProposal | transferProposal | 积分转让参数 |  | 包含转入账户，转出账户，积分种类，数额和校验nonce值  
pubKey | string | 验证公钥 |  |   
signature | string | 验证签名 |  |   
  
### transferProposal  
字段名称 | 字段类型 | 说明 | 备注  
- | - | - | -   
payerName | string | 积分转出账户 | 必须链上存在,payerName是其公钥的SHA256,表示为16进制字符串  
payeeName | string | 积分转入账户 | 若不存在，自动创建账户  
pointsName | string | 转让积分名称 | 必须链上存在  
amount | uint32 | 积分转让数额 | 不能为0，且要小于转出方余额  
nonce | uint | 待验证的nonce | 需为转出账户Account中的TxCount+1  
  
  
## queryAccountHistory  
  
字段名称 | 字段类型 | 说明 | 备注  
- | - | - | -   
accoutName | string | 积分历史查询对象 | 对应Account中的accountName  

  
  
## registerPoints  
与链上Points结构相同  
  
字段名称 | 字段类型 | 说明 | 备注  
- | - | - | -   
name | string | 积分名称 | Key  
ceiling | uint64 | 积分最大发放限额 | 最大限额，默认值为60,000,000，年增长5%，起点为2018年  
issuer | string | 发放此积分的账户 | 为签发者公钥的SHA256哈希  
issueCount | string | 此积分发放的次数累计 | 用于验证nonce，默认为0