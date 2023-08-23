package wallee

import (
	"fmt"
	"net/http"
)

func (c *WalleeClient) CreateRefund(rr *CreateRefundRequest) (*CreateRefundResponse, error) {
	path := fmt.Sprintf("/api/refund/refund?spaceId=%s", c.SpaceID)
	rrsp := &CreateRefundResponse{}
	err := c.request(path, http.MethodPost, rr, &rrsp)
	return rrsp, err
}
