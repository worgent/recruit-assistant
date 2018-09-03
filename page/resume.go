package page

/**
处理推荐简历信息
1.筛选合适简历，并进行沟通；
2.统计简历重复率
3.统计简历质量

*/

import (
	"errors"
	"fmt"
	"github.com/fedesog/webdriver"
	cf "goBoss/config"
	"goBoss/module/entity"
	"goBoss/module/util"
	"log"
	"strings"
	"time"
)

var resumeFirst = true
var candidateUtil = util.CandidateUtil{}
var logUtil = util.MyLog{}
var communicateLimit = cf.RConfig.CommunicateLimit
var resumePageLimit = cf.RConfig.ResumePageLimit
var candidateLimit = cf.RConfig.ResumeFilterLimit

type Resume struct {
	Eg Engineer
	//ResumeList   map[string]map[string]string
	//ReplyList map[string]bool
	CandidateMap map[string]*entity.Candidate
	//所有候选人ID数组
	CandidateList []string
	//沟通候选人ID数组
	CommunicateCandidateList []string
}

//处理入口
func (m *Resume) Run() {
	time.Sleep(cf.RConfig.WebOperationInterval * time.Second)
	m.EnterResume()
	time.Sleep(cf.RConfig.WebOperationInterval * time.Second)
	m.DealResumes()
	//todo 增加数据库导入
}

func (m *Resume) EnterResume() {
	//
	err := m.Eg.GetElement("首页", "推荐牛人").Click(m.Eg.Session())
	Assert(err)
	//
}

//处理简历，沟通达到指定数目停止；或者翻页达到指定次数；或者推荐简历数目达到指定书目
func (m *Resume) DealResumes() {
	dealCandidateNum := 0
	dealPageCount := 1

	logUtil.Info(fmt.Sprintf("开始处理候选人简历"))

	for {
		m.DealOnePage(&dealCandidateNum)
		m.NextPage(&dealPageCount)
		time.Sleep(cf.RConfig.WebOperationInterval * time.Second)

		//沟通超过10个人，停止
		if len(m.CommunicateCandidateList) >= communicateLimit {
			logUtil.Info(fmt.Sprintf("简历处理停止，达到沟通上限%d，本次沟通人次%d个\n",
				communicateLimit, len(m.CommunicateCandidateList)))
			break
		}
		if dealCandidateNum >= candidateLimit {
			logUtil.Info(fmt.Sprintf("简历处理停止，达到简历处理上限%d，本次处理人次%d个\n",
				candidateLimit, dealCandidateNum))

			break
		}
		if dealPageCount >= candidateLimit {
			logUtil.Info(fmt.Sprintf("简历处理停止，达到简历处理翻页上限%d，本次处理页面%d个\n",
				resumePageLimit, dealPageCount))

			break
		}

		logUtil.Info(fmt.Sprintf("候选人简历第%d页, 目前已处理%d个简历，沟通%d人\n",
			dealPageCount, dealCandidateNum, len(m.CommunicateCandidateList)))

		dealPageCount += 1
	}

	logUtil.Info(fmt.Sprintf("候选人简历处理结束\n"))
	logUtil.Info(fmt.Sprintf("本次处理页面%d个， 筛选简历%d个， 沟通人次%d个\n",
		dealPageCount, dealCandidateNum, len(m.CommunicateCandidateList)))
	for i, UID := range m.CommunicateCandidateList {
		logUtil.Info(fmt.Sprintf("index-%d 沟通人信息%s\n",
			i+1, (m.CandidateMap[UID]).BriefDump()))
	}
}

func (m *Resume) NextPage(dealCount *int) {
	//没有新消息或者没有消息
	logUtil.Info(fmt.Sprintf("候选人列表翻页至[%d]\n", (*dealCount + 1)))
	//m.Eg.Session().Refresh()
	var args = []interface{}{}
	//每次增加高度5000
	script := fmt.Sprintf("$('#container').scrollTop(%d)", (*dealCount+1)*5000)
	ret, e := m.Eg.Session().ExecuteScript(script, args)
	logUtil.Debug(fmt.Sprintf("ExecuteScript %s, ret %s\n", script, string(ret)))
	Assert(e)
}

//检查候选人列表，并根据情况发起沟通
func (m *Resume) DealOnePage(dealCandidateNum *int) {
	// 获取候选人列表
	ResumeList, e := m.Eg.GetElement("推荐牛人", "候选人列表").GetElements(m.Eg.Session())
	Assert(e)
	logUtil.Info(fmt.Sprintf("本页抓取候选人[%d]条\n",
		len(ResumeList)-*dealCandidateNum))

	if len(ResumeList)-*dealCandidateNum == 0 {
		// 消息列表为空， 持续检查
		logUtil.Info(fmt.Sprintf("本页候选人列表为空, 请检查!"))
		return
	}

	//处理新增候选人；
	//解析候选人数据
	for i, ms := range ResumeList {
		if i < *dealCandidateNum {
			continue
		}
		info, e := m.getBaseInfo(ms)
		if e != nil {
			logUtil.Warn(fmt.Sprintf("候选人信息获取失败, %s，%s", e.Error(), info.BriefDump()))
			continue
		}
		//先判断原来是不是有，
		_, exist := m.CandidateMap[info.UID]
		if exist == true {
			logUtil.Info(fmt.Sprintf("候选人信息已经存在: %s\n", i+1, info.BriefDump()))
			continue
		}
		logUtil.Debug(fmt.Sprintf("第[%d]条候选人信息为: %s\n", i+1, info.BriefDump()))

		m.CandidateList = append(m.CandidateList, info.UID)
		m.CandidateMap[info.UID] = &info

		//已沟通的标记为通过
		if info.Communicatetatus == "继续沟通" {
			info.DealResult = 0
			logUtil.Info(fmt.Sprintf("候选人已经沟通过: %s\n", i+1, info.BriefDump()))
			continue
		}

		//初步筛选候选人
		info.DealResult = candidateUtil.FilterCandidate(info)
		//只要0-为通过
		if info.DealResult != 0 {
			logUtil.Info(fmt.Sprintf("候选人不合符要求: %s\n", info.BriefDump()))
			continue
		}

		//通过筛选，点击详情获取信息
		ms.Click()
		time.Sleep(cf.RConfig.WebOperationInterval * time.Second)
		detailElement, e := m.Eg.GetElement("推荐牛人", "详情页面").GetEle(m.Eg.Session())
		Assert(e)
		info.Content, _ = detailElement.Text()
		e = m.getDetailInfo(detailElement, &info)
		if e != nil {
			logUtil.Error(fmt.Sprintf("候选人获得详细信息出错: %s\n", info.BriefDump()))
			continue
		}
		logUtil.Debug(fmt.Sprintf("第[%d]条候选人更新教育详细信息为: %s\n", i+1, info.EducationMapDump()))

		//关闭详情
		clickElement, e := m.Eg.GetElement("推荐牛人", "详情关闭按钮").GetEle(m.Eg.Session())
		Assert(e)
		clickElement.Click()
		time.Sleep(cf.RConfig.WebOperationInterval * time.Second)

		//进行二次筛选，主要筛选毕业时间，工作经历和薪资的关系
		info.DealResult = candidateUtil.FilterCandidateDetail(info)
		//只要0-为通过
		if info.DealResult != 0 {
			logUtil.Info(fmt.Sprintf("候选人二面不合符要求: %s\n", info.EducationMapDump()))
			continue
		}

		//通过筛选，进行沟通
		e = m.communicateCandidate(ms)
		if e != nil {
			logUtil.Error(fmt.Sprintf("沟通候选人%s-%s失败，原因%s\n",
				info.UID, info.Name, e.Error()))
			continue
		}

		m.CommunicateCandidateList = append(m.CommunicateCandidateList, info.UID)
		logUtil.Info(fmt.Sprintf("沟通候选人%s-%s成功\n",
			info.UID, info.Name))

		time.Sleep(cf.RConfig.WebOperationInterval * time.Second)
	}

	*dealCandidateNum = len(ResumeList)
}

//模拟点击打招呼
func (m *Resume) communicateCandidate(element webdriver.WebElement) error {
	ele, e := element.FindElement(webdriver.CSS_Selector, ".sider-op .btn-greet")
	if e == nil {
		ele.Click()
		return nil
	} else {
		return e
	}
}

//获取元素指定样式的文本
func (m *Resume) getElementText(element webdriver.WebElement, css string) (string, error) {
	ele, e := element.FindElement(webdriver.CSS_Selector, css)
	if e != nil {
		return "", e
	}
	return ele.Text()
}

//获取一组元素的文本
func (m *Resume) getElementTexts(element webdriver.WebElement, css string) ([]string, error) {
	eleTexts := []string{}
	eles, e := element.FindElements(webdriver.CSS_Selector, css)
	if e != nil {
		return eleTexts, e
	}
	for _, ele := range eles {
		eleText, _ := ele.Text()
		eleTexts = append(eleTexts, eleText)
	}
	return eleTexts, nil
}

//是否包含相应元素
func (m *Resume) containsElement(element webdriver.WebElement, css string) bool {
	_, e := element.FindElement(webdriver.CSS_Selector, css)
	if e == nil {
		return true
	} else {
		return false
	}
}

func (m *Resume) getBaseInfo(element webdriver.WebElement) (entity.Candidate, error) {
	info := entity.Candidate{}
	info.CreatedTime = time.Now().Unix()
	info.UID, _ = element.GetAttribute("data-uid")
	eles, _ := element.FindElements(webdriver.CSS_Selector, ".info-labels .label-text")
	if len(eles) != 6 {
		log.Println(".info-labels .label-text should be 6, %d", len(eles))
		return info, errors.New(".info-labels .label-text should be 6")
	}
	info.District, _ = eles[0].Text()
	info.ExperienceYearStr, _ = eles[1].Text()
	info.ExperienceYear = candidateUtil.GetExperienceYear(info.ExperienceYearStr)

	info.Education, _ = eles[2].Text()
	info.AgeStr, _ = eles[3].Text()
	info.Age = candidateUtil.GetAgeYear(info.AgeStr)

	info.JobSeekingStatus, _ = eles[4].Text()
	info.ActiveTime, _ = eles[5].Text()

	info.Name, _ = m.getElementText(element, ".geek-name")
	info.Advantage, _ = m.getElementText(element, ".advantage")
	info.RecommendReason, _ = m.getElementText(element, ".recommend-reason")
	info.JobHistory, _ = m.getElementText(element, ".experience")
	info.ExpectSalaryStr, _ = m.getElementText(element, ".badge-salary")
	info.ExpectLow, info.ExpectHigh = candidateUtil.GetExpectSalary(info.ExpectSalaryStr)

	//性别 fz-male, fz-female
	if m.containsElement(element, "fz-male") {
		info.Sex = "男"
	} else if m.containsElement(element, "fz-female") {
		info.Sex = "女"
	} else {
		info.Sex = "未知"
		logUtil.Warn(fmt.Sprintf("性别未知，info：%s", info.Dump()))

	}

	eles, _ = element.FindElements(webdriver.CSS_Selector, ".chat-info .text p")
	if len(eles) != 3 {
		logUtil.Error(fmt.Sprintf("%s, .chat-info .text p should be 3,info:%s", info.Dump()))
		return info, errors.New(".chat-info .text p should be 3")
	}
	tempStr, _ := eles[2].Text()
	tempStrs := strings.Split(tempStr, "•")
	if len(tempStrs) != 2 {
		logUtil.Error(fmt.Sprintf(".chat-info .text p 3 should be split to 2 with •,info:%s", info.Dump()))
		return info, errors.New(".chat-info .text p 3 should be split to 2 with •")
	}

	//获得教育信息
	info.School = tempStrs[0]
	info.Special = tempStrs[1]

	e := errors.New("")
	info.Communicatetatus, e = m.getElementText(element, ".sider-op .btn-greet")
	if e != nil {
		info.Communicatetatus, e = m.getElementText(element, ".sider-op .btn-continue")
		if e != nil {
			logUtil.Error(fmt.Sprintf("Communicatetatus not found,info:%s", info.Dump()))
			return info, errors.New("Communicatetatus not found")
		}
	}

	//for k, v := range info {
	//	info[k] = strings.Replace(v, " ", "", -1)
	//}

	return info, nil
}

func (m *Resume) getDetailInfo(element webdriver.WebElement, info *entity.Candidate) error {
	resumeItemEles, _ := element.FindElements(webdriver.CSS_Selector,
		".dialog-resume-full .resume-dialog .resume-item")
	if len(resumeItemEles) <= 0 {
		logUtil.Error(fmt.Sprintf("resumeItem 应该大于0，info:%s",
			info.Dump()))
		return errors.New("resumeItem 应该大于0")
	}
	for i, resumeItemEle := range resumeItemEles {
		if i == 0 {
			continue
		}
		titleEle, e := resumeItemEle.FindElement(webdriver.CSS_Selector,
			".title")
		resumeItemStr, _ := resumeItemEle.Text()
		if e != nil {
			logUtil.Warn(fmt.Sprintf("resumeitem获取title失败，content：%s", resumeItemStr))
			continue
		}
		titleStr, _ := titleEle.Text()
		logUtil.Info(fmt.Sprintf("开始处理[%s]信息", titleStr))

		if titleStr == "教育经历" {
			historyItemEles, _ := resumeItemEle.FindElements(webdriver.CSS_Selector,
				".history-list .history-item")
			if len(historyItemEles) <= 0 {
				logUtil.Error(fmt.Sprintf("historyItemEles 应该大于0，info:%s, css :.history-list .history-item",
					info.Dump()))
				return errors.New("historyItemEles 应该大于0")
			}
			educationMap := []entity.Study{}
			for _, historyItemEle := range historyItemEles {
				study := entity.Study{}
				study.Content, _ = historyItemEle.Text()
				periodEle, e := historyItemEle.FindElement(webdriver.CSS_Selector,
					".period")
				Assert(e)
				periodStr, _ := periodEle.Text()
				study.InTime, study.OutTime = candidateUtil.GetEducationPeriod(periodStr)

				study.School, _ = m.getElementText(historyItemEle, ".name b")

				vlineTexts, _ := m.getElementTexts(historyItemEle, ".name .vline")
				if len(vlineTexts) != 2 {
					logUtil.Error(fmt.Sprintf("name vline 应该是两个值-%d，专业-学历，info:%s",
						len(vlineTexts), info.Dump()))
					return errors.New("name vline 应该是两个值，专业-学历")
				}
				study.Special, study.Education = vlineTexts[0], vlineTexts[1]

				schoolTagsEle, e := historyItemEle.FindElement(webdriver.CSS_Selector,
					".text .school-tags")
				if e == nil {
					schoolTagsStr, _ := schoolTagsEle.Text()
					if strings.Contains(schoolTagsStr, "211院校") {
						study.Is211 = true
					}
					if strings.Contains(schoolTagsStr, "985院校") {
						study.Is985 = true
					}
				}

				educationMap = append(educationMap, study)
				logUtil.Debug(fmt.Sprintf("处理完一个学历情况，详情%+v", study))
			}
			info.EducationMap = educationMap

		} else if titleStr == "项目经验" {
			historyItemEles, _ := element.FindElements(webdriver.CSS_Selector,
				".history-list .history-item")
			if len(historyItemEles) <= 0 {
				logUtil.Error(fmt.Sprintf("historyItemEles 应该大于0，info:%s, css :.history-list .history-item",
					info.Dump()))
				return errors.New("historyItemEles 应该大于0")
			}
			historyList := []entity.ProjectExperience{}
			for _, historyItemEle := range historyItemEles {
				one := entity.ProjectExperience{}
				one.Content, _ = historyItemEle.Text()

				historyList = append(historyList, one)
			}
			info.ProjectExperienceMap = historyList
		} else if titleStr == "工作经历" {
			historyItemEles, _ := element.FindElements(webdriver.CSS_Selector,
				".history-list .history-item")
			if len(historyItemEles) <= 0 {
				logUtil.Error(fmt.Sprintf("historyItemEles 应该大于0，info:%s, css :.history-list .history-item",
					info.Dump()))
				return errors.New("historyItemEles 应该大于0")
			}
			historyList := []entity.WorkExperience{}
			for _, historyItemEle := range historyItemEles {
				one := entity.WorkExperience{}
				one.Content, _ = historyItemEle.Text()

				historyList = append(historyList, one)
			}
			info.WorkExperienceMap = historyList
		} else {
			logUtil.Info(fmt.Sprintf("[%s]信息暂不处理", titleStr))

		}
	}

	return nil
}
