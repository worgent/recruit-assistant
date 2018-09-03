package util

import "log"

func Assert(err error) {
	if err != nil {
		log.Printf("Error: %s", err.Error())
		// panic("程序遇到问题啦, 请检查截图和日志...")
	}
}
