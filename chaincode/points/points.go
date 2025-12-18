package points

import (
	"bytes"
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gitlab.bigtree.com/dashu-blockchain/credit/pkg/errors"
)
import ts "github.com/golang/protobuf/ptypes/timestamp"

type paramIssuePoints struct {
	IssueProposal *issueProposal `json:"issueProposal,omitempty" valid:"required"`
	PubKey        string         `json:"pubKey,omitempty" valid:"required"`
	Signature     string         `json:"signature,omitempty" valid:"required"`
}

type issueProposal struct {
	AccountName string `json:"accountName,omitempty" valid:"required"`
	PointsName  string `json:"pointsName,omitempty" valid:"required"`
	Amount      uint32 `json:"amount,omitempty" valid:"required"`
	Nonce       uint   `json:"nonce,omitempty" valid:"required"`
}

func newParamIssuePoints() *paramIssuePoints {
	return &paramIssuePoints{
		IssueProposal: new(issueProposal),
	}
}

func issuePoints(stub shim.ChaincodeStubInterface, body []byte) ([]byte, error) {
	param := newParamIssuePoints()
	if err := json.Unmarshal(body, param); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, err := govalidator.ValidateStruct(param); !ok {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}

	//验证签名
	if _, err := validateSignature(param.IssueProposal, param.PubKey, param.Signature); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	acc := newAccount()
	acc.AccountName = param.IssueProposal.AccountName
	acc.Balances = append(acc.Balances, &balance{param.IssueProposal.PointsName, param.IssueProposal.Amount})
	//验证公钥是否是issuer, 是否存在积分种类，验证nonce，验证是否发放量超限，更新限额
	for _, b := range acc.Balances {
		_, err := multiStepValidate(stub, b, param.IssueProposal.Nonce, param.PubKey)
		if err != nil {
			return nil, errors.Error(errors.StatusBadRequest, err)
		}
	}
	//若账户不存在, 进行逐步验证后,创建新账户
	if !validateAccount(stub, acc.GetKey()) {
		accBytes, err := createAccount(stub, acc)
		if err != nil {
			return nil, errors.Error(errors.StatusBadRequest, "Failed to create account.")
		}
		return accBytes, nil
	}

	//若账户存在则更新积分余额
	resultBytes, _ := updateAccount(stub, acc, true)
	if resultBytes == nil {
		return nil, errors.Error(errors.StatusBadRequest, "Failed to update the balances.")
	}
	return resultBytes, nil
}

type paramQueryBalance struct {
	AccountName string `json:"accountName,omitempty"`
	PointsName  string `json:"pointsName,omitempty"`
}

func queryBalance(stub shim.ChaincodeStubInterface, body []byte) ([]byte, error) {
	param := new(paramQueryBalance)
	if err := json.Unmarshal(body, param); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, err := govalidator.ValidateStruct(param); !ok {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}

	if param.AccountName == "" {
		p, err := queryPointsByKey(stub, param.PointsName)
		if err != nil {
			return nil, errors.Error(errors.StatusNotFound, err)
		}
		if p == nil {
			return nil, errors.Error(errors.StatusBadRequest, "Points doesn't exit.")
		}
		resultBytes, err := json.Marshal(p)
		if err != nil {
			return nil, errors.Error(errors.StatusBadRequest, err)
		}
		return resultBytes, nil
	}

	queryBalance := new(balance)
	if !validateAccount(stub, param.AccountName) {
		return nil, errors.Error(errors.StatusBadRequest, "Account doesn't exit.")
	}

	if param.PointsName == "" {
		resultAccount, _ := queryAllBalances(stub, param.AccountName)
		resultBytes, err := json.Marshal(resultAccount)
		if err != nil {
			return nil, errors.Error(errors.StatusBadRequest, err)
		}
		return resultBytes, nil
	}
	if p, _ := queryPointsByKey(stub, param.PointsName); p == nil {
		return nil, errors.Error(errors.StatusBadRequest, "Points doesn't exit.")
	}

	if !isPointsInAccount(stub, param.AccountName, param.PointsName) {
		queryBalance = &balance{param.PointsName, 0}
	} else {
		err := *new(error)
		queryBalance, err = queryBalanceByKey(stub, param.AccountName, param.PointsName)
		if err != nil {
			return nil, errors.Error(errors.StatusNotFound, err)
		}
		if queryBalance == nil {
			return nil, errors.Errorf(errors.StatusNotFound, "Cannot find %v in %v.", param.PointsName, param.AccountName)
		}
	}

	resultBytes, err := json.Marshal(queryBalance)
	if err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	return resultBytes, nil
}

type paramTransferPoints struct {
	TransferProposal *transferProposal `json:"transferProposal,omitempty" valid:"required"`
	PubKey           string            `json:"pubKey,omitempty" valid:"required"`
	Signature        string            `json:"signature,omitempty" valid:"required"`
}

type transferProposal struct {
	PayerName  string `json:"payerName,omitempty" valid:"required"`
	PayeeName  string `json:"payeeName,omitempty" valid:"required"`
	PointsName string `json:"pointsName,omitempty" valid:"required"`
	Amount     uint32 `json:"amount,omitempty" valid:"required"`
	Nonce      uint   `json:"nonce,omitempty" valid:"required"`
}

func newParamTransferPoints() *paramTransferPoints {
	return &paramTransferPoints{
		TransferProposal: new(transferProposal),
	}
}

func transferPoints(stub shim.ChaincodeStubInterface, body []byte) ([]byte, error) {
	param := newParamTransferPoints()
	if err := json.Unmarshal(body, param); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, err := govalidator.ValidateStruct(param); !ok {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}

	//验证签名
	if _, err := validateSignature(param.TransferProposal, param.PubKey, param.Signature); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	//验证账户是否存在
	if !validateAccount(stub, param.TransferProposal.PayerName) {
		return nil, errors.Error(errors.StatusBadRequest, "Payer account doesn't exit.")
	}
	//payer账户是否与公钥一一对应
	if makeHash(param.PubKey) != param.TransferProposal.PayerName {
		return nil, errors.Error(errors.StatusBadRequest, "Payer account doesn't match public key.")
	}
	//nonce值验证
	if !validateTransferNonce(stub, param.TransferProposal.PayerName, param.TransferProposal.Nonce) {
		return nil, errors.Error(errors.StatusBadRequest, "The nonce doesn't match.")
	}
	//验证是否有此类积分
	if !isPointsInAccount(stub, param.TransferProposal.PayerName, param.TransferProposal.PointsName) {
		return nil, errors.Error(errors.StatusBadRequest, "Payer doesn't have this points.")
	}
	b, err := queryBalanceByKey(stub, param.TransferProposal.PayerName, param.TransferProposal.PointsName)
	if err != nil {
		return nil, errors.Error(errors.StatusNotFound, "Payer doesn't have this points.")
	}
	if param.TransferProposal.Amount == 0 || b.PointsBalance < param.TransferProposal.Amount {
		return nil, errors.Error(errors.StatusBadRequest, "Transfer amount is zero or payer doesn't have sufficient points.")
	}
	var payeeBytes []byte
	if !validateAccount(stub, param.TransferProposal.PayeeName) {
		//若payee账户不存在，为其创建账户
		payeeAccount := newAccount()
		payeeAccount.AccountName = param.TransferProposal.PayeeName
		payeeAccount.Balances = append(payeeAccount.Balances, &balance{param.TransferProposal.PointsName, param.TransferProposal.Amount})
		payeeBytes, err = createAccount(stub, payeeAccount)
		if err != nil {
			return nil, errors.Error(errors.StatusBadRequest, "Failed to create payee's account.")
		}
	} else {
		//若payee账户存在，则更新相关积分
		updatePayee := newAccount()
		updatePayee.AccountName = param.TransferProposal.PayeeName
		updatePayee.Balances = append(updatePayee.Balances, &balance{param.TransferProposal.PointsName, param.TransferProposal.Amount})
		//传入增量，相加
		payeeBytes, err = updateAccount(stub, updatePayee, true)
		if err != nil {
			return nil, errors.Error(errors.StatusBadRequest, "Failed to update payee's balance.")
		}
	}

	updatePayer := newAccount()
	updatePayer.AccountName = param.TransferProposal.PayerName
	updatePayer.TxCount = 1 //TxCount+1,此处为增量
	updatePayer.Balances = append(updatePayer.Balances, &balance{param.TransferProposal.PointsName, param.TransferProposal.Amount})
	//payer账户已经有了此类积分，则将积分值相减
	payerBytes, err := updateAccount(stub, updatePayer, false)
	if err != nil {
		return nil, errors.Error(errors.StatusBadRequest, "Failed to update payer's balance.")
	}

	var buffer bytes.Buffer //输出的bytes形式是 {"Payer":"","Payee":""}
	buffer.WriteString(`{"Payer": `)
	buffer.Write(payerBytes)
	buffer.WriteString(`,"Payee": `)
	buffer.Write(payeeBytes)
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

type paramQueryHistory struct {
	AccountName string `json:"accountName,omitempty" valid:"required"`
}

type KeyModification struct {
	TxId      string        `json:"tx_id"`
	TxType    string        `json:"txType"`
	TxDetails string        `json:"txDetails"`
	Balances  *Account      `json:"balances"`
	Timestamp *ts.Timestamp `json:"timestamp"`
	IsDelete  bool          `json:"is_delete,omitempty"`
}

func queryAccountHistory(stub shim.ChaincodeStubInterface, body []byte) ([]byte, error) {
	param := new(paramQueryHistory)
	if err := json.Unmarshal(body, param); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, err := govalidator.ValidateStruct(param); !ok {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}

	accountBytes, err := stub.GetState(param.AccountName)
	if err != nil {
		return nil, errors.Error(errors.StatusNotFound, err)
	}
	if accountBytes == nil {
		return nil, errors.Error(errors.StatusNotFound, "Account cannot be found.")
	}
	acc := newAccount()
	if err := json.Unmarshal(accountBytes, acc); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, _ := govalidator.ValidateStruct(acc); !ok {
		return nil, errors.Error(errors.StatusBadRequest, "History query can only be invoked by account.")
	}

	historyIter, err := stub.GetHistoryForKey(param.AccountName)
	if err != nil {
		return nil, errors.Error(errors.StatusNotFound, err)
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	theFirstElement := true

	for historyIter.HasNext() {
		keyMod, err := historyIter.Next()
		if err != nil {
			return nil, errors.Error(errors.StatusNotFound, err)
		}

		keyM := KeyModification{
			Timestamp: &ts.Timestamp{},
		}

		//查询结果中有byte形式的结果，为了使结果更易读，新建结构体转为string输出结果
		keyM.TxId = keyMod.TxId
		keyM.TxType = ""    //API server处填入
		keyM.TxDetails = "" //API server处填入
		acc := newAccount()
		json.Unmarshal(keyMod.Value, acc)
		keyM.Balances = acc
		keyM.Timestamp.Nanos = keyMod.Timestamp.Nanos
		keyM.Timestamp.Seconds = keyMod.Timestamp.Seconds
		keyM.IsDelete = keyMod.IsDelete
		resultBytes, err := json.Marshal(keyM)
		if err != nil {
			return nil, errors.Error(errors.StatusBadRequest, err)
		}
		if resultBytes == nil {
			return nil, errors.Error(errors.StatusBadRequest, "Queried account doesn't have history.")
		}
		if !theFirstElement {
			buffer.WriteString(",")
		}
		buffer.Write(resultBytes)
		theFirstElement = false
	}
	buffer.WriteString("]")
	defer historyIter.Close()

	return buffer.Bytes(), nil
}

func registerPoints(stub shim.ChaincodeStubInterface, body []byte) ([]byte, error) {
	param := new(Points)
	if err := json.Unmarshal(body, param); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, err := govalidator.ValidateStruct(param); !ok {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}

	if pointsBytes, _ := stub.GetState(param.Name); pointsBytes != nil {
		//链上已有同名的Points
		return nil, errors.Errorf(errors.StatusBadRequest, "There already is a point named %v.on chain.", param.Name)
	}

	if param.Ceiling == 0 {
		param.Ceiling = baseForMaxCirculation
	}

	param.Circulation = 0 //积分创建时，发行量必须为0

	pBytes, _ := json.Marshal(param)
	if pBytes == nil {
		return nil, errors.Error(errors.StatusBadRequest, "The registered points is empty.")
	}

	if !param.HasKey() {
		return nil, errors.Error(errors.StatusBadRequest, "Points must have a name as the key.")
	}
	if err := stub.PutState(param.GetKey(), pBytes); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	return body, nil
}
