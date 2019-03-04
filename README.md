# recruit-assistant

基于go语言的boss直聘招聘辅助工具，支持学历，学校(985,211)，专业，薪资，经验等方面的筛选，并进行主动沟通。

详细配置：
{

  "MysqlConnectStr": "root:@tcp(localhost:3306)/boss",
  "ChromeApp": "/Applications/Google\\ Chrome.app/Contents/MacOS/Google\\ Chrome",
  "CommunicateLimit": 5,
  "ResumeFilterLimit": 100,
  "ResumePageLimit": 10,
  "EducationList": [
    "本科", "硕士"
  ],
  "SpecialList": [
    "计算机", "软件", "网络", "通信"
  ],
  "AgeLowLimit": 20,
  "AgeHighLimit": 32,

  "ExperienceLowLimit": 0,
  "ExperienceHighLimit": 8,

  "SalaryLowLimit": 12,
  "SalaryHighLimit": 16,

  "ActiveTimeList": [
    "刚刚活跃", "今日活跃"
  ],

  "JobSeekingStatusList": [
    "离职-随时到岗", "在职-月内到岗", "在职-暂不考虑"
  ],

  "SalaryExperienceConfig" : {
    "0" : [6,8],
    "1" : [7,9],
    "2" : [10,12],
    "3" : [10,12],
    "4" : [10,15],
    "5" : [10,15]
  }
}
