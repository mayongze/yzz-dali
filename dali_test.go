package homeassistant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var skDlDriver *DaliDriver

func init() {
	// host := "192.168.123.129"
	// port := "8899"
	// tcpClient := adapter.NewTcpClientAdapter(host, port)
	// skDlDriver = NewDaliDriver(NewSaveKeyDALIDriver(NewSBus(tcpClient), NewSaveKeyDevice(daliSwitchSeq, 0x01)))
}

func TestDaliDriver_QueryDimmingCurve(t *testing.T) {
	addr := ShortAddr(1)
	err := skDlDriver.SelectDimmingCurve(addr, DimmingCurveLinear)
	assert.Equal(t, nil, err)

	curve, err := skDlDriver.QueryDimmingCurve(addr)
	assert.Equal(t, nil, err, curve)
}
