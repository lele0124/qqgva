package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/api"

var (
	Router                    = new(router)
	apiMerchant               = api.Api.Merchant
	apiMerchantBusinessRecord = api.Api.MerchantBusinessRecord
)

type router struct {
	Merchant               merchant
	MerchantBusinessRecord merchant_record
}
