
package request
import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)
type MerchantSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
       MerchantName  *string `json:"merchantName" form:"merchantName"` 
       ContactPerson  *string `json:"contactPerson" form:"contactPerson"` 
       ContactPhone  *string `json:"contactPhone" form:"contactPhone"` 
       Address  *string `json:"address" form:"address"` 
       BusinessScope  *string `json:"businessScope" form:"businessScope"` 
       IsEnabled  *bool `json:"isEnabled" form:"isEnabled"` 
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
