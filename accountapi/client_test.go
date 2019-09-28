package accountapi

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"os"
)

func getHostAddress() string {
	if value, ok := os.LookupEnv("API_HOST_ADDR"); ok {
		return value
	}
	return "http://localhost:8080"
}

var _ = Describe("The Form3 Accounts Client", func() {
	var (
		client   *Client
		expected *Account
		actual   *Account
		err      error
	)

	BeforeEach(func() {
		client = NewClient(getHostAddress())
		expected = &Account{
			ID:             uuid.New().String(),
			Type:           "accounts",
			OrganisationID: uuid.New().String(),
			Attributes: Attributes{
				Country:                     "GB",
				AccountClassification:       "Personal",
				BankID:                      "400300",
				Bic:                         "NWBKGB22",
				Iban:                        "GB11NWBK40030041426819",
				Title:                       "Ms",
				FirstName:                   "Samantha",
				BankAccountName:             "Samantha Holder",
				AlternativeBankAccountNames: []string{"Sam Holder"},
			},
		}
	})

	Describe("creating an account", func() {
		It("should create a new account if the ID is unique", func() {
			Expect(client.Create(expected)).Should(PointTo(MatchFields(IgnoreExtras, Fields{
				"ID":             Equal(expected.ID),
				"OrganisationID": Equal(expected.OrganisationID),
				"Attributes":     Equal(expected.Attributes),
			})))
		})

		It("should create a personal account classification", func() {

			expected.Attributes.AccountClassification = "Personal"
			Expect(client.Create(expected)).Should(PointTo(MatchFields(IgnoreExtras, Fields{
				"Attributes": MatchFields(IgnoreExtras, Fields{
					"AccountClassification": Equal("Personal"),
				}),
			})))

		})

		It("should create a business account classification", func() {
			expected.Attributes.AccountClassification = "Business"
			Expect(client.Create(expected)).Should(PointTo(MatchFields(IgnoreExtras, Fields{
				"Attributes": MatchFields(IgnoreExtras, Fields{
					"AccountClassification": Equal("Business"),
				}),
			})))
		})

		It("should create an account number if just iban is set", func() {
			expected.Attributes.Iban = "GB11NWBK40030041426819"
			Expect(client.Create(expected)).Should(PointTo(MatchFields(IgnoreExtras, Fields{
				"Attributes": MatchFields(IgnoreExtras, Fields{
					"AccountNumber": BeEmpty(),
				}),
			})))
		})

		It("should fail with an  error if the country code format is invalid", func() {
			expected.Attributes.Country = "gb"
			_, err = client.Create(expected)
			Expect(err).Should(HaveOccurred())
		})

		It("should fail with an  error for an unsupported account classification", func() {
			expected.Attributes.AccountClassification = "IncorrectClassification"
			_, err = client.Create(expected)
			Expect(err).Should(HaveOccurred())
		})

		It("should fail with an  error if the ID is not a valid UUID", func() {
			expected.ID = "invaliduuid"
			_, err = client.Create(expected)
			Expect(err).Should(HaveOccurred())
		})

		It("should fail with an  error if the Organisation ID is not a valid UUID", func() {
			expected.OrganisationID = "invaliduuid"
			_, err = client.Create(expected)
			Expect(err).Should(HaveOccurred())
		})

		It("should fail to create the account and report an error for duplicate IDs", func() {
			actual, err = client.Create(expected)
			actual, err = client.Create(expected)
			Expect(err).Should(HaveOccurred())
			Expect(actual).To(BeNil())
		})

	})

	Describe("fetching an account", func() {

		It("should retrieve the account if it exists with the given ID", func() {
			actual, err = client.Create(expected)
			Expect(client.Fetch(actual.ID)).Should(Equal(actual))
		})

		It("should error if there's no account with the given ID", func() {
			actual, err = client.Fetch(uuid.New().String())
			Expect(err).Should(HaveOccurred())
			Expect(actual).Should(BeNil())
		})

		Describe("deleting an account", func() {

			Context("given an existing account", func() {
				BeforeEach(func() {
					actual, err = client.Create(expected)
					Expect(actual).ShouldNot(BeNil())
				})

				It("should delete the account", func() {
					Expect(client.Delete(actual)).Should(Succeed())
					actual, err = client.Fetch(actual.ID)
					Expect(actual).Should(BeNil())
				})
			})

			Context("given the account has already been deleted", func() {
				var accountCreated *Account
				BeforeEach(func() {
					accountCreated, _ = client.Create(expected)
					Expect(accountCreated).ShouldNot(BeNil())
					Expect(client.Delete(accountCreated)).Should(Succeed())
					actual, _ = client.Fetch(accountCreated.ID)
					Expect(actual).Should(BeNil())
				})

				It("should delete the account", func() {
					Expect(client.Delete(accountCreated)).Should(Succeed())
				})
			})

		})

	})

	Describe("listing accounts", func() {

		Context("given multiple accounts exist", func() {
			BeforeEach(func() {
				// ensure the accounts service has a minimum number of accounts
				for i := 0; i < 11; i++ {
					uuid := uuid.New().String()
					expected.ID = uuid
					Expect(client.Create(expected)).Should(Succeed())
				}
			})
		})

		It("should list accounts per a page size and page number", func() {
			Expect(client.List(1, 0)).Should(HaveLen(1))
			Expect(client.List(3, 1)).Should(HaveLen(3))
			Expect(client.List(5, 2)).Should(HaveLen(5))
			Expect(client.List(10, 0)).Should(HaveLen(10))
		})

		It("should return all accounts if page size is 0", func() {
			actual, err := client.List(0, 0)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(actual)).Should(BeNumerically(">=", 10)) // at least the minium created
		})

		It("should error if page size isnegative ", func() {
			_, err := client.List(-1, 0)
			Expect(err).Should(HaveOccurred())
		})

		It("should error if page number is negative ", func() {
			_, err := client.List(0, -1)
			Expect(err).Should(HaveOccurred())
		})

	})

})
