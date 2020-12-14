package main

import (
	"net/url"
	"strings"
	"time"

	"fmt"
	"log"

	"github.com/tealeg/xlsx"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var (
	webDriver selenium.WebDriver
	service   *selenium.Service
)
var (
	jobDetailList  []*JobDetailModel
	jobDetailModel *JobDetailModel
	detailUrls     []string
)
var (
	file     *xlsx.File
	sheet    *xlsx.Sheet
	row      *xlsx.Row
	cell     *xlsx.Cell
	excelErr error
)

func main() {
	StartWebDriverByChrome()
	LoginZhiPin()
	//RecommendZhiPin()
	SearchList("java技术总监")
	TaskFinish()
}

func StartWebDriverByChrome() (selenium.WebDriver, bool) {
	opts := []selenium.ServiceOption{}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	// 禁止加载图片，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 1,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			//"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
		},
	}
	caps.AddChrome(chromeCaps)
	var err error
	// 启动chromedriver，端口号可自定义
	service, err = selenium.NewChromeDriverService("./chromedriver", 9516, opts...)
	if err != nil {
		log.Printf("Error starting the ChromeDriver server: %v", err)
	}
	// 调起chrome浏览器
	webDriver, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9516))
	if err != nil {
		panic(err)
		return webDriver, false
	}
	return webDriver, true
}
func RecommendZhiPin() {
	targeUrl := "https://www.zhipin.com/web/geek/recommend"
	err := webDriver.Get(targeUrl)

	if err != nil {
		panic(fmt.Sprintf("Failed to load page: %s\n", err))
	}
	log.Println(webDriver.GetCookies())
	time.Sleep(5 * time.Second)
}
func LoginZhiPin() {
	targeUrl := "https://login.zhipin.com/?ka=header-login"
	err := webDriver.Get(targeUrl)

	if err != nil {
		panic(fmt.Sprintf("Failed to load page: %s\n", err))
	}
	log.Println(webDriver.GetCookies())
	time.Sleep(15 * time.Second)
}
func SearchList(keyWord string) bool {
	file = xlsx.NewFile()
	sheet, excelErr = file.AddSheet("Sheet1")
	if excelErr != nil {
		fmt.Printf(excelErr.Error())
	}
	row = sheet.AddRow()
	row.SetHeightCM(0.5)
	cell = row.AddCell()
	cell.Value = "公司名称"
	cell = row.AddCell()
	cell.Value = "公司财务"
	cell = row.AddCell()
	cell.Value = "规模"
	cell = row.AddCell()
	cell.Value = "行业"
	cell = row.AddCell()
	cell.Value = "公司网址"
	cell = row.AddCell()
	cell.Value = "招聘职位"
	cell = row.AddCell()
	cell.Value = "薪资"
	cell = row.AddCell()
	cell.Value = "地区"
	cell = row.AddCell()
	cell.Value = "工作经验"
	cell = row.AddCell()
	cell.Value = "学历"
	cell = row.AddCell()
	cell.Value = "福利"
	cell = row.AddCell()
	cell.Value = "职位描述"

	escapeKeyword := url.QueryEscape(keyWord)
	targeUrl := "https://www.zhipin.com/job_detail/?query=" + escapeKeyword + "&city=101010100&position=100101"
	err := webDriver.Get(targeUrl)
	if err != nil {
		panic(fmt.Sprintf("Failed to load page: %s\n", err))
		return false
	}
	jobListItem, err := webDriver.FindElement(selenium.ByClassName, "job-list")
	if err != nil {
		panic(fmt.Sprintf("Failed to load job list: %s\n", err))
		return false
	}
	jobTabItem, err := jobListItem.FindElement(selenium.ByClassName, "job-tab")
	if err != nil {
		panic(fmt.Sprintf("Failed to load job tab: %s\n", err))
		return false
	}
	totalCountStr, err := jobTabItem.GetAttribute("data-rescount")
	if err != nil {
		panic(fmt.Sprintf("Failed to load job total count: %s\n", err))
		return false
	}
	fmt.Println(totalCountStr)
	jobItemDetailList, err := jobListItem.FindElements(selenium.ByClassName, "primary-box")
	if err != nil {
		panic(fmt.Sprintf("Failed to load job detail item: %s\n", err))
		return false
	}

	for _, webItem := range jobItemDetailList {
		jobUrl, detailErr := webItem.GetAttribute("href")
		if detailErr != nil {
			panic(fmt.Sprintf("Failed to load job detail Url href : %s\n", err))
			return false
		}
		ka, detailErr := webItem.GetAttribute("ka")
		if detailErr != nil {
			panic(fmt.Sprintf("Failed to load job detail Url ka: %s\n", err))
			return false
		}
		data_lid, detailErr := webItem.GetAttribute("data-lid")
		if detailErr != nil {
			panic(fmt.Sprintf("Failed to load job detail Url lid: %s\n", err))
			return false
		}
		targeUrl = "https://www.zhipin.com" + jobUrl + "?ka=" + ka + "&lid=" + data_lid
		detailUrls = append(detailUrls, targeUrl)
	}
	var dataRow *xlsx.Row
	for _, detailUrl := range detailUrls {
		dataRow = sheet.AddRow()

		err := webDriver.Get(detailUrl)
		if err != nil {
			fmt.Sprintf("Failed to load detail page: %s\n", err)
			continue
		}
		jobDetailModel = &JobDetailModel{}
		mainContent, err := webDriver.FindElement(selenium.ByID, "main")
		if err != nil {
			fmt.Sprintf("Failed to load detail page: %s\n", err)
			continue
		}
		infoCompany, err := mainContent.FindElement(selenium.ByClassName, "sider-company")
		if err != nil {
			fmt.Sprintf("Failed to load detail page infoCompany: %s\n", err)
			continue
		}
		infoPrimary, err := mainContent.FindElement(selenium.ByClassName, "info-primary")
		if err != nil {
			fmt.Sprintf("Failed to load detail page infoPrimary: %s\n", err)
			continue
		}
		fmt.Println(infoPrimary)
		InfoJobDetail, err := mainContent.FindElement(selenium.ByClassName, "job-detail")
		time.Sleep(5 * time.Second)
		var detailItemErr error
		//公司
		jobCompanyNameItems, detailItemErr := infoCompany.FindElements(selenium.ByCSSSelector, "div.company-info a")
		if detailItemErr != nil {
			panic(fmt.Sprintf("Failed to load jobCompanyNameItems: %s\n", err))
		}
		for index, _ := range jobCompanyNameItems {
			if index == 1 {
				jobDetailModel.JobCompanyName, detailItemErr = jobCompanyNameItems[1].Text()
			}
		}

		JobCompanyItems, detailItemErr := infoCompany.FindElements(selenium.ByCSSSelector, "p")
		if detailItemErr != nil {
			panic(fmt.Sprintf("Failed to load JobCompanyItems: %s\n", err))
		}
		for index, _ := range JobCompanyItems {
			if index == 1 {
				jobDetailModel.JobCompanyFiance, detailItemErr = JobCompanyItems[1].Text()
			} else if index == 2 {
				jobDetailModel.JobCompanyScale, detailItemErr = JobCompanyItems[2].Text()
			} else if index == 3 {
				jobDetailModel.JobCompanyIndustry, detailItemErr = JobCompanyItems[3].Text()
			} else if index == 4 {
				jobDetailModel.JobCompanyUrl, detailItemErr = JobCompanyItems[4].Text()
			}
		}

		// 职位
		jobNameItem, detailItemErr := infoPrimary.FindElement(selenium.ByCSSSelector, ".name h1")
		if detailItemErr != nil {
			panic(fmt.Sprintf("Failed to load jobNameItem: %s\n", err))
		}
		jobDetailModel.JobName, detailItemErr = jobNameItem.Text()
		jobSalaryItem, detailItemErr := infoPrimary.FindElement(selenium.ByCSSSelector, ".name span")
		if detailItemErr != nil {
			panic(fmt.Sprintf("Failed to load jobSalaryItem: %s\n", err))
		}
		jobDetailModel.JobSalary, detailItemErr = jobSalaryItem.Text()
		jobCityItem, detailItemErr := infoPrimary.FindElement(selenium.ByCSSSelector, "p a.text-city")
		if detailItemErr != nil {
			panic(fmt.Sprintf("Failed to load jobCityItem: %s\n", err))
		}
		jobDetailModel.JobCity, detailItemErr = jobCityItem.Text()
		metaItem, detailItemErr := webDriver.FindElement(selenium.ByXPATH, "//meta[@name='description']")

		metaItemText, detailItemErr := metaItem.GetAttribute("content")
		metaItemArray := strings.Split(metaItemText, "，")
		for _, metaItemValue := range metaItemArray {
			metaItemKeyPair := strings.Split(metaItemValue, "：")
			if len(metaItemKeyPair) == 2 {
				if metaItemKeyPair[0] == "要求" {
					jobDetailModel.JobYears = metaItemKeyPair[1]
				} else if metaItemKeyPair[0] == "学历" {
					jobDetailModel.JobEducation = metaItemKeyPair[1]
				}
			}
		}
		jobTagItems, detailItemErr := infoPrimary.FindElements(selenium.ByCSSSelector, ".job-tags span")
		if detailItemErr != nil {
			panic(fmt.Sprintf("Failed to load jobCityItem: %s\n", err))
		}
		var jboTags []string
		if len(jobTagItems) > 0 {
			for _, jobTag := range jobTagItems {
				jobTagValue, detailItemErr := jobTag.Text()
				if detailItemErr == nil {
					jboTags = append(jboTags, jobTagValue)
				}
			}
			jobDetailModel.JobTag = strings.Join(jboTags, " ")
		} else {
			jobDetailModel.JobTag = ""
		}

		jobDescItem, detailItemErr := InfoJobDetail.FindElement(selenium.ByCSSSelector, ".detail-content .text")
		if detailItemErr != nil {
			panic(fmt.Sprintf("Failed to load jobDescItem: %s\n", err))
		}
		jobDetailModel.JobDesc, detailItemErr = jobDescItem.Text()
		fmt.Println(jobDetailModel)
		jobDetailList = append(jobDetailList, jobDetailModel)

		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobCompanyName)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobCompanyFiance)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobCompanyScale)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobCompanyIndustry)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobCompanyUrl)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobName)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobSalary)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobCity)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobYears)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobEducation)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobTag)
		cell = dataRow.AddCell()
		cell.SetValue(jobDetailModel.JobDesc)
		time.Sleep(20 * time.Second)
	}
	excelErr = file.Save(keyWord + "_北京.xlsx")
	return true
}

type JobDetailModel struct {
	JobCompanyName     string //公司名称
	JobCompanyFiance   string //公司财务状况
	JobCompanyScale    string //公司规模
	JobCompanyIndustry string //公司行业
	JobCompanyUrl      string //公司网址
	JobName            string
	JobSalary          string
	JobCity            string
	JobYears           string
	JobEducation       string
	JobTag             string
	JobDesc            string
}

func TaskFinish() {
	service.Stop()
	webDriver.Quit()
}

