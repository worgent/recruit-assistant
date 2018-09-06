package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	Config  UserConfig
	RConfig ResumeConfig

	RSConfig ResumeSimulateConfig
	Environ Env
)

const (
	//ChromeReg = `SOFTWARE\Google\Chrome\BLBeacon`
)

//用户配置
type UserConfig struct {
	Directory     string   `json:"directory"`
	User          string   `json:"user"`
	Password      string   `json:"password"`
	Receiver      string   `json:"receiver"`
	Sender        string   `json:"sender"`
	SenderPwd     string   `json:"sender_pwd"`
	MailPort      int      `json:"mail_port"`
	AppID         string   `json:"app_id"`
	APIKey        string   `json:"api_key"`
	SecretKey     string   `json:"secret_key"`
	BlackList     []string `json:"black_list"`
	ResumeKeyword string   `json:"resume_keyword"`
	Retry         int      `json:"retry"`
	Delay         int      `json:"delay"`
	AutoResume    bool     `json:"auto_resume"`
	AutoDownload  bool     `json:"auto_download"`
	DriverUrl     string   `json:"driver_url"`
	LoginURL      string   `json:"login_url"`
	Host          string   `json:"host"`
	LoginJSON     string   `json:"login_json"`
	MsgPage       string   `json:"msg_page"`
	JobJSON       string   `json:"job_json"`
	HisMsg        string   `json:"his_msg"`
	ResumeURL     string   `json:"resume_url"`
	WebTimeout    int      `json:"web_timeout"`
	Headless      bool     `json:"headless"`
	StarCompany   []string `json:"star_company"`
	StarReply     string   `json:"star_reply"`
	BlackReply    string   `json:"black_reply"`
	CommonReply   string   `json:"common_reply"`
	ExpectSalary  int64    `json:"expect_salary"`
	MailServer    string   `json:"mail_server"`
}

//简历招聘配置
type ResumeConfig struct {
	//数据库配置
	MysqlConnectStr string `json:"MysqlConnectStr"`

	ChromeApp	string `json:"ChromeApp"`
	//沟通上限
	CommunicateLimit int `json:"CommunicateLimit"`
	//简历筛选上限
	ResumeFilterLimit int `json:"ResumeFilterLimit"`
	//简历翻页上限
	ResumePageLimit int `json:"ResumePageLimit"`
	//学历范围
	EducationList []string `json:"EducationList"`
	//专业范围，特殊，为包含关键字，非完全匹配
	SpecialList []string `json:"SpecialList"`
	//年龄范围
	AgeLowLimit  int `json:"AgeLowLimit"`
	AgeHighLimit int `json:"AgeHighLimit"`

	//经验范围
	ExperienceLowLimit  int `json:"ExperienceLowLimit"`
	ExperienceHighLimit int `json:"ExperienceHighLimit"`

	//薪资范围,特殊，下限不高于，上限不高于；非范围内
	SalaryLowLimit  int `json:"SalaryLowLimit"`
	SalaryHighLimit int `json:"SalaryHighLimit"`

	//上线活跃时间，刚刚活跃, 今日活跃，   3日内活跃
	ActiveTimeList []string `json:"ActiveTimeList"`

	//求职状态，离职-随时到岗，在职-考虑机会，在职-月内到岗，在职-暂不考虑
	JobSeekingStatusList []string `json:"JobSeekingStatusList"`

	Be985 bool `json:"Be985"`
	Be211 bool `json:"Be211"`

	//薪酬-经验对应
	SalaryExperienceConfig map[int][]int `json:"SalaryExperienceConfig"`
}

type ResumeSimulateConfig struct{
	//默认延时
	WebOperationDelay time.Duration `json:"WebOperationDelay"`
	//登录完成后延时
	LoginDelay time.Duration `json:"LoginDelay"`
	//点击推荐牛人后延时
	RecommendDelay time.Duration `json:"RecommendDelay"`
	//翻页延时
	NextPageDelay time.Duration `json:"NextPageDelay"`
	//点击简历详情延时
	OpenResumeDetailDelay time.Duration `json:"OpenResumeDetailDelay"`
	//关闭简历详情延时
	CloseResumeDetailDelay time.Duration `json:"CloseResumeDetailDelay"`
	//读简历详情延时
	ReadBriefDelay time.Duration `json:"ReadBriefDelay"`
	//读简历详情延时
	ReadDetailDelay time.Duration `json:"ReadDetailDelay"`
	//沟通延时
	CommunicateDelay time.Duration `json:"CommunicateDelay"`



}

type Env struct {
	Root       string
	Sys        string
	DriverName string
	DriverZip  string
	QrcodeFile string
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func init() {
	Environ.Root = GetCurrentDirectory()
	// Environ.Root, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	Environ.Sys = runtime.GOOS
	Environ.QrcodeFile = "qrcode.png"

	//检查文件，如不存在
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/config/data.json", Environ.Root))
	if err != nil {
		log.Printf("打开用户配置文件失败，切换/Users/hadoop/go/src/goBoss! Error: %s", err.Error())
		Environ.Root = "/Users/hadoop/go/src/goBoss"
	}

	// 解析json
	data, err = ioutil.ReadFile(fmt.Sprintf("%s/config/data.json", Environ.Root))
	if err != nil {
		log.Panicf("打开用户配置文件失败! Error: %s", err.Error())
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Panicf("解析用户配置文件data.json失败!Error: %s", err.Error())
	}

	//解析resume-data
	data, err = ioutil.ReadFile(fmt.Sprintf("%s/config/data-resume.json", Environ.Root))
	if err != nil {
		log.Panicf("打开招聘配置文件失败! Error: %s", err.Error())
	}
	err = json.Unmarshal(data, &RConfig)
	if err != nil {
		log.Panicf("解析招聘配置文件data-resume.json失败!Error: %s", err.Error())
	}

	//解析simulate-data
	data, err = ioutil.ReadFile(fmt.Sprintf("%s/config/data-simulate.json", Environ.Root))
	if err != nil {
		log.Panicf("打开模拟配置文件失败! Error: %s", err.Error())
	}
	err = json.Unmarshal(data, &RSConfig)
	if err != nil {
		log.Panicf("解析模拟配置文件data-simulate.json失败!Error: %s", err.Error())
	}
}
