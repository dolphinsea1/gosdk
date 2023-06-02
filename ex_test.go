package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/gjson"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_DiffLevel(t *testing.T) {
	var a, b map[string]interface{}
	if err := json.Unmarshal([]byte(json1), &a); err != nil {
		fmt.Println("error:", err)
		return
	}
	if err := json.Unmarshal([]byte(json2), &b); err != nil {
		fmt.Println("error:", err)
		return
	}

	start := time.Now()
	var ignoreField = []string{"dstAddrs", "srcAddrs"}
	var keyField = []string{"dstAddrs"}
	diff, err := CompareJSONWithLevel(json1, json2, 3, keyField, ignoreField)
	if err != nil {
		// 解析 JSON 失败，输出错误信息
		fmt.Errorf("%s", err)
	}
	for _, report := range diff {
		fmt.Printf("action:%v\ndiff:%v\n", report.Action, report.Diff)
	}

	fmt.Printf("耗时：%v\n", time.Since(start))
	// 输出到文件
	err = OutputDiffReportToFile(diff, "diff_report.txt")
	if err != nil {
		// 输出失败，输出错误信息
		fmt.Println(err)
	} else {
		fmt.Println("Diff report has been saved to diff_report")
	}

}

// /======================================

func OutputDiffReportToFile(reports []DiffRes, filename string) error {
	// 打开输出文件并准备写入数据
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 json.MarshalIndent() 函数将差异报告转换成 JSON 格式
	data, err := json.MarshalIndent(reports, "", "  ")
	if err != nil {
		return err
	}

	// 将数据写入文件
	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
func TestName(t *testing.T) {
	res, _ := CompareWithLevel(json1, json2, 3, []string{}, []string{"srcZones"})
	bytes, _ := json.Marshal(res)
	fmt.Println(string(bytes))

	/*for _, re := range res {
		bytes, _ := json.Marshal(re)
		fmt.Println(string(bytes))
	}*/
	/*results := gjson.Parse(json1).Array()
	for _, result := range results {
		fmt.Println(result.Get())
	}
	*/
}

func CompareWithLevel(json1, json2 string, level int, keyField, ignoreField []string) ([]DiffRes, error) {
	// 解析JSON字符串
	oldJson := gjson.Parse(json1)
	newJson := gjson.Parse(json2)

	reports := diffList(oldJson, newJson, level, keyField, ignoreField)
	res := make([]DiffRes, 0)
	for s, diffReports := range reports {
		res = append(res, DiffRes{
			Action: s,
			Diff:   diffReports,
		})
	}

	return res, nil
}

func diffList(oldJson, newJson gjson.Result, level int, keyField, ignoreField []string) map[string][]DiffReport {
	reports := make(map[string][]DiffReport)
	if level <= 1 {
		if !cmp.Equal(oldJson, newJson) {
			reports[Modify] = append(reports[Modify], DiffReport{
				OldData: oldJson.Str,
				NewData: newJson.Str,
			})
		}
		return reports
	}

	// 获取所有键名
	oldField := make([]string, 0)
	getAllKeys(oldJson, &oldField, level)
	newField := make([]string, 0)
	getAllKeys(newJson, &newField, level)
	fmt.Println(oldField, "\n", newField)
	// 比对根据层级比对

	for _, field := range oldField {
		//存在则忽略不对比
		if containsList(ignoreField, field) {
			continue
		}
		if len(strings.Split(field, ".")) == (level - 1) {
			oldData := oldJson.Get(field)
			newData := newJson.Get(field)
			for _, s := range newField {
				s1 := strings.Split(field, ",")
				s2 := strings.Split(s, ",")
				if len(s1) == len(s2) && strings.EqualFold(s1[len(s1)-1], s2[len(s2)-1]) {
					newData = newJson.Get(s)
				}
			}
			if newData.Exists() {
				if !cmp.Equal(oldData.Raw, newData.Raw) {
					reports[Modify] = append(reports[Modify], DiffReport{
						Field:   field,
						OldData: oldData.Raw,
						NewData: newData.Raw,
					})
				}
			} else {
				reports[Delete] = append(reports[Delete], DiffReport{
					Field:   field,
					OldData: oldData.Raw,
				})
			}
		}
	}
	for _, field := range newField {
		if containsList(ignoreField, field) {
			continue
		}
		if len(strings.Split(field, ".")) == (level - 1) {
			oldData := oldJson.Get(field)
			newData := newJson.Get(field)
			for _, s := range newField {
				s1 := strings.Split(field, ",")
				s2 := strings.Split(s, ",")
				if len(s1) == len(s2) && strings.EqualFold(s1[len(s1)-1], s2[len(s2)-1]) {
					oldJson = oldJson.Get(s)
				}
			}
			if !oldData.Exists() {
				reports[Add] = append(reports[Add], DiffReport{
					Field:   field,
					NewData: newData.Raw,
				})
			}
		}
	}

	return reports
}

func Test(t *testing.T) {
	m := make(map[int][]interface{})
	traverseJson(gjson.Parse(json1), "", 3, m)
	for i, v := range m {
		fmt.Println(i, v)
	}

}

// 递归遍历 JSON 对象
func traverseJson(result gjson.Result, indent string, level int, m map[int][]interface{}) {
	// 获取当前节点的数据类型
	valueType := result.Type

	// 输出当前节点的数据类型和值
	fmt.Printf("%s- 类型：%v，值：%v\n", indent, valueType.String(), result.Value())
	if level == 0 {
		return
	}
	// 如果当前节点是对象类型，需要遍历子节点
	if valueType == gjson.JSON {
		result.ForEach(func(key, value gjson.Result) bool {
			// 输出下一层级节点前的缩进
			nextIndent := fmt.Sprintf("%s", indent)
			if gjson.Parse(result.Str).IsArray() {
				m[level] = append(m[level], result.Value())
				traverseJson(value, fmt.Sprintf("%s%s", nextIndent, key.String()), level-1, m)
			} else if gjson.Parse(result.Str).IsObject() {
				m[level] = append(m[level], result.Value())
				traverseJson(value, fmt.Sprintf("%s%s", nextIndent, key.String()), level-1, m)
			} else {
				m[level] = append(m[level], result.Value())
			}
			// 递归遍历子节点

			return true // 必须返回 true，否则 ForEach 函数会提前结束
		})
	}
}

// 递归函数，获取JSON对象内部所有层级的键名，包括数组下标
func getAllKeys(value gjson.Result, keys *[]string, maxLevel int, parentKeys ...string) {
	if value.IsObject() {
		*keys = append(*keys, concatenateKeys(parentKeys...))
		if 0 <= maxLevel {
			value.ForEach(func(key, value gjson.Result) bool {
				getAllKeys(value, keys, maxLevel-1, append(parentKeys, key.String())...)
				if maxLevel == 0 {
					return false
				}
				return true
			})
		}
	} else if value.IsArray() {
		*keys = append(*keys, concatenateKeys(parentKeys...))
		if 0 <= maxLevel {
			value.ForEach(func(idx, value gjson.Result) bool {
				getAllKeys(value, keys, maxLevel-1, append(parentKeys, strconv.Itoa(int(idx.Int())))...)
				if maxLevel == 0 {
					return false
				}
				return true
			})

		}
	} else {
		if 0 <= maxLevel {
			*keys = append(*keys, concatenateKeys(parentKeys...))
		}
	}
}

// 辅助函数，将所有键名拼接成一个字符串
func concatenateKeys(keys ...string) string {
	result := ""
	for _, key := range keys {
		if result != "" {
			result = fmt.Sprintf("%s.%s", result, key)
		} else {
			result = key
		}
	}
	return result
}

// 判断一个字符串是否在一个字符串数组中
func containsList(arr []string, str string) bool {
	split := strings.Split(str, ".")
	for _, s := range arr {
		for _, s2 := range split {
			if s == s2 {
				return true
			}
		}

	}
	return false
}

//-----------------------------------------------------

func Test2(t *testing.T) {

	res, _ := CompareJSON(DiffReq{
		JsonOld:     json1,
		JsonNew:     json2,
		DeLevel:     3,
		KeyField:    nil,
		IgnoreField: []string{"group"},
	})
	bytes, _ := json.Marshal(res)
	fmt.Println(string(bytes))
}

// 比较两个 JSON 字符串，返回差异数据列表
func CompareJSON(req DiffReq) (*[]DiffRes, error) {
	oldJSON := []byte(req.JsonOld)
	newJSON := []byte(req.JsonNew)
	diffReports := make(map[string][]DiffReport, 0)
	var oldObject, newObject interface{}

	if err := json.Unmarshal(oldJSON, &oldObject); err != nil {
		return nil, fmt.Errorf("解析旧JSON数据失败：%v", err)
	}

	if err := json.Unmarshal(newJSON, &newObject); err != nil {
		return nil, fmt.Errorf("解析新JSON数据失败：%v", err)
	}

	compare("", oldObject, newObject, req.DeLevel, req.KeyField, req.IgnoreField, diffReports)
	res := make([]DiffRes, 0)
	for k, report := range diffReports {
		res = append(res, DiffRes{Action: k, Diff: report})
	}
	return &res, nil
}

// 对比两个值
func compare(fieldName string, oldVal, newVal interface{}, deLevel int, keyFields, ignoreFields []string, diffReports map[string][]DiffReport) {
	if deLevel == 1 {
		if !cmp.Equal(oldVal, newVal) {
			addDiffReport(fieldName, oldVal, newVal, diffReports)
		}
		return
	}

	oldKind := reflect.TypeOf(oldVal).Kind()
	newKind := reflect.TypeOf(newVal).Kind()
	fmt.Println(oldKind, newKind)
	if oldKind != newKind {
		addDiffReport(fieldName, oldVal, newVal, diffReports)
		return
	}

	switch oldKind {
	case reflect.Map:
		oldMap := oldVal.(map[string]interface{})
		newMap := newVal.(map[string]interface{})

		for key, oldItem := range oldMap {
			newItem, ok := newMap[key]
			if !ok {
				addDiffReport(fmt.Sprintf("%s.%s", fieldName, key), oldItem, nil, diffReports)
				continue
			}

			if isKeyField(keyFields, fieldName, key) && !reflect.DeepEqual(oldItem, newItem) {
				compare(fmt.Sprintf("%s.%s", fieldName, key), oldItem, newItem, deLevel-1, keyFields, ignoreFields, diffReports)
			} else if !isIgnoreField(ignoreFields, fieldName, key) {
				compare(fmt.Sprintf("%s.%s", fieldName, key), oldItem, newItem, deLevel, keyFields, ignoreFields, diffReports)
			}
		}

		for key, newItem := range newMap {
			_, ok := oldMap[key]
			if !ok {
				if !isIgnoreField(ignoreFields, fieldName, key) {
					addDiffReport(fmt.Sprintf("%s.%s", fieldName, key), nil, newItem, diffReports)
				}
			}
		}
	case reflect.Slice, reflect.Array:
		oldSlice := oldVal.([]interface{})
		newSlice := newVal.([]interface{})

		if deLevel == 2 {
			if len(oldSlice) != len(newSlice) {
				addDiffReport(fieldName, oldVal, newVal, diffReports)
				return
			}

			for i := range oldSlice {
				if !reflect.DeepEqual(oldSlice[i], newSlice[i]) {
					addDiffReport(fmt.Sprintf("%s[%d]", fieldName, i), oldSlice[i], newSlice[i], diffReports)
				}
			}
		} else {
			for i, oldItem := range oldSlice {
				if len(newSlice) <= i {
					addDiffReport(fmt.Sprintf("%s[%d]", fieldName, i), oldItem, nil, diffReports)
					continue
				}

				compare(fmt.Sprintf("%s[%d]", fieldName, i), oldItem, newSlice[i], deLevel-1, keyFields, ignoreFields, diffReports)
			}

			for i := len(oldSlice); i < len(newSlice); i++ {
				addDiffReport(fmt.Sprintf("%s[%d]", fieldName, i), nil, newSlice[i], diffReports)
			}
		}
	default:
		if !reflect.DeepEqual(oldVal, newVal) {
			addDiffReport(fieldName, oldVal, newVal, diffReports)
		}
	}
}

// 添加差异数据
func addDiffReport(fieldName string, oldVal, newVal interface{}, diffReports map[string][]DiffReport) {
	action := Modify
	if oldVal == nil {
		action = Add
	}
	if newVal == nil {
		action = Delete
	}
	diffReports[action] = append(diffReports[action], DiffReport{
		Field:   fieldName,
		OldData: oldVal,
		NewData: newVal,
	})
}

// 判断是否为需要对比key
func isKeyField(keyFields []string, fieldName string, key string) bool {
	if len(keyFields) == 0 {
		return true
	}

	for _, item := range keyFields {
		if strings.HasSuffix(item, key) || strings.HasSuffix(item, fmt.Sprintf(".%s", key)) ||
			strings.HasSuffix(item, fmt.Sprintf(".%s.", key)) || strings.HasPrefix(item, fieldName) {
			return true
		}
	}

	return false
}

// 判断是否为忽略字段
func isIgnoreField(ignoreFields []string, fieldName string, key string) bool {
	if len(ignoreFields) == 0 {
		return false
	}

	for _, item := range ignoreFields {
		fmt.Println(ignoreFields, fieldName, key)

		fmt.Println(strings.Contains(key, item) || strings.Contains(fmt.Sprintf(".%s", key), item) ||
			strings.Contains(fmt.Sprintf(".%s.", key), item) || strings.Contains(fieldName, item))
		if strings.Contains(key, item) || strings.Contains(fmt.Sprintf(".%s", key), item) ||
			strings.Contains(fmt.Sprintf(".%s.", key), item) || strings.Contains(fieldName, item) {

			return true
		}
	}

	return false
}
