package paypal

import (
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
)

const (
	sandboxEndpoint					= "https://api-3t.sandbox.paypal.com/nvp"
	productionEndpoint				= "https://api-3t.paypal.com/nvp"
	version							= "84"
)

type PayPalClient struct {
	username				string
	password				string
	signature				string
	usesSandbox				bool
	client					*http.Client
}

type PayPalDigitalGood struct {
	Name			string
	Amount			float64
	Quantity		int16
}

type PayPalResponse struct {
	Ack				string
	CorrelationId	string
	Timestamp		string
	Version			string
	Build			string
	Values			url.Values
}

func NewClient(username, password, signature string, usesSandbox bool) *PayPalClient {
	return &PayPalClient{username, password, signature, usesSandbox, new(http.Client)}
}

func (pClient *PayPalClient) SetExpressCheckoutDigitalGoods(paymentAmount float64, currencyCode string, returnURL, cancelURL string, goods[]PayPalDigitalGood) (*PayPalResponse, error) {
	values := url.Values{}
	values.Set("METHOD", "SetExpressCheckout")
	values.Add("PAYMENTREQUEST_0_AMT", fmt.Sprintf("%.2f", paymentAmount))
	values.Add("PAYMENTREQUEST_0_PAYMENTACTION", "Sale");
	values.Add("PAYMENTREQUEST_0_CURRENCYCODE", currencyCode);
	values.Add("RETURNURL", returnURL);
	values.Add("CANCELURL", cancelURL);
	values.Add("REQCONFIRMSHIPPING", "0");
	values.Add("NOSHIPPING", "1");
	values.Add("SOLUTIONTYPE", "Sole");

	for i := 0; i < len(goods); i++ {
		good := goods[i]

		values.Add(fmt.Sprintf("%s%d", "L_PAYMENTREQUEST_0_NAME", i), good.Name)
		values.Add(fmt.Sprintf("%s%d", "L_PAYMENTREQUEST_0_AMT", i), fmt.Sprintf("%.2f", good.Amount))
		values.Add(fmt.Sprintf("%s%d", "L_PAYMENTREQUEST_0_QTY", i), fmt.Sprintf("%d", good.Quantity))
		values.Add(fmt.Sprintf("%s%d", "L_PAYMENTREQUEST_0_ITEMCATEGORY", i), "Digital")
	}

	return pClient.PerformRequest(values)
}

func (pClient *PayPalClient) PerformRequest(values url.Values) (*PayPalResponse, error) {
	values.Add("USER", pClient.username);
	values.Add("PWD", pClient.password);
	values.Add("SIGNATURE", pClient.signature);
	values.Add("VERSION", version);

	endpoint := productionEndpoint
	if pClient.usesSandbox {
		endpoint = sandboxEndpoint
	}

	formResponse, err := pClient.client.PostForm(endpoint, values)
	defer formResponse.Body.Close()
	if err != nil { return nil, err }

	body, err := ioutil.ReadAll(formResponse.Body)
	if err != nil { return nil, err }

	responseValues, err := url.ParseQuery(string(body))
	response := new(PayPalResponse)
	if err == nil {
		response.Ack = responseValues.Get("ACK")
		response.CorrelationId = responseValues.Get("CORRELATIONID")
		response.Timestamp = responseValues.Get("TIMESTAMP")
		response.Version = responseValues.Get("VERSION")
		response.Build = responseValues.Get("2975009")
		response.Values = responseValues
	}

	return response, err
}
