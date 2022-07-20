package main

import (
	"encoding/json"
	"fmt"
)

// ledger.go defines request and response for interacting with
// hl fabric for ssi
type TX_TYPE string

const (
	NYM_TX      TX_TYPE = "1"
	SCHEMA_TX   TX_TYPE = "101"
	CRED_DEF_TX TX_TYPE = "102"
)

type Request struct {
	Identifier string                 `json:"identifier"`
	Operation  map[string]interface{} `json:"operation"`
	Signature  string                 `json:"signature"`
}

type ReadRequest struct {
	Type TX_TYPE `json:"type"`
}

type NYM_ROLE_REQ string

const (
	TRUSTEE_ROLE_REQ  NYM_ROLE_REQ = "0"
	ENDORSER_ROLE_REQ NYM_ROLE_REQ = "101"
)

type RESPONSE_TYPE string

const (
	READ_RESPONSE_TYPE  RESPONSE_TYPE = "READ"
	WRITE_RESPONSE_TYPE RESPONSE_TYPE = "WRITE"
)

type SuccessResponse struct {
	Op   RESPONSE_TYPE `json:"op"`
	Type TX_TYPE       `json:"type"`
	// anyone in network can read, so read transaction will not include
	// identifier field
	Identifier string      `json:"identifier,omitempty"`
	Data       interface{} `json:"data"`
}

type ErrorResponse struct {
	Op   RESPONSE_TYPE `json:"op"`
	Type TX_TYPE       `json:"type"`
	// anyone in network can read, so read transaction will not include
	// identifier field
	Identifier string `json:"identifier,omitempty"`
	Reason     string `json:"reason"`
}

// tx data in operation
//
// nym request
type NymRequest struct {
	Alias  string       `json:"alias,omitempty"`
	Dest   string       `json:"dest"`
	Role   NYM_ROLE_REQ `json:"role,omitempty"`
	Type   TX_TYPE      `json:"type"`
	Verkey string       `json:"verkey"`
}

// from : unmarshal operation into NymRequest
func (r *NymRequest) from(operation map[string]interface{}) error {
	raw, _ := json.Marshal(operation)
	err := json.Unmarshal(raw, r)
	if err != nil {
		return err
	}
	if r.Dest == "" {
		return fmt.Errorf("%w: 'dest' is empty", ErrInvalidRequest)
	} else if r.Verkey == "" {
		return fmt.Errorf("%w: 'verkey' is empty", ErrInvalidRequest)
	}

	return nil
}

// id: returns key of nym to fetch from worldstate
func (r *NymRequest) id() string {
	return r.Dest
}

//
// schema request
type SchemaRequest struct {
	Type TX_TYPE     `json:"type"`
	Data *SchemaData `json:"data"`
}

type SchemaData struct {
	AttrNames []string `json:"attr_names"`
	Name      string   `json:"name"`
	Version   string   `json:"version"`
}

// from : unmarshal operation into SchemaRequest
func (r *SchemaRequest) from(operation map[string]interface{}) error {
	raw, _ := json.Marshal(operation)
	err := json.Unmarshal(raw, r)
	if err != nil {
		return err
	}

	if r.Data == nil {
		return fmt.Errorf("%w: 'data' is null", ErrInvalidRequest)
	} else if len(r.Data.AttrNames) == 0 {
		return fmt.Errorf("%w: schema 'attr_names' is empty", ErrInvalidRequest)
	} else if r.Data.Name == "" {
		return fmt.Errorf("%w: schema 'name' is empty", ErrInvalidRequest)
	} else if r.Data.Version == "" {
		return fmt.Errorf("%w: schema 'version' is empty", ErrInvalidRequest)
	}

	return nil
}

// id: returns key of schema to fetch from worldstate
func (r *SchemaRequest) id(caller string) string {
	return fmt.Sprintf("%s:2:%s:%s", caller, r.Data.Name, r.Data.Version)
}

//
// credential definition request
type CredentialDefinitionRequest struct {
	Data          *CredentialDefinitionData `json:"data"`
	Type          TX_TYPE                   `json:"type"`
	SignatureType string                    `json:"signature_type"`
	Tag           string                    `json:"tag"`
	SchemaID      string                    `json:"schemaId"`
}

func (r *CredentialDefinitionRequest) from(operation map[string]interface{}) error {
	raw, _ := json.Marshal(operation)
	err := json.Unmarshal(raw, r)
	if err != nil {
		return err
	}

	if r.SchemaID == "" {
		return fmt.Errorf("%w: 'schemaId' is empty", ErrInvalidRequest)
	} else if r.SignatureType != "CL" {
		return fmt.Errorf("%w: only CL signature_type is supported", ErrInvalidRequest)
	} else if r.Tag == "" {
		return fmt.Errorf("%w: 'tag' is empty", ErrInvalidRequest)
	} else if r.Data == nil {
		return fmt.Errorf("%w: 'data' is null", ErrInvalidRequest)
	}

	return nil
}

func (r *CredentialDefinitionRequest) id(caller string) string {
	return fmt.Sprintf("%s:3:%s:%s:%s", caller, r.SignatureType, r.SchemaID, r.Tag)
}

//
type ReadIDRequest struct {
	ID string `json:"id"`
}

func (r *ReadIDRequest) from(operation string) error {
	err := json.Unmarshal([]byte(operation), r)
	if err != nil {
		return err
	}

	if r.ID == "" {
		return fmt.Errorf("%w: 'id' is empty", ErrInvalidRequest)
	}

	return nil
}

func (r *ReadIDRequest) id() string {
	return r.ID
}
