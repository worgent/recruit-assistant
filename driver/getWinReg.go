// 感谢Wusuluren的帮助, 使用windows命令行直接操作注册表

package driver

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

//package main
//
//import (
//	"fmt"
//	"log"
//	"github.com/golang/sys/windows/registry"
//	"goBoss/config"
//	"strings"
//)
//
//func main() {
//	k, err := registry.OpenKey(registry.CURRENT_USER, config.ChromeReg, registry.ALL_ACCESS)
//	if err != nil {
//		log.Fatal("获取Windows Chrome版本失败!请检查Chrome是否安装 Error: ", err)
//	}
//	s, _, err := k.GetStringValue("version")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer k.Close()
//	verList := strings.Split(s, ".")
//	ver := verList[0]
//	fmt.Println(ver)
//}

func getWinChromeVersion() string {
	cmd := exec.Command("reg", "query", `HKEY_CURRENT_USER\Software\Google\Chrome\BLBeacon`, "/v", "version")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		log.Fatal("获取Windows Chrome版本失败!请检查Chrome是否安装 Error: ", err)
	}
	ver := string(cmdOutput.Bytes())
	verList := strings.Split(ver, " ")
	ver = verList[len(verList)-1]
	//fmt.Println(ver)
	verList = strings.Split(ver, ".")
	ver = verList[0]
	//fmt.Println(verList[0])
	return ver
}
