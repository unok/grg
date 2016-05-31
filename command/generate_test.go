package command

import (
	"os"
	"strings"
	"testing"
)

const GITHUB_TEST_URL = "https://github.com/unok/test_proj_for_grg/issues"

func TestGenerateCommand_args_ok(t *testing.T) {
	var c = &GenerateCommand{}
	var args = []string{GITHUB_TEST_URL}
	actual := c.Run(args)
	expected := 0
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGenerateCommand_args_ng(t *testing.T) {
	var c = &GenerateCommand{}
	var args = []string{}
	actual := c.Run(args)
	expected := 1
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGenerateCommand_args_not_github_url(t *testing.T) {
	var c = &GenerateCommand{}
	var args = []string{"aaaaa"}
	actual := c.Run(args)
	expected := 2
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGenerateCommand_env_not_found(t *testing.T) {
	var c = &GenerateCommand{}
	var args = []string{GITHUB_TEST_URL}
	os.Unsetenv("GRG_TOKEN")
	actual := c.Run(args)
	expected := 3
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGenerateCommand_adjust_none(t *testing.T) {
	str := strings.Repeat("a", max_cell_size)
	actual := AdjustCellSize(str)
	expected := str
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGenerateCommand_adjust_cut(t *testing.T) {
	str := strings.Repeat("a", max_cell_size+1)
	actual := AdjustCellSize(str)
	expected := str
	if actual == expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
	if len(actual) == 32767 {
		t.Errorf("got %v\nwant %v", len(actual), 32767)
	}
}
