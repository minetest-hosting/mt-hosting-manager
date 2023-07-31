package wallee_test

import (
	"mt-hosting-manager/api/wallee"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://app-wallee.com/en-us/doc/api/web-service#_example
func TestCreateMac(t *testing.T) {
	userID := "2481632"
	key := "OWOMg2gnaSx1nukAM6SN2vxedfY1yLPONvcTKbhDv7I="
	method := http.MethodGet
	path := "/api/transaction/read?spaceId=12&id=1"
	ts := int64(1425387916)

	mac, err := wallee.CreateMac(userID, key, method, path, ts)
	assert.NoError(t, err)
	assert.Equal(t, "tz7RqRBH8GEw29tj2/x76opKmlbthgSFJLqiwkKIUQlKiM1HaHeisC/IQzaEgALMgwSI0kvJdPAmbS11oxzz4Q==", mac)
}
