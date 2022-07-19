package main

// ssiType.go defines nym, schema, credential definition type
// for storing on worldstate

// Nym transaction

type NYM_ROLES string

const (
	TRUSTEE_ROLE  NYM_ROLES = "TRUSTEE"
	ENDORSER_ROLE NYM_ROLES = "ENDORSER"
)

type Nym struct {
	Alias  string    `json:"alias,omitempty"`
	Did    string    `json:"did"`
	Role   NYM_ROLES `json:"role,omitempty"`
	Verkey string    `json:"verkey"`
	Ver    string    `json:"ver"`
}

// Schema transaction

type Schema struct {
	ID        string   `json:"id"`
	AttrNames []string `json:"attrNames"`
	Name      string   `json:"name"`
	Version   string   `json:"version"`
	Ver       string   `json:"ver"`
}

// Credential definition

type CredentialDefinition struct {
	ID       string                    `json:"id"`
	SchemaID string                    `json:"schemaId"`
	Type     string                    `json:"type"`
	Tag      string                    `json:"tag"`
	Value    *CredentialDefinitionData `json:"value"`
	Ver      string                    `json:"ver"`
}

type CredentialDefinitionData struct {
	Primary    *Primary    `json:"primary"`
	Revocation *Revocation `json:"revocation,omitempty"`
}

type Primary struct {
	N     string                 `json:"n"`
	R     map[string]interface{} `json:"r"`
	Rctxt string                 `json:"rctxt"`
	S     string                 `json:"s"`
	Z     string                 `json:"z"`
}

type Revocation struct {
	G      string `json:"g"`
	G_dash string `json:"g_dash"`
	H      string `json:"h"`
	H0     string `json:"h0"`
	H1     string `json:"h1"`
	H2     string `json:"h2"`
	H_cap  string `json:"h_cap"`
	Htilde string `json:"htilde"`
	Pk     string `json:"pk"`
	U      string `json:"u"`
	Y      string `json:"y"`
}
