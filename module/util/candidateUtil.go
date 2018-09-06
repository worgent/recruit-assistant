package util

import (
	"fmt"
	cf "goBoss/config"
	"goBoss/module/entity"
	"strconv"
	"strings"

	entity2 "goBoss/module/dao"
	"time"
)

var logUtil = MyLog{}

//候选人信息
type CandidateUtil struct {
}

//获取期望薪资上下限
func (c *CandidateUtil) GetEducationPeriod(educationPeriod string) (int, int) {

	educationPeriod = strings.Replace(educationPeriod, "至今",
		string(time.Now().Year()), -1)
	educationPeriodStrs := strings.Split(educationPeriod, "-")
	if len(educationPeriodStrs) != 2 {
		logUtil.Error(fmt.Sprintf("educationPeriod illegal, %s", educationPeriod))
		return 0, 0
	} else {
		logUtil.Debug(fmt.Sprintf("educationPeriodStrs %s-%s",
			educationPeriodStrs[0], educationPeriodStrs[1]))

		inTime, e1 := strconv.Atoi(educationPeriodStrs[0])
		outTime, e2 := strconv.Atoi(educationPeriodStrs[1])
		Assert(e1)
		Assert(e2)
		if e2 != nil {
			logUtil.Error(fmt.Sprintf("educationPeriod outTime illegal, %s", educationPeriodStrs[1]))
			outTime = time.Now().Year()
		}

		return inTime, outTime
	}
}

//获取期望薪资上下限
func (c *CandidateUtil) GetExpectSalary(salaryStr string) (int, int) {
	salaryStr = strings.Replace(salaryStr, "k", "", -1)
	salaryStrs := strings.Split(salaryStr, "-")
	if len(salaryStrs) != 2 {
		return 0, 0
	} else {
		salaryLow, e1 := strconv.Atoi(salaryStrs[0])
		salaryHigh, e2 := strconv.Atoi(salaryStrs[1])
		Assert(e1)
		Assert(e2)

		return salaryLow, salaryHigh
	}
}

//获得工作经验值，应届生
func (c *CandidateUtil) GetExperienceYear(experienceStr string) int {
	if experienceStr == "应届生" {
		return 0
	} else {
		experienceStr = strings.Replace(experienceStr, "年以内", "", -1)
		experienceStr = strings.Replace(experienceStr, "年以上", "", -1)
		experienceStr = strings.Replace(experienceStr, "年", "", -1)

		age, e := strconv.Atoi(experienceStr)
		Assert(e)
		return age
	}
}

//获得年龄
func (c *CandidateUtil) GetAgeYear(ageStr string) int {
	ageStr = strings.Replace(ageStr, "岁", "", -1)
	age, e := strconv.Atoi(ageStr)
	Assert(e)
	return age

}

/**
获取详情信息
进行二次筛选，主要筛选毕业时间，工作经历和薪资的关系
## **毕业年限-薪资
	18年 8k以内
	17年 10k以内
	16年 12k以内
	15年 14k以内

返回 0-通过，101-没有教育经历，102-没有毕业时间, 103-不符合总和薪资要求， 104-

*/
func (c *CandidateUtil) FilterCandidateDetail(candidate entity.Candidate) int {
	if len(candidate.EducationMap) == 0 {
		logUtil.Info(fmt.Sprintf("候选人没有教育经历: %s-%s, 学历-%s\n",
			candidate.UID, candidate.Name,
			candidate.Education))
		return 101
	}
	firstEducation := candidate.EducationMap[0]
	if firstEducation.OutTime == 0 {
		logUtil.Info(fmt.Sprintf("候选人没有毕业时间: %s-%s, 学历-%s\n",
			candidate.UID, candidate.Name,
			firstEducation.Content))
		return 102
	}
	//假定毕业后即工作
	workYear := time.Now().Year() - firstEducation.OutTime

	//薪酬范围
	salaryConfig := cf.RConfig.SalaryExperienceConfig

	salaryLimit, e := salaryConfig[workYear]
	if e != true {
		logUtil.Info(fmt.Sprintf("候选人工作年限没有能匹配的薪资要求: %s-%s，毕业年限-%d, 期望薪资-%s\n",
			candidate.UID, candidate.Name,
			workYear, candidate.ExpectSalaryStr))
		return 103
	}
	if candidate.ExpectHigh > salaryLimit[1] {
		logUtil.Info(fmt.Sprintf("候选人不符合总和薪资要求: %s-%s，毕业年限-%d, 期望薪资-%s\n",
			candidate.UID, candidate.Name,
			workYear, candidate.ExpectSalaryStr))
		return 103
	}

	// 学历, 加强判断 详情里面是否有专科
	// 本科以上
	acceptEducations := cf.RConfig.EducationList
	if c.inArray(acceptEducations, candidate.Education) == false {
		logUtil.Info(fmt.Sprintf("候选人不符合学历要求: %s-%s, 学历-%s\n",
			candidate.UID, candidate.Name,
			candidate.Education))
		return 2
	}

	//符合条件
	return 0
}

/**
处理候选人简历
返回 0-通过，1-活跃时间不符，2-学历不符， 3-专业不符， 4-薪资不符，
	5-年龄不符， 6-经验不符，7-求职状态不符, 8-大学不在库中, 9-既不是985，211，也不是大学，也不是公办，也不是工科
*/
func (c *CandidateUtil) FilterCandidate(candidate entity.Candidate) int {
	//活跃时间 最近时间，刚刚活跃, 今日活跃，   3日内活跃，本月活跃，本周活跃
	activeTimeList := cf.RConfig.ActiveTimeList
	if c.partInArray(activeTimeList, candidate.ActiveTime) == false {
		logUtil.Info(fmt.Sprintf("候选人不符合活跃时间要求: %s-%s，活跃时间-%s\n",
			candidate.UID, candidate.Name,
			candidate.ActiveTime))
		return 1
	}

	// 学历, 加强判断 详情里面是否有专科
	// 本科以上
	acceptEducations := cf.RConfig.EducationList
	if c.inArray(acceptEducations, candidate.Education) == false {
		logUtil.Info(fmt.Sprintf("候选人不符合学历要求: %s-%s, 学历-%s\n",
			candidate.UID, candidate.Name,
			candidate.Education))
		return 2
	}
	//
	// 设定范围内，计算机，软件工程，自动化等等
	acceptSpecials := cf.RConfig.SpecialList
	if c.partInArray(acceptSpecials, candidate.Special) == false {
		logUtil.Info(fmt.Sprintf("候选人不符合专业要求: %s-%s，专业-%s\n",
			candidate.UID, candidate.Name,
			candidate.Special))
		return 3
	}

	//todo 学校筛选，二本以上； 优先筛选 985，211，导入数据库
	collegeDao := entity2.CollegeDao{}
	college, _ := collegeDao.FindCollege(candidate.School)
	if college.IsFind() == false {
		logUtil.Info(fmt.Sprintf("候选人不符合学校要求，不在库中: %s-%s，学校-%s\n",
			candidate.UID, candidate.Name,
			candidate.School))
		return 8
	}
	//既不是985，211，也不是大学，也不是公办，也不是工科
	if college.Is985() == false &&
		college.Is211() == false &&
		college.IsUniversity() == false &&
		college.IsGovermentCreate() == false &&
		college.IsEngineerType() == false {
		return 9
	}

	//期望薪资在指定范围, 超过12-20的先pass掉
	if candidate.ExpectLow > cf.RConfig.SalaryLowLimit ||
		candidate.ExpectHigh > cf.RConfig.SalaryHighLimit {
		logUtil.Info(fmt.Sprintf("候选人不符合薪资要求: %s-%s，薪资-%s\n",
			candidate.UID, candidate.Name,
			candidate.ExpectSalaryStr))
		return 4
	}

	//年龄 大于32
	if candidate.Age > cf.RConfig.AgeHighLimit ||
		candidate.Age < cf.RConfig.AgeLowLimit {
		logUtil.Info(fmt.Sprintf("候选人不符合年龄要求: %s-%s，年龄-%s\n",
			candidate.UID, candidate.Name,
			candidate.AgeStr))
		return 5
	}

	//经验大于 10
	if candidate.ExperienceYear > cf.RConfig.ExperienceHighLimit ||
		candidate.ExperienceYear < cf.RConfig.ExperienceLowLimit {
		logUtil.Info(fmt.Sprintf("候选人不符合工作年龄要求: %s-%s，工作经验-%d\n",
			candidate.UID, candidate.Name,
			candidate.ExperienceYear))
		return 6
	}

	// 求职状态，离职-随时到岗，在职-考虑机会，在职-月内到岗，在职-暂不考虑
	jobSeekingStatusList := cf.RConfig.JobSeekingStatusList
	if c.partInArray(jobSeekingStatusList, candidate.JobSeekingStatus) == false {
		logUtil.Info(fmt.Sprintf("候选人不符合求职状态要求: %s-%s，求职状态-%s\n",
			candidate.UID, candidate.Name,
			candidate.JobSeekingStatus))
		return 7
	}

	//符合条件
	return 0
}

//判断元素是否在数组中
func (c *CandidateUtil) inArray(array []string, element string) bool {
	for _, a := range array {
		if element == a {
			return true
		}
	}
	return false
}

//判断字符串是否包含数组元素
func (c *CandidateUtil) partInArray(array []string, element string) bool {
	for _, a := range array {
		if strings.Contains(element, a) == true {
			return true
		}
	}
	return false
}
