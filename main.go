package main

import (
	"fmt"
	"github.com/fedesog/webdriver"
	cf "goBoss/config"
	"goBoss/driver"
	"goBoss/page"
	"log"
	"os"
	//"goBoss/module/entity"
	dao "goBoss/module/dao"
	"goBoss/module/entity"
	"goBoss/module/util"
	"time"
)

var logUtil = util.MyLog{}

func main() {
	//collegeTest()
	//return
	setLog()
	driver.SetDriver() // 自动获取浏览器驱动
	chromeDriver := webdriver.NewChromeDriver(fmt.Sprintf("%s/driver/%s", cf.Environ.Root, cf.Environ.DriverName))
	engine := &page.Engine{Dr: chromeDriver}
	lg := &page.Login{Eg: engine}
	lg.Eg.Start()
	lg.Eg.OpenBrowser()
	//测试先关掉登录
	lg.Login()

	candidateMap := make(map[string]*entity.Candidate)
	candidateList := []string{}
	resume := &page.Resume{
		Eg: engine, CandidateList: candidateList,
		CandidateMap: candidateMap,
	}
	resume.Run()

	//time.Sleep(10 * time.Second)

	defer page.TearDown(engine)
}

func setLog() {
	//set logfile Stdout
	logFile, logErr := os.OpenFile(fmt.Sprintf("%s/boss-%s.log",
		cf.Environ.Root, time.Now().Format("2006-01-02-15-04-05")),
		os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", logFile, "cServer start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func test() {
	str := "1撒zxz是谁我我说-22_-laoYu#$@sd兰考县"
	r := []rune(str)
	//fmt.Println("rune=", r)
	strSlice := []string{}
	cnstr := ""
	for i := 0; i < len(r); i++ {
		if r[i] <= 40869 && r[i] >= 19968 {
			cnstr = cnstr + string(r[i])
			strSlice = append(strSlice, cnstr)

		}
		//fmt.Println("r[", i, "]=", r[i], "string=", string(r[i]))
	}
	if 0 == len(strSlice) {
		//无中文，需要跳过，后面再找规律
	}
	fmt.Println("原字符串:", str, "    提取出的中文字符串:", cnstr)
	fmt.Println(strSlice)

}
func collegeTest() {
	collegeDao := dao.CollegeDao{}

	college, e := collegeDao.FindCollege("北京大")

	util.Assert(e)

	logUtil.Debug(college.Dump())

}
