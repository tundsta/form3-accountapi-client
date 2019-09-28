package accountapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

//Client configuration
type Client struct {
	baseURL string
}

// Attributes - acount attributes sub resource
type Attributes struct {
	AccountClassification       string   `json:"account_classification"`
	AccountMatchingOptOut       bool     `json:"account_matching_opt_out"`
	AccountNumber               string   `json:"account_number"`
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names"`
	BankAccountName             string   `json:"bank_account_name"`
	BankID                      string   `json:"bank_id"`
	BankIDCode                  string   `json:"bank_id_code"`
	BaseCurrency                string   `json:"base_currency"`
	Bic                         string   `json:"bic"`
	Country                     string   `json:"country"`
	FirstName                   string   `json:"first_name"`
	Iban                        string   `json:"iban"`
	JointAccount                bool     `json:"joint_account"`
	SecondaryIdentification     string   `json:"secondary_identification"`
	Title                       string   `json:"title"`
}

// Account resource defining acccount attributes and metadata
type Account struct {
	Attributes     Attributes `json:"attributes"`
	CreatedOn      time.Time  `json:"created_on"`
	ID             string     `json:"id"`
	ModifiedOn     time.Time  `json:"modified_on"`
	OrganisationID string     `json:"organisation_id"`
	Type           string     `json:"type"`
	Version        int        `json:"version"`
}

// accountBody resource body wrapping the Account when interacting with API
type accountBody struct {
	Account2 *Account `json:"data"`
	// Links Links `json:"links"`
}

// accountListBody wrapping a list of accounts when interacting with the API
type accountListBody struct {
	Accounts []*Account `json:"data"`
}

// NewClient - create a new client with specified config
func NewClient(baseURL string) *Client {
	return &Client{baseURL}
}

// Create a new account for the given object. Error is returned if a failure is encountered
func (c *Client) Create(a *Account) (*Account, error) {

	u := fmt.Sprintf("%v/v1/organisation/accounts", c.baseURL)
	req, err := c.newRequest("POST", u, &accountBody{a})
	if err != nil {
		return nil, err
	}

	var account Account
	b := accountBody{&account}
	err = c.doRequest(req, &b)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Delete the given account object. An error is thrown if this fails
func (c *Client) Delete(a *Account) error {
	u := fmt.Sprintf("%v/v1/organisation/accounts/%v?version=%v", c.baseURL, a.ID, a.Version)
	req, err := c.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	err = c.doRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// List all accounts by specifying the size of the page and page number to return.
func (c *Client) List(pageSize int, pageNumber int) ([]*Account, error) {
	u := fmt.Sprintf("%v/v1/organisation/accounts?page[size]=%v&page[number]=%v", c.baseURL, pageSize, pageNumber)

	req, err := c.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var b *accountListBody

	err = c.doRequest(req, &b)
	if err != nil {
		return nil, err
	}

	return b.Accounts, nil

}

// Fetch an account by a given account ID
func (c *Client) Fetch(id string) (*Account, error) {
	u := fmt.Sprintf("%v/v1/organisation/accounts/%v", c.baseURL, id)

	req, err := c.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var b *accountBody
	err = c.doRequest(req, &b)
	if err != nil {
		return nil, err
	}

	return b.Account2, nil

}

func (c *Client) newRequest(method, url string, body interface{}) (*http.Request, error) {
	var buf = new(bytes.Buffer)

	if body != nil {
		b, err := json.Marshal(body)
		buf = bytes.NewBuffer(b)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, v interface{}) error {

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if s := resp.StatusCode; s < 200 || s >= 300 {
		var str interface{}
		body, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &str)
		if err != nil {
			e := fmt.Errorf("Error unmarshaling response error: %v", err)
			return e
		}

		e := fmt.Errorf("%v Error: %v", resp.StatusCode, str)
		return e
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err == io.EOF {
		err = nil //  empty response body
	}
	if err != nil {
		return err
	}

	return nil
}
