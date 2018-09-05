package util

import (

	"time"
	"fmt"
)

/**
	模拟相关
 */
type SimulateUtil struct {
}

//推荐筛选
func (s *SimulateUtil) Delay(delay time.Duration) {
	time.Sleep(delay * time.Second)
	logUtil.Info(fmt.Sprintf("延时[%d]s", delay))
}


