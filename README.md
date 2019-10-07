# Form3 Accounts Client

Candidate Name: Olatunde Awoyemi

Client Library for the Form3 [Accounts API](https://api-docs.form3.tech/api.html?shell#organisation-accounts) managing Bank Accounts registered with Form3 for the validation and allocation of inbound payments.


## Usage

Create a new client and call an appropriate function.

```go
client = NewClient("http://localhost:8080")

account, err := client.Create(account &Account{
                                                ID:             uuid.New().String(),
                                                Type:           "accounts",
                                                OrganisationID: uuid.New().String(),
                                                Attributes: Attributes{
                                                    Country:    "GB",
                                                    ..
                                                    ..
                                              })
                                                    
```

|function|function signature|
|---|---|
|create Account|`func (c *Client) Create(a *Account) (*Account, error)`|
|fetch Account|`func (c *Client) Fetch(id string) (*Account, error)`|
|delete Account|`func (c *Client) Delete(a *Account) error`|
|list Accounts|`func (c *Client) List(pageSize int, pageNumber int) ([]*Account, error)`|

Note, errors returned from the API are represented simply as an `error` with error response body expressed as a string.

## Running tests

Integration tests are defined via the [Ginkgo](http://onsi.github.io/ginkgo/) BDD test framework.

Ensure the full stack is running via (this also executes the tests):
```bash
docker-compose up
```

To rerun the tests, either run via:

```bash
go test
```

or install [Ginkgo](https://onsi.github.io/ginkgo/#getting-ginkgo) and run:

```
ginkgo
```

