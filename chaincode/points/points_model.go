package points

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/btcsuite/btcutil/base58"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gitlab.bigtree.com/dashu-blockchain/credit/pkg/errors"
)

const (
	baseForMaxCirculation         = 60000000
	yearlyGrowthForMAxCirculation = 0.05
)

type Account struct {
	AccountName string     `json:"accountName,omitempty" valid:"required"`
	Balances    []*balance `json:"balances,omitempty"`
	TxCount     uint       `json:"txCount"`
}

type balance struct {
	PointsName    string `json:"pointsName,omitempty" valid:"required"`
	PointsBalance uint32 `json:"pointsBalance"`
}

type Points struct {
	Name        string `json:"name,omitempty" valid:"required"`
	Circulation uint64 `json:"circulation"`
	Ceiling     uint64 `json:"ceiling"`
	Issuer      string `json:"issuer,omitempty" valid:"required"` //Issuer应为签发者公钥的SHA256哈希
	IssueCount  uint   `json:"issueCount"`
}

func newAccount() *Account {
	return &Account{
		Balances: make([]*balance, 0),
	}
}

func (a *Account) HasKey() bool {
	return a.AccountName != ""
}

func (a *Account) GetKey() string {
	return a.AccountName
}

func (p *Points) HasKey() bool {
	return p.Name != ""
}

func (p *Points) GetKey() string {
	return p.Name
}

func (p *Points) isHigherThanCeiling() bool {
	return p.Circulation > p.Ceiling
}

func makeHash(key string) string {
	sha := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sha[:])
}

func createAccount(stub shim.ChaincodeStubInterface, a *Account) ([]byte, error) {
	a.TxCount = 0 //TxCount在账户创建时必须为0
	accBytes, _ := json.Marshal(a)
	if accBytes == nil {
		return nil, errors.Error(errors.StatusBadRequest, "The account is empty.")
	}

	if !a.HasKey() {
		return nil, errors.Error(errors.StatusBadRequest, "Account must have a name as the key.")
	}
	if err := stub.PutState(a.GetKey(), accBytes); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}

	return accBytes, nil
}

//RSA的
//func extractPubKey(pub string) (*rsa.PublicKey, error) {
//	block, _ := pem.Decode([]byte(pub))
//	if block == nil {
//		return nil, errors.Error(errors.StatusBadRequest, "Pem decode error.")
//	}
//
//	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		return nil, errors.Error(errors.StatusBadRequest, "Could not parse DER encoded public key (encryption key)")
//	}
//	publicKey, isRSAPublicKey := publicKeyInterface.(*rsa.PublicKey)
//	if !isRSAPublicKey {
//		return nil, errors.Error(errors.StatusBadRequest, "Public key parsed is not an RSA public key")
//	}
//	return publicKey, nil
//}

func extractPubKey(pub string) *ecdsa.PublicKey {
	pubBytes := base58.Decode(pub)
	x, y := elliptic.Unmarshal(elliptic.P256(), pubBytes)
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
}

//// RSA
// *PSSOptions == nil,
//func verifySig(pub *rsa.PublicKey, sig []byte, hashed []byte) bool {
//	hash := crypto.SHA256
//	//opts := &rsa.PSSOptions{
//	//	SaltLength: -1,
//	//}
//	return rsa.VerifyPKCS1v15(pub, hash, hashed, sig) == nil
//}

func verifySig(pub *ecdsa.PublicKey, sig string, hashed []byte) bool {
	r, s, err := extractSignature(sig)
	if err != nil {
		return false
	}
	result := ecdsa.Verify(pub, hashed, r, s)
	return result
}

func extractSignature(signature string) (rint, sint *big.Int, err error) {
	byterun, err := hex.DecodeString(signature)
	if err != nil {
		err = errors.New(errors.StatusBadRequest, "decrypt error, "+err.Error())
		return
	}
	r, s, err := UnmarshalECDSASignature(byterun)
	return r, s, err
}

func UnmarshalECDSASignature(raw []byte) (*big.Int, *big.Int, error) {
	sig := new(ECDSASignature)
	_, err := asn1.Unmarshal(raw, sig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed unmashalling signature [%s]", err)
	}

	if sig.R == nil {
		return nil, nil, errors.New(errors.StatusBadRequest, "invalid signature, R must be different from nil")
	}
	if sig.S == nil {
		return nil, nil, errors.New(errors.StatusBadRequest, "invalid signature, S must be different from nil")
	}

	if sig.R.Sign() != 1 {
		return nil, nil, errors.New(errors.StatusBadRequest, "invalid signature, R must be larger than zero")
	}
	if sig.S.Sign() != 1 {
		return nil, nil, errors.New(errors.StatusBadRequest, "invalid signature, S must be larger than zero")
	}

	return sig.R, sig.S, nil
}

type ECDSASignature struct {
	R, S *big.Int
}

func validateSignature(proposal interface{}, pub string, sig string) (bool, error) {
	pubKey := extractPubKey(pub)
	if pubKey == nil {
		return false, errors.Error(errors.StatusBadRequest, "Extracted public key is empty.")
	}

	proposalBytes, err := json.Marshal(proposal)
	if err != nil {
		return false, err
	}
	if proposalBytes == nil {
		return false, errors.Error(errors.StatusBadRequest, "The input params are empty.")
	}
	h := crypto.Hash.New(crypto.SHA256)
	h.Write(proposalBytes)
	hashed := h.Sum(nil)
	//decodedSig, _ := hex.DecodeString(sig) // 从16进制转码
	if !verifySig(pubKey, sig, hashed) {
		return false, errors.Error(errors.StatusBadRequest, "Verification failed.")
	}
	return true, nil
}

func multiStepValidate(stub shim.ChaincodeStubInterface, b *balance, nonce uint, pub string) ([]byte, error) {
	pubKey := extractPubKey(pub)
	if pubKey == nil {
		return nil, errors.Error(errors.StatusBadRequest, "Extracted public key is empty.")
	}

	if !validateIssuer(stub, b.PointsName, makeHash(pub)) {
		return nil, errors.Errorf(errors.StatusBadRequest, "Pub key doesn't match the corresponding issuer for %v.", b.PointsName)
	}
	targetPoints, err := queryPointsByKey(stub, b.PointsName)
	if err != nil {
		return nil, err
	}
	if targetPoints == nil {
		return nil, errors.Error(errors.StatusBadRequest, "Points doesn't exit.")
	}
	//判断nonce
	if targetPoints.IssueCount != nonce-1 {
		return nil, errors.Error(errors.StatusBadRequest, "The nonce doesn't match.")
	}
	targetPoints.IssueCount += 1
	//判断是否超过发放总量限制
	targetPoints.Circulation += uint64(b.PointsBalance)
	newCeiling, err := latestCeiling(stub, targetPoints.Name)
	if err != nil {
		return nil, err
	}
	targetPoints.Ceiling = newCeiling
	if targetPoints.isHigherThanCeiling() {
		return nil, errors.Error(errors.StatusBadRequest, "Circulation is higher than the ceiling.")
	}
	resultBytes, err := updatePoints(stub, targetPoints)
	if err != nil {
		return nil, err
	}
	return resultBytes, nil
}

//判断账户是否存在
func validateAccount(stub shim.ChaincodeStubInterface, acountKey string) bool {
	if _, err := queryAccountByKey(stub, acountKey); err == nil {
		return true
	} else {
		return false
	}
}

func validateIssuer(stub shim.ChaincodeStubInterface, pointsName string, issuer string) bool {
	pointBytes, err := stub.GetState(pointsName)
	if err != nil {
		return false
	}
	if pointBytes == nil {
		return false
	}
	p := new(Points)
	if err := json.Unmarshal(pointBytes, p); err != nil {
		return false
	}
	if ok, _ := govalidator.ValidateStruct(p); !ok {
		return false
	}

	return p.Issuer == issuer
}

func validateTransferNonce(stub shim.ChaincodeStubInterface, accountKey string, nonce uint) bool {
	payer, _ := queryAccountByKey(stub, accountKey)
	if payer == nil {
		return false
	}
	if payer.TxCount == nonce-1 {
		return true
	} else {
		return false
	}
}

func queryAccountByKey(stub shim.ChaincodeStubInterface, accountKey string) (*Account, error) {
	accountBytes, err := stub.GetState(accountKey)
	if err != nil {
		return nil, err
	}
	if accountBytes == nil {
		return nil, errors.ErrNotFound
	}

	account := new(Account)
	if err = json.Unmarshal(accountBytes, account); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, err := govalidator.ValidateStruct(account); !ok {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}

	return account, nil
}

func queryBalanceByKey(stub shim.ChaincodeStubInterface, accountKey string, pointsName string) (*balance, error) {
	accountBytes, err := stub.GetState(accountKey)
	if err != nil {
		return nil, err
	}
	if accountBytes == nil {
		return nil, errors.ErrNotFound
	}

	account := new(Account)
	if err = json.Unmarshal(accountBytes, account); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, err := govalidator.ValidateStruct(account); !ok {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}

	for _, b := range account.Balances {
		if b.PointsName == pointsName {
			return b, nil
		}
	}

	return nil, errors.ErrNotFound
}

func queryAllBalances(stub shim.ChaincodeStubInterface, accountKey string) (*Account, error) {
	accountBytes, err := stub.GetState(accountKey)
	if err != nil {
		return nil, err
	}
	if accountBytes == nil {
		return nil, errors.ErrNotFound
	}

	account := new(Account)
	if err = json.Unmarshal(accountBytes, account); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, err := govalidator.ValidateStruct(account); !ok {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}

	return account, nil
}

func queryPointsByKey(stub shim.ChaincodeStubInterface, pointsName string) (*Points, error) {
	// 链上查找相关积分，若无则返回nil, nil
	pointsBytes, err := stub.GetState(pointsName)
	if err != nil {
		return nil, err
	}
	if pointsBytes == nil {
		return nil, err //若链上key不存在，会在此处返回nil, nil
	}

	resultPoints := new(Points)
	if err := json.Unmarshal(pointsBytes, resultPoints); err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	ceiling, err := latestCeiling(stub, pointsName)
	if err != nil {
		return nil, err
	}
	resultPoints.Ceiling = ceiling

	return resultPoints, nil
}

//判断账户中是否存在此类积分
func isPointsInAccount(stub shim.ChaincodeStubInterface, accountKey string, pointsName string) bool {
	if _, err := queryBalanceByKey(stub, accountKey, pointsName); err == nil {
		return true
	} else {
		return false
	}
}

func updatePoints(stub shim.ChaincodeStubInterface, points *Points) ([]byte, error) {
	pointsBytes, err := json.Marshal(points)
	if err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if pointsBytes == nil {
		return nil, errors.ErrNotFound
	}

	if err := stub.PutState(points.Name, pointsBytes); err != nil {
		return nil, err
	}
	return pointsBytes, nil
}

//更新积分发放上限数额，每年以固定速度增长
func latestCeiling(stub shim.ChaincodeStubInterface, PointsName string) (uint64, error) {
	currentTime := time.Now()
	defaultCeiling := baseForMaxCirculation * math.Pow(1+yearlyGrowthForMAxCirculation, float64(currentTime.Year()-2018))
	pointsBytes, err := stub.GetState(PointsName)
	if err != nil {
		return 0, err
	}
	if pointsBytes == nil {
		return 0, errors.Error(errors.StatusBadRequest, "Points name donesn't match any data on chain.")
	}

	points := new(Points)
	if err := json.Unmarshal(pointsBytes, points); err != nil {
		return 0, errors.Error(errors.StatusBadRequest, err)
	}
	if ok, err := govalidator.ValidateStruct(points); !ok {
		return 0, errors.Error(errors.StatusBadRequest, err)
	}

	if points.Ceiling == 0 {
		return uint64(defaultCeiling), nil
	}

	return points.Ceiling, nil
}

//传入增量数据，plus=true时是余额增加，plus=false时是余额减少
func updateAccount(stub shim.ChaincodeStubInterface, a *Account, plus bool) ([]byte, error) {
	account, err := queryAccountByKey(stub, a.GetKey())
	if err != nil {
		return nil, err
	}

	for _, b := range a.Balances {
		//更新余额
		if !isPointsInAccount(stub, account.GetKey(), b.PointsName) { //若链上账户不存在此类积分，为其添加新积分种类
			account.Balances = append(account.Balances, b)
		} else {
			for i, balance := range account.Balances {
				if balance.PointsName == b.PointsName && plus {
					//如果增量数据是增加
					account.Balances[i].PointsBalance += b.PointsBalance
				} else if balance.PointsName == b.PointsName && !plus {
					//如果增量数据是减少
					account.Balances[i].PointsBalance = account.Balances[i].PointsBalance - b.PointsBalance
				}
			}
		}
	}
	//更新TxCount
	if a.TxCount == 1 {
		account.TxCount += 1
	}

	accountBytes, err := json.Marshal(account)
	if err != nil {
		return nil, errors.Error(errors.StatusBadRequest, err)
	}
	if err := stub.PutState(account.GetKey(), accountBytes); err != nil {
		return nil, err
	}
	return accountBytes, nil
}

//func updateCount(stub shim.ChaincodeStubInterface, input interface{}) ([]byte, error) { //每次转让积分后调用此函数使txCount加1
//	var resultBytes []byte
//	var key string
//	switch i := input.(type) {
//	case *Account:
//		result, err := queryAccountByKey(stub, i.AccountName)
//		if err != nil {
//			return nil, err
//		}
//		result.TxCount += 1
//		key = result.GetKey()
//		resultBytes, err = json.Marshal(result)
//		if err != nil {
//			return nil, errors.Error(errors.StatusBadRequest, err)
//		}
//	case *Points:
//		result, err := queryPointsByKey(stub, i.Name)
//		if err != nil {
//			return nil, err
//		}
//		result.IssueCount += 1
//		key = result.GetKey()
//		resultBytes, err = json.Marshal(result)
//		if err != nil {
//			return nil, errors.Error(errors.StatusBadRequest, err)
//		}
//	default:
//		return nil, errors.Error(errors.StatusBadRequest, "Error in updateCount.")
//	}
//
//	if err := stub.PutState(key, resultBytes); err != nil {
//		return nil, err
//	}
//	return resultBytes, nil
//
//}
