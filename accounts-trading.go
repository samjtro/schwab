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
	accountEndpoint        string = "https://api.schwabapi.com/trader/v1"
	endpointAccountNumbers string = accountEndpoint + "/accounts/accountNumbers"
	endpointAccounts       string = accountEndpoint + "/accounts"
	endpointAccount        string = accountEndpoint + "/accounts/%s"
	// endpointUserPreference string = accountEndpoint + "/userPreference"
	endpointOrders        string = accountEndpoint + "/orders"
	endpointAccountOrders string = accountEndpoint + "/accounts/%s/orders"
	endpointAccountOrder  string = accountEndpoint + "/accounts/%s/orders/%s"
	// endpointPreviewOrder  string = accountEndpoint + "/accounts/%s/previewOrder"
	// endpointTransactions string = accountEndpoint + "/accounts/%s/transactions"
	endpointTransaction string = accountEndpoint + "/accounts/%s/transactions/%s"
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
		if len(order.OrderLegCollection) == 0 {
			order.OrderLegCollection = append(order.OrderLegCollection, OrderLeg{})
		}
		order.OrderLegCollection[0].Instruction = instruction
	}
}

// Set SingleLegOrder.Quantity
func Quantity(quantity int) SingleLegOrderComposition {
	return func(order *SingleLegOrder) {
		if len(order.OrderLegCollection) == 0 {
			order.OrderLegCollection = append(order.OrderLegCollection, OrderLeg{})
		}
		order.OrderLegCollection[0].Quantity = quantity
	}
}

// Set SingleLegOrder.Instrument
func Instrument(instrument SimpleOrderInstrument) SingleLegOrderComposition {
	return func(order *SingleLegOrder) {
		if len(order.OrderLegCollection) == 0 {
			order.OrderLegCollection = append(order.OrderLegCollection, OrderLeg{})
		}
		order.OrderLegCollection[0].Instrument = instrument
	}
}

// Submit a single-leg order for the specified encrypted account ID
func (agent *Agent) SubmitSingleLegOrder(hashValue string, order *SingleLegOrder) error {
	orderJson, err := sonic.Marshal(order)
	if err != nil {
		return err
	}
	fmt.Println(string(orderJson))
	req, err := http.NewRequest("POST", fmt.Sprintf(endpointAccountOrders, hashValue), strings.NewReader(string(orderJson)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := agent.Handler(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body")
	}
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
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

// WIP: do not use
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
	err = sonic.Unmarshal(body, &orders)
	if err != nil {
		fmt.Println(body)
		return []FullOrder{}, err
	}
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
