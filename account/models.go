package account

import "time"

// Represents top level model for Organisation/Account model
// Look into https://api-docs.form3.tech/api.html#organisation-accounts-resource for more details
type Account struct {
	Id             string      `json:"id,omitempty"`
	Type           string      `json:"type,omitempty"`
	OrganisationId string      `json:"organisation_id,omitempty"`
	Attributes     *Attributes `json:"attributes"`
	Version        int         `json:"version,omitempty"`
	CreatedOn      time.Time   `json:"created_on,omitempty"`
	ModifiedOn     time.Time   `json:"modified_on,omitempty"`
}

// Represents custom attributes for Organisation/Account model as well as standard attributes
// Look into https://api-docs.form3.tech/api.html#organisation-accounts-resource for more details
//
// Deprecated fields were skipped
type Attributes struct {
	Country                 string   `json:"country"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	BankId                  string   `json:"bank_id,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	CustomerId              string   `json:"customer_id,omitempty"`
	Name                    []string `json:"name"`
	AlternativeNames        []string `json:"alternative_names"`
	AccountClassification   string   `json:"account_classification,omitempty"`
	JointAccount            bool     `json:"joint_account,omitempty"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Switched                bool     `json:"switched,omitempty"`
	Status                  string   `json:"status,omitempty"`
	BankIdCode              string   `json:"bank_id_code,omitempty"`
}

type CreateAccountRequest struct {
	Account *Account `json:"data"`
}

type CreateAccountResponse struct {
	Account *Account `json:"data"`
}

type FetchAccountRequest struct {
	Id string
}

func (r *FetchAccountRequest) Validate() error {
	if len(r.Id) <= 0 {
		return &ValidationError{Message: "id cannot be empty"}
	}
	return nil
}

type FetchAccountResponse struct {
	Account *Account `json:"data"`
}

type ListAccountsRequest struct {
	PageNumber int
	PageSize   int
}

func (r *ListAccountsRequest) Validate() error {
	if r.PageNumber < 0 {
		return &ValidationError{Message: "pageNumber has to larger than zero"}
	}

	if r.PageSize < 0 {
		return &ValidationError{Message: "pageSize has to larger than zero"}
	}
	return nil
}

type ListAccountResponse struct {
	Accounts []*Account `json:"data"`
}

type DeleteAccountRequest struct {
	Id      string
	Version int
}

func (r *DeleteAccountRequest) Validate() error {
	if len(r.Id) <= 0 {
		return &ValidationError{Message: "id cannot be empty"}
	}

	if r.Version < 0 {
		return &ValidationError{Message: "version has to be larger than zero"}
	}
	return nil
}
