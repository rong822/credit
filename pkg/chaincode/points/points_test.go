package points

import (
	"encoding/json"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

//使用此单元测试前，需要把chaincode中的验证去掉

func TestMain(m *testing.M) {
	m.Run()
}

func tIssuePoints() {
	stdData := `{
  "issueProposal": {
    "accountName": "c1",
    "pointsName": "BOCpoints",
    "amount": 100,
    "nonce": 2
  },
"pubKey": "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCvBGcCN6LRryZSKWbziij5X6Dd\nUeMxVSLXp+6QoG7RqloWWV/dVjP0R3h+f7hy2Pd7qtqjdKYCygaY/S2+iCfpsrEz\nOhksH7NySf8R0b9Z9QHzTWb9gQ+Gixv5P3if1E2mENcaG/uBizkWCQ2o+3gy+gK/\nq5DIuLQycb9rDVheFQIDAQAB\n-----END PUBLIC KEY-----\n",
  "signature": "7a938b65db6cf5ee169d11ecd55e6c2a886cb7e309d5eb8b77c078d6f34224baef56eb9c61f0502ce28e416fce1b2e33a17d79f441a6d13e835514996cbc4accaca6468a0f0bbdb71f55b62cc5c0d28233954e62cb519deee8aadde612bb48f8addc4b8f0e4e98c99cef651a049c1d1ebea8c40199a4141e6acc41ecb712a377"
}`
	mockInvoke(MethodIssuePoints, stdData)
}

func TestIssuePoints(t *testing.T) {
	tInit()
	tIssuePoints()

	expectedAccount := new(Account)
	json.Unmarshal([]byte(`{
      "accountName": "a1",
      "balances": [
        {
          "pointsName": "ICBCpoints",
          "pointsBalance": 200
        }
      ]
}`), expectedAccount)
	actualAccount := new(Account)
	json.Unmarshal(stub.State[getKey("a1")], actualAccount)
	assert.Equal(t, expectedAccount, actualAccount, "积分发放：账户已存在此类积分")

	expectedPoints := new(Points)
	json.Unmarshal([]byte(`{
  "name": "ICBCpoints",
  "issuer": "ICBC",
  "ceiling": 60000000,
  "circulation": 100
}`), expectedPoints)
	actualPoints := new(Points)
	json.Unmarshal(stub.State[getKey("ICBCpoints")], actualPoints)
	assert.Equal(t, expectedPoints, actualPoints, "积分发放：账户已存在此类积分,发放总量更新")
}

func tIssuePoints2() {
	stdData := `{
  "issueProposal": {
    "accountName": "a0",
    "pointsName": "ICBCpoints",
    "amount": 100,
    "nonce": 1
  }
}`
	mockInvoke(MethodIssuePoints, stdData)
}

func TestIssuePoints2(t *testing.T) {
	tInit()
	tIssuePoints2()

	expectedAccount := new(Account)
	json.Unmarshal([]byte(`{
      "accountName": "a0",
      "balances": [
        {
          "pointsName": "ICBCpoints",
          "pointsBalance": 100
        }
      ]
}`), expectedAccount)
	actualAccount := new(Account)
	json.Unmarshal(stub.State[getKey("a0")], actualAccount)
	assert.Equal(t, expectedAccount, actualAccount, "积分发放：无账户，创建新的账户")

	expectedPoints := new(Points)
	json.Unmarshal([]byte(`{
  "name": "ICBCpoints",
  "issuer": "ICBC",
  "ceiling": 60000000,
  "circulation": 100
}`), expectedPoints)
	actualPoints := new(Points)
	json.Unmarshal(stub.State[getKey("ICBCpoints")], actualPoints)
	assert.Equal(t, expectedPoints, actualPoints, "积分发放：已有账户，创建新的积分种类,发放总量更新")
}

func tQueryBalance() peer.Response {
	stdData := `{
  "accountName": "a1",
  "pointsName": "ICBCpoints"
}`
	return mockInvoke(MethodQueryBalance, stdData)
}

func TestQueryBalance(t *testing.T) {
	tInit()
	res := tQueryBalance()

	expectedBalance := new(balance)
	json.Unmarshal([]byte(`{
          "pointsName": "ICBCpoints",
          "pointsBalance": 100
}`), expectedBalance)
	actualBalance := new(balance)
	json.Unmarshal(res.Payload, actualBalance)
	assert.Equal(t, expectedBalance, actualBalance, "积分查询")
}

func tQueryBalance2() peer.Response {
	stdData := `{
"accountName": "a8",
  "pointsName": "ICBCpoints"
}`
	return mockInvoke(MethodQueryBalance, stdData)
}

func TestQueryBalance2(t *testing.T) {
	tInit()
	res := tQueryBalance2()
log.Println(string(res.Payload))
	expectedBalance := new(balance)
	json.Unmarshal([]byte(`{
          "pointsName": "ABCpoints",
          "pointsBalance": 0
}`), expectedBalance)
	actualBalance := new(balance)
	json.Unmarshal(res.Payload, actualBalance)
	assert.Equal(t, expectedBalance, actualBalance, "积分查询:账户没有相关积分，返回0")
}

func tQueryBalance3() peer.Response {
	stdData := `{
    "accountName": "a1"
  }`
	return mockInvoke(MethodQueryBalance, stdData)
}

func TestQueryBalance3(t *testing.T) {
	tInit()
	res := tQueryBalance3()

	expectedBalance := new(balance)
	json.Unmarshal([]byte(`{
  "accountName": "a1",
  "balances": [
    {
      "pointsName": "ICBCpoints",
      "pointsBalance": 100
    }
  ]
}`), expectedBalance)
	actualBalance := new(balance)
	json.Unmarshal(res.Payload, actualBalance)
	assert.Equal(t, expectedBalance, actualBalance, "积分查询:没有指定积分种类，查询账户全部积分")
}

func tTransferPoints() {
	stdData := `{
  "signature": "472c6b9fd70f0761161ff3a6b3e8dfd0ae146afdaf2582e5ec555b5a6df710cbdcaedbbfeea0afabfb321efd2d047a2d3560f492738c1221a92945a1b580b11d2f15d23c2c7ad98eb66eed34a2359a2b57e454147d74bcacef9e4c02b61ccb18e64aa061a92bb79f33aa0d194df1cd334fff7a8cb1a83ed39006c93db9d22f81",
  "transferProposal": {
    "payerName": "7bb5042b063818959f28afa01062b6c6dcf242442babff6bbea17162a9991841",
    "payeeName": "73690bf1eb790f8cad3771b444de97f69c16502a4a3ea3ede8acc4f2d3af5fc2",
    "pointsName": "TokenF",
    "amount": 10,
    "nonce": 1
  },
  "pubKey": "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCA+qObrAG8slUTRnPONX5SclCK2BDVZ2xA7+vG\n8qz/TdRsD8XQaiuK2xiSILmO4XWqUUo0VJLpcWyFXjNfgP6OL5ntEqNIlwseFJLfWAdlifu96+ZV\n7s/qnZAw2sAhONUWCHk5TBerZ3sOxqD1NB20iiXkEmho8C4Qa62p1rNtpwIDAQAB\n-----END PUBLIC KEY-----\n"
}`
	mockInvoke(MethodTransferPoints, stdData)
}

func TestTransferPoints(t *testing.T) {
	tInit()
	//tIssuePoints2()
	tTransferPoints()

	expectedPayer := new(Account)
	json.Unmarshal([]byte(`{
     "accountName": "7bb5042b063818959f28afa01062b6c6dcf242442babff6bbea17162a9991841",
     "balances": [
       {
         "pointsName": "TokenF",
         "pointsBalance": 90
       }
     ],
"txCount": 1
}`), expectedPayer)
	actualPayer := new(Account)
	json.Unmarshal(stub.State[getKey("7bb5042b063818959f28afa01062b6c6dcf242442babff6bbea17162a9991841")], actualPayer)
	assert.Equal(t, expectedPayer, actualPayer, "积分转让：对比支付方账户余额")

	expectedPayee := new(Account)
	json.Unmarshal([]byte(`{
     "accountName": "73690bf1eb790f8cad3771b444de97f69c16502a4a3ea3ede8acc4f2d3af5fc2",
     "balances": [
       {
         "pointsName": "TokenF",
         "pointsBalance": 10
       }
     ]
}`), expectedPayee)
	actualPayee := new(Account)
	json.Unmarshal(stub.State[getKey("73690bf1eb790f8cad3771b444de97f69c16502a4a3ea3ede8acc4f2d3af5fc2")], actualPayee)
	assert.Equal(t, expectedPayee, actualPayee, "积分转让：对比接收方账户余额")
}
