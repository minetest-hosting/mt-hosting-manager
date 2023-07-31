package wallee

type LineItemType string

const (
	LineItemTypeFee     LineItemType = "FEE"
	LineItemTypeProduct LineItemType = "PRODUCT"
)

type TransactionStateType string

const (
	TransactionStatePending   TransactionStateType = "PENDING"
	TransactionStateConfirmed TransactionStateType = "CONFIRMED"
	TransactionStateFulfilled TransactionStateType = "FULFILL"
)

type LineItem struct {
	AmountIncludingTax float64      `json:"amountIncludingTax"`
	Name               string       `json:"name"`
	Quantity           int          `json:"quantity"`
	Type               LineItemType `json:"type"`
	UniqueID           string       `json:"uniqueId"`
	ID                 *int64       `json:"id"`
}

type TransactionRequest struct {
	Currency   string      `json:"currency"`
	LineItems  []*LineItem `json:"lineItems"`
	FailedURL  string      `json:"failedUrl"`
	SuccessURL string      `json:"successUrl"`
}

type TransactionResponse struct {
	Currency  string               `json:"currency"`
	ID        int64                `json:"id"`
	LineItems []*LineItem          `json:"lineItems"`
	State     TransactionStateType `json:"state"`
}

type FilterOperatorType string

const FilterOperatorEquals FilterOperatorType = "EQUALS"

type FilterTypeType string

const FilterTypeLeaf FilterTypeType = "LEAF"

type TransactionSearchFilter struct {
	FieldName string             `json:"fieldname"`
	Operator  FilterOperatorType `json:"operator"`
	Type      FilterTypeType     `json:"type"`
	Value     string             `json:"value"`
}

type TransactionSearchRequest struct {
	Filter *TransactionSearchFilter `json:"filter"`
}

type TransactionSearchResponse []*TransactionResponse
