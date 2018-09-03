package page

import (
	"fmt"
	dr "github.com/fedesog/webdriver"
	cf "goBoss/config"
	"goBoss/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

var first = true

type Message struct {
	Eg        Engineer
	MsgList   map[string]map[string]string
	ReplyList map[string]bool
}

func (m *Message) Listen() {
	m.EnterMessage()
	m.Receive()

}

func (m *Message) Receive() {
	for {
		fmt.Printf("[%s]---正在获取消息列表\n", time.Now().Format("2006-01-02 15:04:05"))
		m.CheckMsgList()
		fmt.Printf("[%s]---回复列表: %+v其中value为true时代表简历已投递!\n", time.Now().Format("2006-01-02 15:04:05"), m.ReplyList)
		m.ReFetch()
	}
}

func (m *Message) SendJobMail(company, bossName, title string) {
	job_link, err := m.Eg.GetElement("消息页面", "职位链接").Attr(m.Eg.Session(), "href")
	Assert(err)
	job := utils.Spider{Url: job_link}
	desc, err := job.Html(m.Eg.GetElement("消息页面", "岗位要求").Value)
	if err != nil {
		log.Println("获取职位描述失败!Error: ", err)
	}
	addr := job.Text(m.Eg.GetElement("消息页面", "职位地址").Value)
	if desc == "" || addr == "" {
		desc = "职位可能已下架, 未获取到jd"
		addr = "职位可能已下架, 未获取到工作地点"
	}
	bs := m.Eg.ScreenAsBs64()
	d := utils.Mail{
		Subject: title,
		Content: fmt.Sprintf(`
		<html>
			<body>
				<h4>公司: %s</h4>
				<h4>名字: %s</h4>
				<h4>地址: %s</h4>
				<h4>%s</h4>
				<img src="data:image/png;base64, %s"/>
			</body>
		</html>	
		`, company, bossName, addr, desc, bs),
	}
	d.Send() // 发送邮件
}

func (m *Message) SendMsg(companyType, bossName, company string) {
	var reply string
	dialog := m.Eg.GetElement("消息页面", "消息对话框")
	switch {
	case companyType == "star":
		reply = fmt.Sprintf(cf.Config.StarReply, bossName, company)
	case companyType == "black":
		reply = cf.Config.BlackReply
	default:
		reply = cf.Config.CommonReply
	}
	err := dialog.SendKeys(m.Eg.Session(), reply)
	time.Sleep(4 * time.Second)
	Assert(err)
	err = m.Eg.GetElement("消息页面", "发送按钮").Click(m.Eg.Session())
	if err != nil {
		fmt.Printf("[%s]---自动回复失败!内容: %s, 接受者公司: %s, 接受者: %s\n Error: %s\n", time.Now().Format("2006-01-02 15:04:05"), reply, company, bossName, err.Error())
	}
	fmt.Printf("[%s]---自动回复成功!内容: %s, 接受者公司: %s, 接受者: %s\n", time.Now().Format("2006-01-02 15:04:05"), reply, company, bossName)
	m.ReplyList[fmt.Sprintf("%s|%s", company, bossName)] = false
	m.SendJobMail(company, bossName, "回复boss消息成功!")

}

func (m *Message) IsStar(company string) string {
	// 判断是否是大厂
	stars := cf.Config.StarCompany
	black_list := cf.Config.BlackList
	for _, star := range stars {
		if strings.Contains(strings.ToUpper(company), strings.ToUpper(star)) {
			return "star"
		}
	}
	for _, black := range black_list {
		if strings.Contains(strings.ToUpper(company), strings.ToUpper(black)) {
			return "black"
		}
	}
	return "common"
}

func (m *Message) SendInfo(bossName, company string) {
	err := m.Eg.GetElement("消息页面", "发送简历").Click(m.Eg.Session())
	if err != nil {
		fmt.Printf("[%s]---遇到问题: 发送简历给公司: %s Boss: %s 出错!Error: %s\n", time.Now().Format("2006-01-02 15:04:05"), company, bossName, err.Error())
	}
	time.Sleep(2 * time.Second)
	err = m.Eg.GetElement("消息页面", "发送简历确认").Click(m.Eg.Session())
	Assert(err)
	fmt.Printf("[%s]---发送简历给公司: %s Boss: %s 成功!", time.Now().Format("2006-01-02 15:04:05"), company, bossName)
	m.SendJobMail(company, bossName, "成功发送简历给Boss!")
}

func (m *Message) ReFetch() {
	//没有新消息或者没有消息
	fmt.Printf("[%s]---正在重新获取消息\n", time.Now().Format("2006-01-02 15:04:05"))
	m.Eg.Session().Refresh()
	time.Sleep(time.Duration(cf.Config.Delay) * time.Second) // 延迟Delay秒刷新
}

func (m *Message) EnterMessage() {
	time.Sleep(5 * time.Second)
	err := m.Eg.GetElement("首页", "消息").Click(m.Eg.Session())
	Assert(err)
}

func (m *Message) CheckMsgList() {
	time.Sleep(3 * time.Second)
	// 获取消息列表
	messageList, e := m.Eg.GetElement("消息页面", "消息列表").GetElements(m.Eg.Session())
	Assert(e)
	fmt.Printf("\n[%s]---正在抓取最近5条消息\n", time.Now().Format("2006-01-02 15:04:05"))
	if len(messageList) == 0 {
		// 消息列表为空， 持续检查
		log.Printf("消息列表为空, 请检查!")
	}
	for i, ms := range messageList[:5] {
		ms.Click()
		// 输出前10条最新消息
		time.Sleep(2 * time.Second)
		info := m.getInfo()
		fmt.Printf("[%s]---第%d条消息内容为: %+v\n", time.Now().Format("2006-01-02 15:04:05"), i+1, info)
		key := info["company"] + "|" + info["candidateName"]
		//if !first {
		//	// 说明不是第一次运行, 不发消息
		//	company, candidateName := info["company"], info["candidateName"]
		//	if info_bf, ok := m.MsgList[key]; ok {
		//		if info["latest"] != info_bf["latest"] && info_bf["latest"] != "" && info["latest"] != "" {
		//			m.dealMsg(company, candidateName, key, info)
		//		}
		//	} else {
		//		m.dealMsg(company, candidateName, key, info)
		//	}
		//}
		m.MsgList[key] = info
		time.Sleep(1 * time.Second)
	}
	first = false

}

func (m *Message) dealMsg(company, bossName, key string, info map[string]string) {
	star := m.IsStar(company)
	salary := strings.Split(info["money"], "-")
	var low, high string
	if len(salary) > 1 {
		low, high = strings.Replace(salary[0], "K", "", -1), strings.Replace(salary[1], "K", "", -1)
	}
	lowSalary, err := strconv.ParseInt(low, 10, 64)
	if err != nil {
		log.Println("获取最低薪水失败!Error: ", err.Error())
	}
	highSalary, err := strconv.ParseInt(high, 10, 64)
	if err != nil {
		log.Println("获取最高薪水失败!Error: ", err.Error())
	}
	// 没发送过消息才回复
	msg, _ := m.Eg.GetElement("消息页面", "自己的消息").GetElements(m.Eg.Session())
	// 如果预期薪水小于最大-1且大于最低+1, 则继续。如公司薪水为8-15k, 预期为12K, 满足要求。
	if (lowSalary+1) < cf.Config.ExpectSalary && cf.Config.ExpectSalary < (highSalary-1) {
		if status, ok := m.ReplyList[key]; !ok {
			// 发送消息
			if len(msg) == 0 {
				m.SendMsg(star, bossName, company)
			} else {
				// 你们之前沟通过
				m.SendJobMail(company, bossName, "有Boss来新消息了!(沟通过职位)")
			}

		} else {
			if star == "star" {
				// 回复包含简历且未发送过简历
				if strings.Contains(info["latest"], cf.Config.ResumeKeyword) && !status {
					// 发送简历
					m.SendInfo(bossName, company)
					m.ReplyList[key] = true
				}
			}
			// 非大厂不自动发送简历
		}
	} else {
		log.Printf("[%s]---该公司给的待遇不在考虑范围之内!\n", time.Now().Format("2006-01-02 15:04:05"))
	}

}

func (m *Message) getInfo() map[string]string {
	info := make(map[string]string)
	bossEle, err := m.Eg.GetElement("消息页面", "候选人信息").GetElements(m.Eg.Session())
	Assert(err)
	if len(bossEle) > 0 {
		info["candidateName"], _ = bossEle[0].Text()
		//info["company"], _ = bossEle[1].Text()
		//if len(bossEle) < 3 {
		//	info["bossTitle"] = "未获取到boss职位, 可能是猎头"
		//} else {
		//	info["bossTitle"], _ = bossEle[2].Text()
		//}
	}
	//jobEle, err := m.Eg.GetElement("消息页面", "职位信息").GetElements(m.Eg.Session())
	//Assert(err)
	//if len(jobEle) > 0 {
	//	info["position"], _ = jobEle[1].Text()
	//	info["money"], _ = jobEle[2].Text()
	//	info["base"], _ = jobEle[3].Text()
	//}
	for k, v := range info {
		info[k] = strings.Replace(v, " ", "", -1)
	}
	eles, _ := m.Eg.GetElement("消息页面", "聊天内容").GetElements(m.Eg.Session())
	var latest string
	if len(eles) > 0 {
		latest, _ = eles[len(eles)-1].Text()
		if latest == "" {
			// 可能是表情
			emoji, err := eles[len(eles)-1].FindElements(dr.FindElementStrategy("css selector"), "i")
			if err != nil {
				log.Printf("获取emoji表情失败!")
			}
			for _, em := range emoji {
				title, _ := em.GetAttribute("title")
				latest += title
			}
		}
	}
	info["latest"] = latest
	return info
}
