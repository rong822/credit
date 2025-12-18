package points

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/shopspring/decimal"
)

type invokeFunc func(shim.ChaincodeStubInterface, []byte) ([]byte, error)

const (
	ChaincodeName = "scf_points"

	MethodIssuePoints         = "issuePoints"
	MethodQueryBalance        = "queryBalance"
	MethodTransferPoints      = "transferPoints"
	MethodQueryAccountHistory = "queryAccountHistory"
	MethodRegisterIssuer      = "registerPoints"
)

var (
	errFormat = "%v\n"
	logger    = shim.NewLogger("scf_cc_v1")

	invokeFunction = map[string]invokeFunc{
		MethodIssuePoints:         issuePoints,
		MethodQueryBalance:        queryBalance,
		MethodTransferPoints:      transferPoints,
		MethodQueryAccountHistory: queryAccountHistory,
		MethodRegisterIssuer:      registerPoints,
	}
)

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

type Chaincode struct{}

func (c *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Infof("invoke is running %s", function)

	if len(args) != 1 {
		return shim.Error("incorrect number of arguments. Expecting 1")
	}

	if invokeFunction[function] == nil {
		return shim.Error("received unknown function invocation")
	}

	success, err := invokeFunction[function](stub, []byte(args[0]))
	if err != nil {
		logger.Errorf(errFormat, err)
		return shim.Error(err.Error())
	}

	return shim.Success(success)
}
