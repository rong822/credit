package points

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var stub = shim.NewMockStub(ChaincodeName, new(Chaincode))

var accountID = getKey("a1")

func mockInvoke(fn string, body string) peer.Response {
	return stub.MockInvoke("1", [][]byte{[]byte(fn), []byte(body)})
}

func getKey(key string) string {
	//sha := md5.Sum([]byte(key))
	//return hex.EncodeToString(sha[:])
	return key
}

func tInit() {
	stub.State = make(map[string][]byte)

	stub.MockTransactionStart("1")
	stub.PutState("7bb5042b063818959f28afa01062b6c6dcf242442babff6bbea17162a9991841", []byte(`{
  "accountName": "7bb5042b063818959f28afa01062b6c6dcf242442babff6bbea17162a9991841",
  "balances": [
    {
      "pointsName": "TokenF",
      "pointsBalance": 100
    }
  ],
  "txCount": 0
}`))
	stub.MockTransactionEnd("1")

	stub.MockTransactionStart("1")
	stub.PutState(getKey("a2"), []byte(`{
  "accountName": "a2",
  "balances": [
    {
      "pointsName": "TokenF",
      "pointsBalance": 100
    }
  ],
  "txCount": 0
}`))
	stub.MockTransactionEnd("1")

	stub.MockTransactionStart("1")
	stub.PutState(getKey("ICBCpoints"), []byte(`{
  "name": "TokenF",
  "circulation": 0,
"issuer": "ce20de88c17dd6a6f046b9434e78dc45954d4a2103bf96a9fc46189614889815",
"issueCount": 0
}`))
	stub.MockTransactionEnd("1")
}
