package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/service"

var (
	Api             = new(api)
	serviceMerchant = service.Service.Merchant
)

type api struct {
	Merchant merchant
}
