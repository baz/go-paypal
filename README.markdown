go-paypal
---

go-paypal is a package written in Go for accessing PayPal APIs using the ["NVP"](https://cms.paypal.com/us/cgi-bin/?cmd=_render-content&content_ID=developer/e_howto_api_nvp_NVPAPIOverview#id09C2F0G0C7U) format.

Included is a method for using the [Digital Goods for Express Checkout](https://cms.paypal.com/us/cgi-bin/?cmd=_render-content&content_ID=developer/e_howto_api_IntegratingExpressCheckoutDG) payment option.

Quick Start
---
	import "paypal"
	
	client := paypal.NewClient(username, password, signature, true)
	
	goods := make([]paypal.PayPalDigitalGood, 1)
	good := new(paypal.PayPalDigitalGood)
	good.Name, good.Amount, good.Quantity = "Test Good", paymentAmount, 1
	goods[0] = *good
	
	response, _ := client.SetExpressCheckoutDigitalGoods(paymentAmount, currencyCode, returnURL, cancelURL, goods)
	
	fmt.Println(response.Values)
