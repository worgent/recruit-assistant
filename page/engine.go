package page

import (
	"github.com/fedesog/webdriver"
	//"github.com/tebeka/selenium"
	"encoding/base64"
	"errors"
	"fmt"
	cf "goBoss/config"
	"goBoss/utils"
	"log"
	"os"
	"time"
)

type Engine struct {
	Dr *webdriver.ChromeDriver
	Ss *webdriver.Session
}

type Engineer interface {
	MaxWindow() error
	Start()
	OpenBrowser()
	SetWindow(width, height int) error
	GetElement(root, name string) *Element
	Driver() *webdriver.ChromeDriver
	Session() *webdriver.Session
	Screen() ([]byte, error)
	GetUrl() (string, error)
	ScreenShot(key string) string
	ScreenAsBs64() string
}

func (w *Engine) Session() *webdriver.Session {
	return w.Ss
}

func (w *Engine) Driver() *webdriver.ChromeDriver {
	return w.Dr
}

func (w *Engine) GetElement(root, name string) *Element {
	ele, ok := Page[root][name]
	if !ok {
		log.Panicf("page/element.json未找到root: [%s] key: [%s]", root, name)
	}
	return &ele
}

func (w *Engine) MaxWindow() error {
	p := fmt.Sprintf(`{"windowHandle": "current", "sessionId": "%s"}`, w.Ss.Id)
	req := utils.RequestData{
		JSON: p,
	}
	url := fmt.Sprintf("http://127.0.0.1:%d/session/%s/window/current/maximize", w.Dr.Port, w.Ss.Id)
	res, err := utils.HttpPost(url, req)
	if err != nil {
		log.Printf("response: %+v", res)
		return err
	}
	return nil
}

func (w *Engine) Start() {
	var err error
	w.Dr.Start()
	args := make([]string, 0)
	if cf.Config.Headless {
		args = append(args, "--headless")
	}
	desired := webdriver.Capabilities{
		"Platform":           "Mac",
		"goog:chromeOptions": map[string][]string{"args": args, "extensions": []string{}},
		"browserName":        "chrome",
		"version":            "",
		"platform":           "ANY",
	}
	required := webdriver.Capabilities{}
	w.Ss, err = w.Dr.NewSession(desired, required)
	if err != nil {
		log.Printf("open browser failed: %s", err.Error())
	}

}

func (w *Engine) OpenBrowser() {
	//w.Ss.Url("http://www.baidu.com")
	w.Ss.Url(cf.Config.LoginURL)
	err := w.SetWindow(900, 600)
	if err != nil {
		log.Panicf("最大化浏览器失败!!!Msg: %s", err.Error())
	}
	w.Ss.SetTimeoutsImplicitWait(cf.Config.WebTimeout)
}

func (w *Engine) SetWindow(width, height int) error {
	p := fmt.Sprintf(`{"windowHandle": "current", "sessionId": "%s", "height": %d, "width": %d}`, w.Ss.Id, height, width)
	url := fmt.Sprintf("http://127.0.0.1:%d/session/%s/window/current/size", w.Dr.Port, w.Ss.Id)
	res, err := utils.HttpPost(url, utils.RequestData{JSON: p})
	if err != nil {
		return errors.New(fmt.Sprintf(`设置浏览器窗口失败, 请检查!%+v`, res.Error))
	}
	return nil
}

func (w *Engine) Close() {
	w.Ss.CloseCurrentWindow()
	w.Dr.Stop()

}

func (w *Engine) Screen() ([]byte, error) {
	return w.Ss.Screenshot()
}

func (w *Engine) GetUrl() (string, error) {
	return w.Ss.GetUrl()
}

func (w *Engine) ScreenShot(key string) string {
	pic, _ := w.Screen()
	filename := fmt.Sprintf("%s_%s.png", key, time.Now().Format("2006_01_02_15_04_05"))
	filename = fmt.Sprintf("%s/picture/%s", cf.Environ.Root, filename)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Printf("发送消息后截图失败!Error: %s", err.Error())
	}
	f.Write(pic)
	defer f.Close()
	return filename
}

func (w *Engine) ScreenAsBs64() string {
	bt, err := w.Screen()
	if err != nil {
		log.Println("截图出错!Error: ", err.Error())
	}
	bs64 := utils.Encode(base64.StdEncoding, bt)
	return bs64
}
