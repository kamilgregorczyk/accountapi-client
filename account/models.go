package account

// Represents top level model for Organisation/Account model
// Look into https://api-docs.form3.tech/api.html#organisation-accounts-resource for more details
type Account struct {
	Id             string      `json:"id,omitempty"`
	Type           string      `json:"type,omitempty"`
	OrganisationId string      `json:"organisation_id,omitempty"`
	Attributes     *Attributes `json:"attributes"`
	Version        int64       `json:"version,omitempty"`
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
	JointAccount            *bool    `json:"joint_account,omitempty"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
	Status                  string   `json:"status,omitempty"`
	BankIdCode              string   `json:"bank_id_code,omitempty"`
}

type Data struct {
	Data *Account `json:"data"`
}
