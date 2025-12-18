# 积分相关，chaincode接口(logic)文档


## 积分发放

接口名称：issuePoints

实现功能：issuer向积分账户中发放积分

传入参数：[paramIssuePoints](chaincode接口参数数据结构.md#积分发放时参数-paramissuepoints)

返回结果：若成功，以 `json` 格式返回积分发放后账户的积分明细


## 积分查看

接口名称： queryBalance

实现功能：查询账户的积分，查询某种积分的信息，查询账户下的全部积分

传入参数：[paramQueryBalance](chaincode接口参数数据结构.md#积分查询时参数-paramquerybalance)

返回结果：若成功，以 `json` 格式返回积分的查询结果


## 积分转让

接口名称： transferPoints

实现功能：账户间转让积分

传入参数：[paramTransferPoints](chaincode接口参数数据结构.md#积分转让时参数-paramtransferpoints)

返回结果：若成功，以 `json` 格式返回积分转让后的payer和payee积分明细


## 查询积分用户历史

接口名称： queryAccountHistory

实现功能：查询积分账户的历史记录

传入参数：[queryAccountHistory](chaincode接口参数数据结构.md#积分历史查询时参数-paramqueryhistory)

返回结果：若成功，以 `json` 格式返回积分历史查询结果


## 注册新积分种类

接口名称： registerPoints

实现功能：注册新积分种类

传入参数：[points](chaincode接口参数数据结构.md#注册新积分种类-points)

返回结果：若成功，以 `json` 格式返回创建的Points信息
