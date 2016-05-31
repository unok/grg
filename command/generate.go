package command

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/tealeg/xlsx"
	"golang.org/x/oauth2"
	"net/url"
	"os"
	"regexp"
	"strings"
	"syscall"
)

type GenerateCommand struct {
	Meta
}

var max_cell_size = 32767
var page_limit = 20
var closed_bg_color = "808080"

func (c *GenerateCommand) Run(args []string) int {

	// Write your code here
	if len(args) != 1 {
		fmt.Printf("Please execute with issue list URL.\n  ex. grg generate https://github.com/unok/test_proj_for_grg/issues")
		return 1
	}
	var url_string = args[0]

	result, err := url.Parse(url_string)
	if err != nil || result.Host != "github.com" || result.Scheme != "https" {
		fmt.Printf("Please execute with issue list URL.\n  ex. grg generate https://github.com/unok/test_proj_for_grg/issues")
		return 2
	}
	var grg_token = os.Getenv("GRG_TOKEN")
	if len(grg_token) == 0 {
		fmt.Printf("GRG_TOKEN is not defined.  Please set personal API token.\n")
		return 3
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: grg_token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	match_base_url := regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)(/.*)?`)
	if !match_base_url.MatchString(url_string) {
		fmt.Printf("Please execute with issue list URL.\n  ex. grg generate https://github.com/unok/test_proj_for_grg/issues")
		return 5
	}
	match_list := match_base_url.FindAllStringSubmatch(url_string, 1)
	issue_service := *(client.Issues)
	list_option := github.ListOptions{PerPage: 100, Page: 1}
	list_option_repo := &github.IssueListByRepoOptions{ListOptions: list_option, State: "open"}
	owner, repo := match_list[0][1], match_list[0][2]

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	syscall.Unlink("issues.xlsx")
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "#"
	cell = row.AddCell()
	cell.Value = "タイトル"
	sheet.Col(1).Width = 80.0
	cell = row.AddCell()
	cell.Value = "ステータス"
	cell = row.AddCell()
	cell.Value = "担当者"
	sheet.Col(3).Width = 20.0
	cell = row.AddCell()
	cell.Value = "作成者"
	sheet.Col(4).Width = 20.0
	cell = row.AddCell()
	cell.Value = "作成日時"
	cell = row.AddCell()
	cell.Value = "更新日時"
	cell = row.AddCell()
	cell.Value = "ラベル"
	sheet.Col(7).Width = 20.0
	cell = row.AddCell()
	cell.Value = "リンク"
	sheet.Col(8).Width = 50.0
	cell = row.AddCell()
	cell.Value = "詳細"
	sheet.Col(9).Width = 70.0
	cell = row.AddCell()
	cell.Value = "最新コメント"
	sheet.Col(10).Width = 70.0
	cell = row.AddCell()
	cell.Value = ""
	cell = row.AddCell()
	cell.Value = ""
	sheet.Col(12).Width = 20.0
	for {
		issues, resp, err := issue_service.ListByRepo(owner, repo, list_option_repo)
		if err != nil {
			return 4
		}
		for i := 0; i < len(issues); i++ {
			style := xlsx.NewStyle()
			style.Alignment = xlsx.Alignment{Vertical: "top"}
			issue := issues[i]
			if strings.Contains(*issue.HTMLURL, "/pull/") {
				continue
			}
			row = sheet.AddRow()
			row.Height = 20.0
			if *issue.State == "closed" {
				style.Fill = *xlsx.NewFill("solid", closed_bg_color, "")
			} else {
				style.Fill = *xlsx.NewFill("", "", "")
			}
			cell = row.AddCell()
			cell.SetInt(*issue.Number)
			cell.SetStyle(style)
			cell = row.AddCell()
			cell.Value = *issue.Title
			cell.SetStyle(&xlsx.Style{Alignment: xlsx.Alignment{WrapText: true, Vertical: "top"}, Fill: style.Fill})
			cell = row.AddCell()
			cell.Value = *issue.State
			cell.SetStyle(style)
			cell = row.AddCell()
			if issue.Assignee != nil {
				cell.Value = *issue.Assignee.Login
			}
			cell.SetStyle(style)
			cell = row.AddCell()
			cell.Value = *issue.User.Login
			cell.SetStyle(style)
			cell = row.AddCell()
			cell.SetDateTime(*issue.CreatedAt)
			cell.SetStyle(style)
			cell = row.AddCell()
			cell.SetDateTime(*issue.UpdatedAt)
			cell.SetStyle(style)
			cell = row.AddCell()
			label_string := ""
			for _, v := range issue.Labels {
				label_string += fmt.Sprintf("%s\n", *v.Name)
			}
			cell.Value = label_string
			cell.SetStyle(&xlsx.Style{Alignment: xlsx.Alignment{WrapText: true, Vertical: "top"}, Fill: style.Fill})
			cell = row.AddCell()
			cell.Value = *issue.HTMLURL
			cell.SetStyle(style)
			cell = row.AddCell()
			cell.Value = AdjustCellSize(*issue.Body)
			cell.SetStyle(&xlsx.Style{Alignment: xlsx.Alignment{WrapText: true, Vertical: "top"}, Fill: style.Fill})

			cell = row.AddCell()
			if *issue.Comments > 0 {
				comments, _, _ := issue_service.ListComments(owner, repo, *issue.Number, nil)
				last_index := len(comments) - 1
				cell.Value = AdjustCellSize(*comments[last_index].Body)
				cell.SetStyle(&xlsx.Style{Alignment: xlsx.Alignment{WrapText: true, Vertical: "top"}, Fill: style.Fill})
				cell = row.AddCell()
				cell.SetDateTime(*comments[last_index].CreatedAt)
				cell.SetStyle(style)
				cell = row.AddCell()
				cell.Value = *comments[last_index].User.Login
				cell.SetStyle(style)
			} else {
				cell.SetStyle(style)
				cell = row.AddCell()
				cell.SetStyle(style)
				cell = row.AddCell()
				cell.SetStyle(style)
			}
		}
		fmt.Printf("Page %d: %s\n", list_option_repo.ListOptions.Page, client.Rate().String())
		if resp.NextPage == 0 || resp.NextPage > page_limit {
			break
		}
		list_option_repo.ListOptions.Page = resp.NextPage
	}

	err = file.Save("issues.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
	return 0
}

func AdjustCellSize(str string) string {
	str = strings.TrimSpace(str)
	r := []rune(str)
	if len(r) > max_cell_size {
		return string(r[0:(max_cell_size - 1)])
	} else {
		return string(r)
	}
}

func (c *GenerateCommand) Synopsis() string {
	return ""
}

func (c *GenerateCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
