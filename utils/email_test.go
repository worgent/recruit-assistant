package utils

import (
	"fmt"
	"testing"
)

func TestMail_Send(t *testing.T) {
	d := &Mail{
		Subject: "自动回复Boss消息成功!",
		Attach:  `C:\Users\Woody\go\src\goBoss\picture\2018_06_07_00_03_07_error.png`,
		Content: fmt.Sprintf(`<h4>内容: %s, 接受者公司: %s, 接受者: %s</h4>`, "我想进你你你你", "111", "222"),
	}
	d.Send()
}
