/*
Copyright (C) 2025 github.com/go-schwab

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, see
<https://www.gnu.org/licenses/>.
*/

package trader

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
)

var (
	accountEndpoint string = "https://api.schwabapi.com/trader/v1"

	// Accounts
	endpointAccountNumbers string = accountEndpoint + "/accounts/accountNumbers"
	endpointAccounts       string = accountEndpoint + "/accounts"
	endpointAccount        string = accountEndpoint + "/accounts/%s"
	// endpointUserPreference string = accountEndpoint + "/userPreference"

	// Orders
	endpointOrders        string = accountEndpoint + "/orders"
	endpointAccountOrders string = accountEndpoint + "/accounts/%s/orders"
	endpointAccountOrder  string = accountEndpoint + "/accounts/%s/orders/%s"
	// endpointPreviewOrder  string = accountEndpoint + "/accounts/%s/previewOrder"

	// Transactions
	// endpointTransactions string = accountEndpoint + "/accounts/%s/transactions"
	endpointTransaction string = accountEndpoint + "/accounts/%s/transactions/%s"
)

type Transaction struct {
	ActivityId     int
	Time           string
	User           User
	Description    string
	AccountNumber  string
	Type           string
	Status         string
	SubAccount     string
	TradeDate      string
	SettlementDate string
	PositionId     int
	OrderId        int
	NetAmount      int
	ActivityType   string
	TransferItems  TransferItems
}

type User struct {
	CdDomainId     string
	Login          string
	Type           string
	UserId         int
	SystemUserName string
	FirstName      string
	LastName       string
	BrokerRepCode  string
}

type TransferItems struct {
	Instrument     InstrumentRef
	Amount         int
	Cost           int
	Price          float64
	FeeType        string
	PositionEffect string
}

type InstrumentRef struct {
	Cusip        string
	Symbol       string
	Description  string
	InstrumentId int
	NetChange    int
	Type         string
}

type AccountNumbers struct {
	AccountNumber string
	HashValue     string
}

type Account struct {
	SecuritiesAccount SecuritiesAccount
	AggregatedBalance AggregatedBalance
}

type SecuritiesAccount struct {
	Type                    string
	AccountNumber           string
	RoundTrips              int
	IsDayTrader             bool
	IsClosingOnlyRestricted bool
	PFCBFlag                bool
	Positions               []Position
	InitialBalances         InitialBalance
	CurrentBalances         CurrentBalance
	ProjectedBalances       ProjectedBalance
}

type Position struct {
	ShortQuantity                  float64
	AveragePrice                   float64
	CurrentDayProfitLoss           float64
	CurrentDayProfitLossPercentage float64
	LongQuantity                   float64
	SettledLongQuantity            float64
	SettledShortQuantity           float64
	AgedQuantity                   float64
	Instrument                     AccountInstrument
	MarketValue                    float64
	MaintenanceRequirement         float64
	AverageLongPrice               float64
	AverageShortPrice              float64
	TaxLotAverageLongPrice         float64
	TaxLotAverageShortPrice        float64
	LongOpenProfitLoss             float64
	ShortOpenProfitLoss            float64
	PreviousSessionLongQuantity    float64
	PreviousSessionShortQuantity   float64
	CurrentDayCost                 float64
}

type AccountInstrument struct {
	AssetType    string
	Cusip        string
	Symbol       string
	Description  string
	InstrumentID int
	NetChange    float64
}

type InitialBalance struct {
	AccruedInterest                  float64
	AvailableFundsNonMarginableTrade float64
	BondValue                        float64
	BuyingPower                      float64
	CashBalance                      float64
	CashAvailableForTrading          float64
	CashReceipts                     float64
	DayTradingBuyingPower            float64
	DayTradingBuyingPowerCall        float64
	DayTradingEquityCall             float64
	Equity                           float64
	EquityPercentage                 float64
	LiquidationValue                 float64
	LongMarginValue                  float64
	LongOptionMarketValue            float64
	LongStockValue                   float64
	MaintenanceCall                  float64
	MaintenanceRequirement           float64
	Margin                           float64
	MarginEquity                     float64
	MoneyMarketFund                  float64
	MutualFundValue                  float64
	RegTCall                         float64
	ShortMarginValue                 float64
	ShortOptionMarketValue           float64
	ShortStockValue                  float64
	TotalCash                        float64
	IsInCall                         bool
	UnsettledCash                    float64
	PendingDeposits                  float64
	MarginBalance                    float64
	ShortBalance                     float64
	AccountValue                     float64
}

/*type CurrentBalance struct {
	AvailableFunds                   float64
	AvailableFundsNonMarginableTrade float64
	BuyingPower                      float64
	BuyingPowerNonMarginableTrade    float64
	DayTradingBuyingPower            float64
	DayTradingBuyingPowerCall        float64
	Equity                           float64
	EquityPercentage                 float64
	LongMarginValue                  float64
	MaintenanceCall                  float64
	MaintenanceRequirement           float64
	MarginBalance                    float64
	RegTCall                         float64
	ShortBalance                     float64
	ShortMarginValue                 float64
	SMA                              float64
	IsInCall                         float64
	StockBuyingPower                 float64
	OptionBuyingPower                float64
}*/

type CurrentBalance struct {
	AccruedInterest                  float64
	CashBalance                      float64
	CashReceipts                     float64
	LongOptionMarketValue            float64
	LiquidationValue                 float64
	LongMarketValue                  float64
	MoneyMarketFund                  float64
	Savings                          float64
	ShortMarketValue                 float64
	PendingDeposits                  float64
	MutualFundValues                 float64
	BondValue                        float64
	ShortOptionMarketValue           float64
	AvailableFunds                   float64
	AvailableFundsNonMarginableTrade float64
	BuyingPower                      float64
	BuyingPowerNonMarginableTrade    float64
	DayTradingBuyingPower            float64
	Equity                           float64
	EquityPercentage                 float64
	LongMarginValue                  float64
	MaintenanceCall                  float64
	MaintenanceRequirement           float64
	MarginBalance                    float64
	RegTCall                         float64
	ShortBalance                     float64
	ShortMarginValue                 float64
	SMA                              float64
}

type ProjectedBalance struct {
	AvailableFunds                   float64
	AvailableFundsNonMarginableTrade float64
	BuyingPower                      float64
	BuyingPowerNonMarginableTrade    float64
	DayTradingBuyingPower            float64
	DayTradingBuyingPowerCall        float64
	Equity                           float64
	EquityPercentage                 float64
	LongMarginValue                  float64
	MaintenanceCall                  float64
	MaintenanceRequirement           float64
	MarginBalance                    float64
	RegTCall                         float64
	ShortBalance                     float64
	ShortMarginValue                 float64
	SMA                              float64
	IsInCall                         bool
	StockBuyingPower                 float64
	OptionBuyingPower                float64
}

type AggregatedBalance struct {
	CurrentLiquidationValue float64
	LiquidationValue        float64
}

type FullOrder struct {
	Session                  string
	Duration                 string
	OrderType                string
	CancelTime               string
	ComplexOrderStrategyType string
	Quantity                 float64
	FilledQuantity           float64
	RemainingQuantity        float64
	RequestedDestination     string
	DestinationLinkName      string
	ReleaseTime              string
	StopPrice                float64
	StopPriceLinkBasis       string
	StopPriceLinkType        string
	StopPriceOffset          float64
	StopType                 string
	Price                    float64
	TaxLotMethod             string
	OrderLegCollection       []FullOrderLeg
	ActivationPrice          float64
	SpecialInstruction       string
	OrderStrategyType        string
	OrderId                  int
	Cancelable               bool
	Editable                 bool
	Status                   string
	EnteredTime              string
	CloseTime                string
	Tag                      string
	AccountNumber            int
	OrderActivityCollection  []FullOrderActivity
	ReplacingOrderCollection string
	ChildOrderStrategies     string
	StatusDescription        string
}

type FullOrderActivity struct {
	ActivityType           string
	ExecutionType          string
	Quantity               float64
	OrderRemainingQuantity float64
	ExecutionLegs          []FullExecutionLeg
}

type FullExecutionLeg struct {
	LegId             int
	Price             float64
	Quantity          float64
	MismarkedQuantity float64
	InstrumentId      int
	Time              string
}

type FullOrderLeg struct {
	OrderLegType   string
	LegId          int
	Instrument     InstrumentRef
	Instruction    string
	PositionEffect string
	Quantity       float64
	QuantityType   string
	DivCapGains    string
	ToSymbol       string
}

type SingleLegOrder struct {
	OrderType   string `default:"MARKET"`
	Session     string `default:"NORMAL"`
	Duration    string `default:"DAY"`
	Strategy    string `default:"SINGLE"`
	Instruction string
	Quantity    float32
	Instrument  SimpleOrderInstrument
}

type MultiLegOrder struct {
	OrderType          string // LIMIT, MARKET
	Session            string // NORMAL
	Duration           string // DAY
	Strategy           string // SINGLE
	OrderLegCollection []SimpleOrderLeg
}

type SimpleOrderLeg struct {
	Instruction string
	Quantity    float32
	Instrument  SimpleOrderInstrument
}

type SimpleOrderInstrument struct {
	Symbol    string
	AssetType string // EQUITY
}

type (
	SingleLegOrderComposition      func(order *SingleLegOrder)
	MultiLegSimpleOrderComposition func(order *MultiLegOrder)
)

// Create a new Market order
func CreateSingleLegOrder(opts ...SingleLegOrderComposition) *SingleLegOrder {
	order := &SingleLegOrder{OrderType: "MARKET"}
	for _, opt := range opts {
		opt(order)
	}
	return order
}

// Set SingleLegOrder.OrderType
func OrderType(t string) SingleLegOrderComposition {
	return func(order *SingleLegOrder) {
		order.OrderType = t
	}
}

// Set SingleLegOrder.Session
func Session(session string) SingleLegOrderComposition {
	return func(order *SingleLegOrder) {
		order.Session = session
	}
}

// Set SingleLegOrder.Duration
func Duration(duration string) SingleLegOrderComposition {
	return func(order *SingleLegOrder) {
		order.Duration = duration
	}
}

// Set SingleLegOrder.Strategy
func Strategy(strategy string) SingleLegOrderComposition {
	return func(order *SingleLegOrder) {
		order.Strategy = strategy
	}
}

// Set SingleLegOrder.Instruction
func Instruction(instruction string) SingleLegOrderComposition {
	return func(order *SingleLegOrder) {
		order.Instruction = instruction
	}
}

// Set SingleLegOrder.Quantity
func Quantity(quantity float32) SingleLegOrderComposition {
	return func(order *SingleLegOrder) {
		order.Quantity = quantity
	}
}

// Set SingleLegOrder.Instrument
func Instrument(instrument SimpleOrderInstrument) SingleLegOrderComposition {
	return func(order *SingleLegOrder) {
		order.Instrument = instrument
	}
}

var OrderTemplate = `
{
  "orderType": "%s",
  "session": "%s",
  "duration": "%s",
  "orderStrategyType": "%s",
  "orderLegCollection": [
    %s
  ]
}
`

var LegTemplate = `
{
  "instruction": "%s",
  "quantity": %f,
  "instrument": {
    "symbol": "%s",
    "assetType": "%s"
  }
},
`

var LegTemplateLast = `
{
  "instruction": "%s",
  "quantity": %f,
  "instrument": {
    "symbol": "%s",
    "assetType": "%s"
  }
},
`

func marshalSingleLegOrder(order *SingleLegOrder) string {
	return fmt.Sprintf(OrderTemplate, order.OrderType, order.Session, order.Duration, order.Strategy, fmt.Sprintf(LegTemplate, order.Instruction, order.Quantity, order.Instrument.Symbol, order.Instrument.AssetType))
}

func marshalMultiLegOrder(order *MultiLegOrder) string {
	var legs string
	// UNTESTED
	for i, leg := range order.OrderLegCollection {
		if i != len(order.OrderLegCollection)-1 {
			legs += fmt.Sprintf(LegTemplate, leg.Instruction, leg.Quantity, leg.Instrument.Symbol, leg.Instrument.AssetType)
		} else {
			legs += fmt.Sprintf(LegTemplateLast, leg.Instruction, leg.Quantity, leg.Instrument.Symbol, leg.Instrument.AssetType)
		}
	}
	return fmt.Sprintf(OrderTemplate)
}

// Submit a single-leg order for the specified encrypted account ID
func (agent *Agent) SubmitSingleLegOrder(hashValue string, order *SingleLegOrder) error {
	orderJson := marshalSingleLegOrder(order)
	req, err := http.NewRequest("POST", fmt.Sprintf(endpointAccountOrders, hashValue), strings.NewReader(orderJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = agent.Handler(req)
	if err != nil {
		return err
	}
	return nil
}

// Get a specific order by account number & order ID
func (agent *Agent) GetOrder(accountNumber, orderID string) (FullOrder, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(endpointAccountOrder, accountNumber, orderID), nil)
	if err != nil {
		return FullOrder{}, err
	}
	resp, err := agent.Handler(req)
	if err != nil {
		return FullOrder{}, err
	}
	var order FullOrder
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FullOrder{}, err
	}
	err = sonic.Unmarshal(body, &order)
	if err != nil {
		return FullOrder{}, err
	}
	return order, nil
}

// fromEnteredTime, toEnteredTime format:
// yyyy-MM-ddTHH:mm:ss.SSSZ
func (agent *Agent) GetAccountOrders(accountNumber, fromEnteredTime, toEnteredTime string) ([]FullOrder, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(endpointAccountOrders, accountNumber), nil)
	if err != nil {
		return []FullOrder{}, err
	}
	q := req.URL.Query()
	q.Add("fromEnteredTime", fromEnteredTime)
	q.Add("toEnteredTime", toEnteredTime)
	req.URL.RawQuery = q.Encode()
	resp, err := agent.Handler(req)
	if err != nil {
		return []FullOrder{}, err
	}
	var orders []FullOrder
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []FullOrder{}, err
	}
	err = sonic.Unmarshal(body, &orders)
	if err != nil {
		return []FullOrder{}, err
	}
	return orders, nil
}

// WIP:
// fromEnteredTime, toEnteredTime format:
// yyyy-MM-ddTHH:mm:ss.SSSZ
func (agent *Agent) GetAllOrders(fromEnteredTime, toEnteredTime string) ([]FullOrder, error) {
	req, err := http.NewRequest("GET", endpointOrders, nil)
	if err != nil {
		return []FullOrder{}, err
	}
	q := req.URL.Query()
	q.Add("fromEnteredTime", fromEnteredTime)
	q.Add("toEnteredTime", toEnteredTime)
	req.URL.RawQuery = q.Encode()
	resp, err := agent.Handler(req)
	if err != nil {
		return []FullOrder{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []FullOrder{}, err
	}
	var orders []FullOrder
	/* TODO:
	err = sonic.Unmarshal(body, &orders)
	isErrNil(err)*/
	fmt.Println(body)
	return orders, nil
}

// Get encrypted account numbers for trading
func (agent *Agent) GetAccountNumbers() ([]AccountNumbers, error) {
	req, err := http.NewRequest("GET", endpointAccountNumbers, nil)
	if err != nil {
		return []AccountNumbers{}, err
	}
	resp, err := agent.Handler(req)
	if err != nil {
		return []AccountNumbers{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []AccountNumbers{}, err
	}
	var accountNumbers []AccountNumbers
	err = sonic.Unmarshal(body, &accountNumbers)
	if err != nil {
		return []AccountNumbers{}, err
	}
	return accountNumbers, nil
}

// Get all accounts associated with the user logged in
func (agent *Agent) GetAccounts(fields ...string) ([]Account, error) {
	var fieldsRequest string = strings.Join(fields, ",")

	req, err := http.NewRequest("GET", endpointAccounts, nil)
	if err != nil {
		return []Account{}, err
	}
	q := req.URL.Query()
	q.Add("fields", fieldsRequest)
	req.URL.RawQuery = q.Encode()
	resp, err := agent.Handler(req)
	if err != nil {
		return []Account{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Account{}, err
	}
	var accounts []Account
	err = sonic.Unmarshal(body, &accounts)
	if err != nil {
		return []Account{}, err
	}
	return accounts, nil
}

// Get account by encrypted account id
func (agent *Agent) GetAccount(id string, fields ...string) (Account, error) {
	var fieldsRequest string = strings.Join(fields, ",")

	req, err := http.NewRequest("GET", fmt.Sprintf(endpointAccount, id), nil)
	if err != nil {
		return Account{}, err
	}
	q := req.URL.Query()
	q.Add("fields", fieldsRequest)
	req.URL.RawQuery = q.Encode()
	resp, err := agent.Handler(req)
	if err != nil {
		return Account{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Account{}, err
	}
	var account Account
	err = sonic.Unmarshal(body, &account)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

// Get all transactions for the user logged in
// func (agent *Agent) GetTransactions() ([]Transaction, error) {}

// Get a transaction for a specific account id
func (agent *Agent) GetTransaction(accountNumber, transactionId string) (Transaction, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(endpointTransaction, accountNumber, transactionId), nil)
	if err != nil {
		return Transaction{}, err
	}
	resp, err := agent.Handler(req)
	if err != nil {
		return Transaction{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Transaction{}, err
	}
	var transaction Transaction
	err = sonic.Unmarshal(body, &transaction)
	if err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}
