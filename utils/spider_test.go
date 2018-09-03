package utils

import (
	"fmt"
	"testing"
)

func TestSpider_Text(t *testing.T) {
	page := &Spider{
		Url: "https://www.zhipin.com/job_detail/1d0c0bbdf3a2adab1Xd_2t-4F1c~.html",
	}
	jobInfo, _ := page.Html(".detail-content .job-sec .text")
	fmt.Println(jobInfo)
	addr := page.Text(".location-address")
	fmt.Println(jobInfo, addr)
}
