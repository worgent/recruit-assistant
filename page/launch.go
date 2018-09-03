package page

import (
	"bytes"
	"fmt"
	cf "goBoss/config"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Login struct {
	Eg Engineer
}

func Assert(err error) {
	if err != nil {
		log.Printf("Error: %s", err.Error())
		// panic("程序遇到问题啦, 请检查截图和日志...")
	}
}

//func (w *Login) sendCode() {
//	// 识别验证码
//	for {
//		// image, err := w.Session.FindElement(lg["验证码"].Method, lg["验证码"].Value)
//		image := w.Eg.GetElement("登录页面", "验证码")
//		src, err := image.Attr(w.Eg.Session(), "src")
//		Assert(err)
//		// src, _ := image.GetAttribute("src")
//		code := verify.GetCode(src)
//		if len(code) != 4 {
//			// 验证码识别有误
//			time.Sleep(3 * time.Second / 2)
//			fmt.Println("验证码长度不为4, 重新获取!")
//			image.Click(w.Eg.Session())
//			continue
//		} else {
//			err = w.Eg.GetElement("登录页面", "验证码输入框").SendKeys(w.Eg.Session(), code)
//			Assert(err)
//
//			err = w.Eg.GetElement("登录页面", "登录").Click(w.Eg.Session())
//			Assert(err)
//			time.Sleep(3 * time.Second / 2)
//			text, _ := w.Eg.GetElement("登录页面", "验证码错误").Text(w.Eg.Session())
//			// Assert(err)
//			if text == "" {
//				// 登录成功, break
//				fmt.Println("恭喜您登录成功...")
//				break
//			} else {
//				fmt.Println("验证码错误, 重新登录...")
//				time.Sleep(3 * time.Second / 2)
//				w.Eg.GetElement("登录页面", "验证码").Click(w.Eg.Session())
//				continue
//			}
//		}
//
//	}
//}

func (w *Login) Login() {

	w.GeneratePic()
	w.OpenQrCode()
	fmt.Printf("[%s]---请在8秒内扫码登录\n", time.Now().Format("2006-01-02 15:04:05"))
	time.Sleep(10 * time.Second)
	url, _ := w.Eg.GetUrl()
	if strings.Contains(url, "login") {
		// 还未登陆成功
		fmt.Printf("[%s]---10秒内未成功登录, 请在8秒内扫码登录\n", time.Now().Format("2006-01-02 15:04:05"))
		w.Login()
	} else {
		fmt.Printf("[%s]---恭喜您登录成功!\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	//err = w.Eg.GetElement("登录页面", "用户名输入框").SendKeys(w.Eg.Session(), cf.Config.User)
	//Assert(err)
	//err = w.Eg.GetElement("登录页面", "密码输入框").SendKeys(w.Eg.Session(), cf.Config.Password)
	//Assert(err)
	//w.sendCode()
}

func (w *Login) GeneratePic() {
	w.Eg.Session().Refresh() // 刷新页面
	// 进入密码登录页面
	err := w.Eg.GetElement("登录页面", "二维码登录").Click(w.Eg.Session())
	Assert(err)
	time.Sleep(1 * time.Second)
	//pic_url, err := w.Eg.GetElement("登录页面", "验证码图片").Attr(w.Eg.Session(), "src")
	//Assert(err)
	//if pic_url == "" {
	//	// 刷新验证码
	//	w.GeneratePic()
	//}
	//req := utils.Request{
	//	Url:    pic_url,
	//	Method: "GET",
	//}
	//result := req.Http()
	//if !result["status"].(bool) {
	//	log.Fatal("获取Boss直聘验证码失败!")
	//}
	//bs := result["result"].([]byte)   // 转换为[]byte
	bs, err := w.Eg.Screen()
	if err != nil {
		log.Fatal("获取二维码图片失败!")
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/picture/%s", cf.Environ.Root, cf.Environ.QrcodeFile), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	f.Write(bs)
	defer f.Close()
}

func (w *Login) OpenQrCode() {
	var cmd *exec.Cmd
	filePath := fmt.Sprintf("%s/picture/%s", cf.Environ.Root, cf.Environ.QrcodeFile)
	switch cf.Environ.Sys {
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", filePath)
	default:
		cmd = exec.Command("sh", "-c", fmt.Sprintf("open %s", filePath))
	}
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		log.Fatal("打开二维码图片失败! Error: ", err)
	}
}
