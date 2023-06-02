package json

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
)

type DiffReport struct {
	Field   string      `json:"field"`
	OldData interface{} `json:"old_data,omitempty"`
	NewData interface{} `json:"new_data,omitempty"`
}

type DiffRes struct {
	Action string       `json:"action"` // ADD  MOD DEL
	Diff   []DiffReport `json:"diff"`
}

type DiffReq struct {
	JsonOld     string   `json:"json_old"`
	JsonNew     string   `json:"json_new"`
	DeLevel     int      `json:"de_level"`
	KeyField    []string `json:"key_field"`
	IgnoreField []string `json:"ignore_field"`
}

const (
	Modify = "MOD"
	Add    = "ADD"
	Delete = "DEL"
)

func CompareJSONWithLevel(json1, json2 string, level int, keyField, ignoreField []string) ([]DiffRes, error) {
	// 解析 JSON 字符串
	var data1 interface{}
	if err := json.Unmarshal([]byte(json1), &data1); err != nil {
		return nil, err
	}

	var data2 interface{}
	if err := json.Unmarshal([]byte(json2), &data2); err != nil {
		return nil, err
	}

	// 初始化差异报告列表
	reports := make([]DiffReport, 0)

	if level > 3 {
		level = 3
	}

	// 比较两个 JSON 对象的差异
	compareJSON(data1, data2, "", &reports, level, keyField, ignoreField)

	date := make(map[string][]DiffReport)
	for _, report := range reports {
		if report.NewData != nil && report.OldData != nil {
			date[Modify] = append(date[Modify], report)
		}
		if report.NewData != nil && report.OldData == nil {
			date[Add] = append(date[Add], report)
		}
		if report.NewData == nil && report.OldData != nil {
			date[Delete] = append(date[Delete], report)
		}
	}
	res := make([]DiffRes, 0)
	for s, diffReports := range date {
		res = append(res, DiffRes{
			Action: s,
			Diff:   diffReports,
		})
	}

	return res, nil
}

func compareJSON(value1, value2 interface{}, key string, reports *[]DiffReport, level int, keyField, ignoreField []string) {
	//比对层级为1时，直接比对整个json是否一致
	if level <= 1 {
		v1, _ := json.Marshal(value1)
		v2, _ := json.Marshal(value2)
		if md5.Sum(v1) != md5.Sum(v2) {
			appendReport(key, value1, value2, reports)
		}
		return
	}

	//比对层级大于2及以上时，需要比对Object的每项field/Array中的每个item，
	//object/array存在ignore_field时需要过滤字段，不进行输出
	switch v1 := value1.(type) {
	/*case nil:
		if value2 != nil {
			appendReport(key, nil, value2, reports)
		}
	case bool:
		v2, ok := value2.(bool)
		if !ok || v1 != v2 {
			appendReport(key, v1, value2, reports)
		}
	case float64:
		v2, ok := value2.(float64)
		if !ok || v1 != v2 {
			appendReport(key, v1, value2, reports)
		}
	case string:
		v2, ok := value2.(string)
		if !ok || v1 != v2 {
			appendReport(key, v1, value2, reports)
		}*/
	case []interface{}:
		//数组类型时需要比对key_field
		if v2, ok := value2.([]interface{}); ok {
			compareSlice(v1, v2, key, reports, level, keyField, ignoreField)
		} else {
			appendReport(key, v1, value2, reports)
		}
	case map[string]interface{}:
		if v2, ok := value2.(map[string]interface{}); ok {
			compareMap(v1, v2, key, reports, level, keyField, ignoreField)
		} else {
			appendReport(key, v1, value2, reports)
		}
	default:
		// 不支持的数据类型
		appendReport(key, v1, value2, reports)
	}
}

func compareSlice(slice1, slice2 []interface{}, key string, reports *[]DiffReport, level int, keyField, ignoreField []string) {
	if level == 0 {
		return
	}
	if contains(ignoreField, key) {
		return
	}
	maxLen := len(slice1)
	if len(slice2) > maxLen {
		maxLen = len(slice2)
	}

	for i := 0; i < maxLen; i++ {
		if i >= len(slice1) {
			appendReport(fmt.Sprintf("%s[%d]", key, i), nil, slice2[i], reports)
		} else if i >= len(slice2) {
			appendReport(fmt.Sprintf("%s[%d]", key, i), slice1[i], nil, reports)
		} else {
			compareJSON(slice1[i], slice2[i], fmt.Sprintf("%s[%d]", key, i), reports, level-1, keyField, ignoreField)
		}
	}
}

func compareMap(map1, map2 map[string]interface{}, key string, reports *[]DiffReport, level int, keyField, ignoreField []string) {
	if level == 0 {
		return
	}

	for k, v1 := range map1 {
		if contains(ignoreField, k) {
			continue
		}
		if v2, ok := map2[k]; ok {
			compareJSON(v1, v2, joinPath(key, k), reports, level-1, keyField, ignoreField)
		} else {
			appendReport(joinPath(key, k), v1, nil, reports)
		}
	}

	for k, v2 := range map2 {
		if contains(ignoreField, k) {
			continue
		}
		if _, ok := map1[k]; !ok {
			appendReport(joinPath(key, k), nil, v2, reports)
		}
	}
}

func appendReport(key string, old, new interface{}, reports *[]DiffReport) {
	diff := DiffReport{
		Field:   key,
		OldData: old,
		NewData: new,
	}

	*reports = append(*reports, diff)
}

func compareObjects(obj1, obj2 map[string]interface{}, path string, reports *[]DiffReport) {
	for k, v1 := range obj1 {

		v2, ok := obj2[k]
		if !ok {
			// obj2 中未找到该 key，说明 obj1 中的字段被删除了
			keyPath := joinPath(path, k)
			*reports = append(*reports, DiffReport{
				Field:   keyPath,
				OldData: v1,
			})
			continue
		}
		// 比较字段值
		switch v1.(type) {
		case map[string]interface{}:
			v2Map, ok := v2.(map[string]interface{})
			if !ok {
				// 类型不匹配，说明 obj1 中的字段被更新了
				keyPath := joinPath(path, k)
				*reports = append(*reports, DiffReport{
					Field:   keyPath,
					OldData: v1,
					NewData: v2,
				})
				continue
			}
			// 递归比较子对象
			compareObjects(v1.(map[string]interface{}), v2Map, joinPath(path, k), reports)
		case []interface{}:
			v2Slice, ok := v2.([]interface{})
			if !ok {
				// 类型不匹配，说明 obj1 中的字段被更新了
				keyPath := joinPath(path, k)
				*reports = append(*reports, DiffReport{
					Field:   keyPath,
					OldData: v1,
					NewData: v2,
				})
				continue
			}
			// 比较数组元素
			compareArrays(v1.([]interface{}), v2Slice, joinPath(path, k), reports)
		default:
			if !reflect.DeepEqual(v1, v2) {
				// 值不相等，说明 obj1 中的字段被更新了
				keyPath := joinPath(path, k)
				*reports = append(*reports, DiffReport{
					Field:   keyPath,
					OldData: v1,
					NewData: v2,
				})
			}
		}
	}

	for k, v2 := range obj2 {
		if _, ok := obj1[k]; !ok {
			// obj1 中未找到该 key，说明 obj2 中的字段被新增了
			keyPath := joinPath(path, k)
			*reports = append(*reports, DiffReport{
				Field:   keyPath,
				NewData: v2,
			})
		}
	}
}

func compareArrays(arr1, arr2 []interface{}, path string, reports *[]DiffReport) {
	len1, len2 := len(arr1), len(arr2)
	for i := 0; i < len1 || i < len2; i++ {
		if i >= len1 {
			// obj1 中的元素不足，说明 obj2 中的元素被新增了
			keyPath := joinPath(path, fmt.Sprint("[", i, "]"))
			*reports = append(*reports, DiffReport{
				Field:   keyPath,
				NewData: arr2[i],
			})
		} else if i >= len2 {
			// obj2 中的元素不足，说明 obj1 中的元素被删除了
			keyPath := joinPath(path, fmt.Sprint("[", i, "]"))
			*reports = append(*reports, DiffReport{
				Field:   keyPath,
				OldData: arr1[i],
			})
		} else {
			switch arr1[i].(type) {
			case map[string]interface{}:
				v2Map, ok := arr2[i].(map[string]interface{})
				if !ok {
					// 类型不匹配，说明 obj1 中的元素被更新了
					keyPath := joinPath(path, fmt.Sprint("[", i, "]"))
					*reports = append(*reports, DiffReport{
						Field:   keyPath,
						OldData: arr1[i],
						NewData: arr2[i],
					})
					continue
				}
				// 递归比较子对象
				compareObjects(arr1[i].(map[string]interface{}), v2Map, joinPath(path, fmt.Sprint("[", i, "]")), reports)
			default:
				if !reflect.DeepEqual(arr1[i], arr2[i]) {
					// 值不相等，说明 obj1 中的元素被更新了
					keyPath := joinPath(path, fmt.Sprint("[", i, "]"))
					*reports = append(*reports, DiffReport{
						Field:   keyPath,
						OldData: arr1[i],
						NewData: arr2[i],
					})
				}
			}
		}
	}
}

func joinPath(path, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}

// 判断一个字符串是否在一个字符串数组中
func contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
