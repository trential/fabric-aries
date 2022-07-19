package main

import (
	"chaincode/logger"
	"chaincode/utils"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func handleNymTx(stub shim.ChaincodeStubInterface, caller *Nym, operation map[string]interface{}) (*Nym, error) {
	var req NymRequest
	err := req.from(operation)
	if err != nil {
		return nil, err
	}
	logger.Debugf(`
handleNymTx:
			alias=%s
			dest=%s
			role=%s
			verkey=%s`, req.Alias, req.Dest, req.Role, req.Verkey)

	logger.Debug("fetch nym transaction from worldstate")
	raw, err := stub.GetState(req.id())
	if err != nil {
		return nil, fmt.Errorf("%w: nym transaction", ErrWorldstateRead)
	}
	if len(raw) == 0 {
		// nym doesn't exists
		return handleCreateNymTx(stub, caller, &req)
	}
	// nym tx already exists, update
	var nym Nym
	json.Unmarshal(raw, &nym)
	return handleUpdateNymTx(stub, caller, &req, &nym)
}

func handleCreateNymTx(stub shim.ChaincodeStubInterface, caller *Nym, req *NymRequest) (*Nym, error) {
	const fnTag = "handleCreateNymTx"
	logger.Debugf("%s: checking caller role", fnTag)
	if caller.Role != TRUSTEE_ROLE {
		return nil, fmt.Errorf("%w: only TRUSTEE can create nym transaction", ErrInvalidRequest)
	}

	logger.Debugf("%s: getting target role", fnTag)
	var role NYM_ROLES
	switch req.Role {
	case TRUSTEE_ROLE_REQ:
		role = TRUSTEE_ROLE
	case ENDORSER_ROLE_REQ:
		role = ENDORSER_ROLE
	default:
		return nil, fmt.Errorf("%w: invalid target role", ErrInvalidRequest)
	}

	nym := Nym{
		Alias:  req.Alias,
		Did:    req.Dest,
		Role:   role,
		Verkey: req.Verkey,
		Ver:    "1.0",
	}
	logger.Debugf("%s: storing nym transaction in worldstate", fnTag)
	raw, _ := json.Marshal(nym)
	err := stub.PutState(req.id(), raw)
	if err != nil {
		return nil, fmt.Errorf("%w: nym transaction", ErrWorldstateWrite)
	}
	return &nym, nil
}

// TODO
func handleUpdateNymTx(stub shim.ChaincodeStubInterface, caller *Nym, req *NymRequest, nym *Nym) (*Nym, error) {
	return nil, fmt.Errorf("%w: update nym transaction", ErrUnimplemented)
}

func handleSchemaTx(stub shim.ChaincodeStubInterface, caller *Nym, operation map[string]interface{}) (*Schema, error) {
	var req SchemaRequest
	err := req.from(operation)
	if err != nil {
		return nil, err
	}
	logger.Debugf(`
handleSchemaTx:
				attr_names=%v
				name=%s
				version=%s`, req.Data.AttrNames, req.Data.Name, req.Data.Version)

	logger.Debug("fetch schema transaction from worldstate")
	raw, err := stub.GetState(req.id(caller.Did))
	if err != nil {
		return nil, fmt.Errorf("%w: schema transaction", ErrWorldstateRead)
	}

	if len(raw) == 0 {
		// schema doesn't exists
		return handleCreateSchemaTx(stub, caller, &req)
	}
	// schema tx already exists, update
	var schema Schema
	json.Unmarshal(raw, &schema)
	return handleUpdateSchemaTx(stub, caller, &req, &schema)
}

func handleCreateSchemaTx(stub shim.ChaincodeStubInterface, caller *Nym, req *SchemaRequest) (*Schema, error) {
	const fnTag = "handleCreateSchemaTx"
	logger.Debugf("%s: checking caller role", fnTag)
	if !(caller.Role == ENDORSER_ROLE || caller.Role == TRUSTEE_ROLE) {
		return nil, fmt.Errorf("%w: only TRUSTEE and ENDORSER can create schema transaction", ErrInvalidRequest)
	}

	id := req.id(caller.Did)
	schema := Schema{
		ID:        id,
		AttrNames: req.Data.AttrNames,
		Name:      req.Data.Name,
		Version:   req.Data.Version,
		Ver:       "1.0",
	}
	logger.Debugf("%s: storing schema transaction in worldstate", fnTag)
	raw, _ := json.Marshal(schema)
	err := stub.PutState(id, raw)
	if err != nil {
		return nil, fmt.Errorf("%w: schema transaction", ErrWorldstateWrite)
	}
	return &schema, nil
}

// TODO
func handleUpdateSchemaTx(stub shim.ChaincodeStubInterface, caller *Nym, req *SchemaRequest, schema *Schema) (*Schema, error) {
	return nil, fmt.Errorf("%w: update schema transaction", ErrUnimplemented)
}

func handleCredentialDefinitionTx(stub shim.ChaincodeStubInterface, caller *Nym, operation map[string]interface{}) (*CredentialDefinition, error) {
	var req CredentialDefinitionRequest
	err := req.from(operation)
	if err != nil {
		return nil, err
	}
	logger.Debugf(`
handleCredentialDefinitionTx:
							signature_type=%s
							tag=%s
							schemaId=%s`, req.SignatureType, req.Tag, req.SchemaID)
	logger.Debug("fetch credential definition transaction from worldstate")
	raw, err := stub.GetState(req.id(caller.Did))
	if err != nil {
		return nil, fmt.Errorf("%w: credential definition transaction", ErrWorldstateRead)

	}
	if len(raw) == 0 {
		// cred def tx doesn't exists, create new
		return handleCreateCredentialDefinitionTx(stub, caller, &req)
	}
	// cred def tx already exists, update
	var credDef CredentialDefinition
	json.Unmarshal(raw, &credDef)
	return handleUpdateCredentialDefinitionTx(stub, caller, &req, &credDef)
}

func handleCreateCredentialDefinitionTx(stub shim.ChaincodeStubInterface, caller *Nym, req *CredentialDefinitionRequest) (*CredentialDefinition, error) {
	const fnTag = "handleCreateCredentialDefinitionTx"
	logger.Debugf("%s: checking caller role", fnTag)

	if !(caller.Role == ENDORSER_ROLE || caller.Role == TRUSTEE_ROLE) {
		return nil, fmt.Errorf("%w: only TRUSTEE and ENDORSER can create schema transaction", ErrInvalidRequest)
	}
	logger.Debugf("%s: fetching schema", fnTag)
	raw, err := stub.GetState(req.SchemaID)
	if err != nil {
		return nil, fmt.Errorf("%w: schema transaction", ErrWorldstateRead)
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("%w: schema", ErrStateNotFound)
	}

	id := req.id(caller.Did)
	credDef := CredentialDefinition{
		ID:       id,
		SchemaID: req.SchemaID,
		Type:     req.SignatureType,
		Tag:      req.Tag,
		Value:    req.Data,
		Ver:      "1.0",
	}
	logger.Debugf("%s: storing credential definition transaction in worldstate", fnTag)
	raw, _ = json.Marshal(credDef)
	err = stub.PutState(id, raw)
	if err != nil {
		return nil, fmt.Errorf("%w: credential definition transaction", ErrWorldstateWrite)
	}
	return &credDef, nil
}

// TODO
func handleUpdateCredentialDefinitionTx(stub shim.ChaincodeStubInterface, caller *Nym, req *CredentialDefinitionRequest, credDef *CredentialDefinition) (*CredentialDefinition, error) {
	return nil, fmt.Errorf("%w: update credential definition transaction", ErrUnimplemented)

}

func handleReadIDRequest(stub shim.ChaincodeStubInterface, operation map[string]interface{}) (interface{}, error) {
	var req ReadIDRequest
	err := req.from(operation)
	if err != nil {
		return nil, err
	}
	logger.Debugf(`
	handleReadIDRequest:
						id=%s`, req.ID)
	raw, err := stub.GetState(req.id())
	if err != nil {
		return nil, fmt.Errorf("%w: read transaction", ErrWorldstateRead)
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("%w: transaction not found", ErrStateNotFound)
	}
	var out interface{}
	json.Unmarshal(raw, &out)
	return out, nil
}

// authenticate: authenticate the caller and returns it
func authenticate(stub shim.ChaincodeStubInterface, req Request, reqMsg string) (*Nym, error) {
	// get caller did
	raw, err := stub.GetState(req.Identifier)
	if err != nil {
		return nil, err
	}
	// caller did not found
	if len(raw) == 0 {
		return nil, fmt.Errorf("%w: DID not found", ErrVerifySignature)
	}
	var nym Nym
	json.Unmarshal(raw, &nym)
	valid := utils.VerifyRequest(nym.Verkey, reqMsg, req.Signature)
	if !valid {
		return nil, fmt.Errorf("%w: invalid signature", ErrVerifySignature)
	}
	return &nym, nil
}

func getTxType(typ interface{}) (TX_TYPE, error) {
	if typ == nil {
		return TX_TYPE(""), fmt.Errorf("%w: tx type not provided", ErrInvalidRequest)
	}
	tp, ok := typ.(string)
	if !ok {
		return TX_TYPE(""), fmt.Errorf("%w: invalid tx type", ErrInvalidRequest)
	}
	txp := TX_TYPE(tp)
	if !(txp == NYM_TX ||
		txp == SCHEMA_TX ||
		txp == CRED_DEF_TX) {
		return txp, fmt.Errorf("%w: tx type not supported", ErrInvalidRequest)
	}
	return txp, nil
}

//
func getTxTypeFromID(id string) TX_TYPE {
	parts := strings.Split(id, ":")
	if len(parts) == 0 {
		return TX_TYPE("")
	} else if len(parts) == 1 {
		return NYM_TX
	}
	var out TX_TYPE
	switch parts[1] {
	case "2":
		out = SCHEMA_TX
	case "3":
		out = CRED_DEF_TX
	}
	return out
}
