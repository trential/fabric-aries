package main

import (
	"chaincode/logger"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

func main() {
	logger.SetupLogger("SSI", logger.DEBUG_LEVEL)
	logger.Info("Start chaincode")
	cc := SSIChaincode{}

	if err := shim.Start(cc); err != nil {
		return
	}
}

type SSIChaincode struct{}

func (SSIChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error(fmt.Errorf("%w: require one argument(genesis)", ErrInvalidRequest).Error())
	}
	var nyms []NymRequest
	err := json.Unmarshal([]byte(args[0]), &nyms)
	if err != nil {
		return shim.Error(fmt.Errorf("%w: invalid request object", ErrInvalidRequest).Error())
	}
	logger.Info("creating genesis trustees")
	for _, nym := range nyms {
		if nym.Dest == "" || nym.Verkey == "" {
			return shim.Error(fmt.Errorf("%w: dest or verkey should not be empty", ErrInvalidRequest).Error())
		}
		_, err := handleCreateNymTx(stub, &Nym{
			Role: TRUSTEE_ROLE,
		}, &NymRequest{
			Alias:  nym.Alias,
			Dest:   nym.Dest,
			Role:   TRUSTEE_ROLE_REQ,
			Type:   NYM_TX,
			Verkey: nym.Verkey,
		})
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	return shim.Success(nil)
}

func (SSIChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, args := stub.GetFunctionAndParameters()

	var err error
	var op RESPONSE_TYPE
	var txType TX_TYPE
	var identifier string
	var data interface{}
	switch fcn {
	case "read":
		op = READ_RESPONSE_TYPE
		if len(args) != 1 {
			err = fmt.Errorf("%w: require one argument(txID)", ErrInvalidRequest)
			break
		}
		id := args[0]
		logger.Debug("get transaction type from id")
		txType = getTxTypeFromID(id)
		data, err = read(stub, id)
	case "write":
		op = WRITE_RESPONSE_TYPE
		if len(args) != 1 {
			err = fmt.Errorf("%w: require one argument(singedLedgerRequest)", ErrInvalidRequest)
			break
		}
		logger.Debug("unmarshal SignedRequestLedger")
		var req Request
		err = json.Unmarshal([]byte(args[0]), &req)
		if err != nil {
			err = fmt.Errorf("%w: invalid ledger request: %v", ErrInvalidRequest, err)
			break
		}
		identifier = req.Identifier
		logger.Debug("authenticate caller")
		var caller *Nym
		caller, err = authenticate(stub, req, args[0])
		if err != nil {
			break
		}
		logger.Debug("get tx type from operations")
		txType, err = getTxType(req.Operation)
		if err != nil {
			break
		}
		logger.Debugf("caller: %s", identifier)
		data, err = write(stub, txType, caller, req.Operation)
	default:
		err = fmt.Errorf("%s: %w", fcn, ErrFunctionNotSupported)
	}

	if err != nil {
		raw, _ := json.Marshal(ErrorResponse{
			Op:         op,
			Type:       txType,
			Identifier: identifier,
			Reason:     err.Error(),
		})
		logger.Info(string(raw))
		return peer.Response{
			Status:  shim.ERROR,
			Payload: raw,
		}
	}
	raw, _ := json.Marshal(SuccessResponse{
		Op:         op,
		Type:       txType,
		Identifier: identifier,
		Data:       data,
	})
	return shim.Success(raw)
}

func write(stub shim.ChaincodeStubInterface, txType TX_TYPE, caller *Nym, operation map[string]interface{}) (interface{}, error) {
	var data interface{}
	var err error
	switch txType {
	case NYM_TX:
		data, err = handleNymTx(stub, caller, operation)
	case SCHEMA_TX:
		data, err = handleSchemaTx(stub, caller, operation)
	case CRED_DEF_TX:
		data, err = handleCredentialDefinitionTx(stub, caller, operation)
	}
	return data, err
}

func read(stub shim.ChaincodeStubInterface, id string) (interface{}, error) {
	logger.Debugf("read transaction id: %s", id)
	raw, err := stub.GetState(id)
	if err != nil {
		return nil, ErrWorldstateRead
	}
	if len(raw) == 0 {
		return nil, ErrStateNotFound
	}
	var data map[string]interface{}
	json.Unmarshal(raw, &data)
	return data, nil
}
