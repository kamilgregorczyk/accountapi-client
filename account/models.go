package account

type Account struct {
	Id             string             `json:"id,omitempty"`
	Type           string             `json:"type,omitempty"`
	OrganisationId string             `json:"organisation_id,omitempty"`
	Attributes     *AccountAttributes `json:"attributes"`
	Version        int64              `json:"version,omitempty"`
}

// Deprecated fields were skipped
type AccountAttributes struct {
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

type AccountData struct {
	Data *Account `json:"data"`
}
