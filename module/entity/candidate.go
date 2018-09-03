package entity

import "fmt"

//教育信息
type Study struct {
	Content string

	//毕业院校
	School string
	//专业
	Special string
	//学历
	Education string
	//在校时间
	InTime  int
	OutTime int

	Is211 bool
	Is985 bool
}

//工作经历
type WorkExperience struct {
	Content string
}

//项目经验
type ProjectExperience struct {
	Content string
}

//候选人信息
type Candidate struct {
	UID string
	//性别，通过详情获取
	Sex    string
	Name   string
	Avatar string

	//年龄
	AgeStr string
	Age    int
	//地区
	District string
	//个人优势
	Advantage string
	//推荐原因
	RecommendReason string

	//经验
	ExperienceYearStr string
	ExperienceYear    int

	//工作经历
	JobHistory string
	//期望薪资
	ExpectSalaryStr string
	ExpectLow       int
	ExpectHigh      int
	//求职状态，
	JobSeekingStatus string

	//毕业院校
	School string
	//专业
	Special string
	//学历
	Education string

	//活跃时间
	ActiveTime string
	//获取时间
	CreatedTime int64

	//沟通状态
	Communicatetatus string

	//处理结果
	DealResult int

	///////以下为扩展信息，打开详情页获取
	//详情文本
	Content string
	//详细教育信息
	EducationMap []Study
	//项目经验信息
	ProjectExperienceMap []ProjectExperience
	//工作经历
	WorkExperienceMap []WorkExperience
}

//dump
func (c *Candidate) Dump() string {
	return fmt.Sprintf("%+v", *c)
}

func (c *Candidate) BriefDump() string {
	return fmt.Sprintf("{ UID:%s, Name:%s, Education:%s, school:%s.%s, "+
		"ExpectSalaryStr:%s, AgeStr:%s, ExperienceYearStr:%s, JobSeekingStatus:%s} ",
		c.UID, c.Name, c.Education, c.School, c.Special,
		c.ExpectSalaryStr, c.AgeStr, c.ExperienceYearStr, c.JobSeekingStatus)
}

func (c *Candidate) EducationMapDump() string {
	ret := ""
	for _, study := range c.EducationMap {
		ret += fmt.Sprintf("%+v\n", study)
	}
	return ret
}

//获得最高学历

//学历中是否有专科
