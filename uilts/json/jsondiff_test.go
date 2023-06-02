package json

import (
	"fmt"
	"testing"
)

const json1 = `[{
		"name": "张so",
		"age": 18,
		"hobbies": ["篮球", "游泳", "旅游"]
	},
	{
		"name": "李四",
		"age": 20,
		"hobbies": ["足球", "游戏"]
	}
]`

const json2 = `[{
		"name": "张三",
		"age": 20,
		"hobbies": ["篮球", "游泳", "旅游", "阅读"]
	},
	{
		"name": "李四",
		"age": 20,
		"hobbies": ["足球", "游戏", "阅读"]
	},
	{
		"name": "王五",
		"age": 22,
		"hobbies": ["音乐", "电影"]
	}
]`

const json3 = `{
		"name": "张so",
		"age": 18,
		"hobbies": ["篮球", "游泳", "旅游"]
	}`

const json4 = `{
		"name": "张三",
		"age": 20,
		"hobbies": ["篮球", "游泳", "旅游", "阅读"],
		"sex":"男"
	}`

func TestName(t *testing.T) {

	var keyField = []string{}
	var ignoreField = []string{"hobbies"}

	diff, err := CompareJSONWithLevel(json1, json2, 2, keyField, ignoreField)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	for _, report := range diff {
		fmt.Printf("action:%v\ndiff:%v\n", report.Action, report.Diff)
	}

	diff1, err := CompareJSONWithLevel(json3, json4, 3, keyField, ignoreField)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	for _, report := range diff1 {
		fmt.Printf("action1:%v\ndiff1:%v\n", report.Action, report.Diff)
	}

}
