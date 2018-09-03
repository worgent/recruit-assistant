package entity

import "fmt"

//教育信息
type College struct {
	Id       int
	Name     string
	Province string
	Distict  string
	City     string
	//211
	Info211 string
	//985
	Info985 string
	//公办，民办
	CreateType string
	//专业类型，综合，工科
	SpecialType string
	//归属
	Belong string
	//大学，学院
	EducationType string
	//批次
	Batch string
}

func (c *College) IsFind() bool {
	if c.Id == 0 {
		return false
	} else {
		return true
	}
}

//dump
func (c *College) Is985() bool {
	if c.Info985 == "985" {
		return true
	} else {
		return false
	}
}

func (c *College) Is211() bool {
	if c.Info211 == "211" {
		return true
	} else {
		return false
	}
}

//公办
func (c *College) IsGovermentCreate() bool {
	if c.CreateType == "公办院校" {
		return true
	} else {
		return false
	}
}

func (c *College) IsUniversity() bool {
	if c.CreateType == "大学" {
		return true
	} else {
		return false
	}
}

func (c *College) IsEngineerType() bool {
	if c.SpecialType == "工科类院校" {
		return true
	} else {
		return false
	}
}

//dump
func (c *College) Dump() string {
	return fmt.Sprintf("%+v", *c)
}
