package service

var Service = new(Services)

type Services struct {
	Merchant *merchant
}

func init() {
	Service.Merchant = Merchant
}
