
package request
import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)
type MerchantBusinessRecordSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
       MerchantID  *string `json:"merchantId" form:"merchantId"` 
       RecordType  *string `json:"recordType" form:"recordType"` 
       Amount  *float64 `json:"amount" form:"amount"` 
       Description  *string `json:"description" form:"description"` 
       RecordTimeRange  []time.Time  `json:"recordTimeRange" form:"recordTimeRange[]"`
    request.PageInfo
}
