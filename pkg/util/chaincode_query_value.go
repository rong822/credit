package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	LessThan           = "<"
	LessThanOrEqual    = "<="
	GreaterThan        = ">"
	GreaterThanOrEqual = ">="
	MustExistIn        = ".."
)

var (
	conditionMap = map[string]string{
		GreaterThan:        "$gt",
		GreaterThanOrEqual: "$gte",
		LessThan:           "$lt",
		LessThanOrEqual:    "$lte",
		MustExistIn:        "$in",
	}
)

type Query struct {
	Selector map[string]interface{}   `json:"selector,omitempty"`
	Sort     []map[string]interface{} `json:"sort,omitempty"`
	Limit    int                      `json:"limit,omitempty"`
	Page     int                      `json:"page,omitempty"`
	Skip     int                      `json:"skip,omitempty"`
}

func NewQuery() *Query {
	return &Query{
		Limit: 20,
	}
}

// ParseCouchDBQueryString 返回 CouchDB 的查询语句
//
// Limit 默认为一页 20 个
// Page 默认为 第 1 页
// Skip 默认为 跳过的个数
//
// 普通字段查询格式:
// {
//   "selector": {
//     "key1": "value1",
//     "key2": "value2"
//   }
// }
//
// 大小查询以及范围查询格式:
// {
//   "selector": {
//     "amount": "100",                               // 查询金额 等于 100
//     "createTime1": ">1531211582",                  // 查询时间 大于 1531211582
//     "amount3": ">=100",                            // 查询金额 大于等于 100
//     "amount4": "<100",                             // 查询金额 小于 100
//     "amount5": "<=100",                            // 查询金额 小于等于 100
//     "createTime2": "1531111582..1531211582"        // 查询金额 1531111582-1531211582 的区间, 包含 1531111582 和 1531211582.
//   }
// }

func ParseQueryString(body []byte) (string, error) {
	query := NewQuery()
	if err := json.Unmarshal(body, query); err != nil {
		return "", err
	}

	for k, v := range query.Selector {
		switch v := v.(type) {
		case string:

			if strings.Contains(v, GreaterThanOrEqual) {
				value := make(map[string]interface{})
				value[conditionMap[GreaterThanOrEqual]] = convertStringToValue(strings.Replace(v, GreaterThanOrEqual, "", 1))
				query.Selector[k] = value
				continue
			}
			if strings.Contains(v, GreaterThan) {
				value := make(map[string]interface{})
				value[conditionMap[GreaterThan]] = convertStringToValue(strings.Replace(v, GreaterThan, "", 1))
				query.Selector[k] = value
				continue
			}
			if strings.Contains(v, LessThanOrEqual) {
				value := make(map[string]interface{})
				value[conditionMap[LessThanOrEqual]] = convertStringToValue(strings.Replace(v, LessThanOrEqual, "", 1))
				query.Selector[k] = value
				continue
			}
			if strings.Contains(v, LessThan) {
				value := make(map[string]interface{})
				value[conditionMap[LessThan]] = convertStringToValue(strings.Replace(v, LessThan, "", 1))
				query.Selector[k] = value
				continue
			}
			if strings.Contains(v, MustExistIn) {
				values := strings.Split(v, MustExistIn)
				value := make(map[string]interface{})
				value[conditionMap[GreaterThanOrEqual]] = convertStringToValue(values[0])
				value[conditionMap[LessThanOrEqual]] = convertStringToValue(values[1])
				query.Selector[k] = value
				continue
			}

		default:

			if strings.Contains(k, "start_") {
				key := strings.Replace(k, "start_", "", 1)
				value := make(map[string]interface{})
				value[conditionMap[GreaterThanOrEqual]] = v
				delete(query.Selector, k)
				if query.Selector[key] != nil {
					value[conditionMap[LessThanOrEqual]] = convertStringToValue(extractMapValue(query.Selector[key]))
					query.Selector[key] = value
				} else {
					query.Selector[key] = value
				}
				continue
			}

			if strings.Contains(k, "end_") {
				key := strings.Replace(k, "end_", "", 1)
				value := make(map[string]interface{})
				value[conditionMap[LessThanOrEqual]] = v
				delete(query.Selector, k)
				if query.Selector[key] != nil {
					value[conditionMap[GreaterThanOrEqual]] = convertStringToValue(extractMapValue(query.Selector[key]))
					query.Selector[key] = value
				} else {
					query.Selector[key] = value
				}
				continue
			}
		}
	}

	// 分页查询
	if query.Page != 0 {
		query.Skip += query.Page * query.Limit
		query.Page = 0
	}

	result, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func convertStringToValue(v string) interface{} {
	if strings.Contains(v, ".") {
		value, _ := strconv.ParseFloat(v, 64)
		return value
	}
	value, _ := strconv.ParseInt(v, 10, 64)
	return value
}

func extractMapValue(i interface{}) string { //extract the value from a map instantiated interface
	str := fmt.Sprintf("%v", i)
	strs := strings.Split(str, ":")
	str = strings.Replace(strs[1], "]", "", 1)
	return str
}
