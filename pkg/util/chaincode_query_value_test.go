package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQuery(t *testing.T) {
	ast := assert.New(t)

	query, _ := ParseQueryString([]byte(`{
  "selector": {
    "gender": "Male"
  }
}`))
	ast.Equal(query, `{"selector":{"gender":"Male"},"limit":20}`, "等于查询")

	query, _ = ParseQueryString([]byte(`{
  "selector": {
    "gender": "Male",
    "createTime": ">1530606578"
  }
}`))
	ast.Equal(query, `{"selector":{"createTime":{"$gt":1530606578},"gender":"Male"},"limit":20}`, "时间: 大于查询")

	query, _ = ParseQueryString([]byte(`{
  "selector": {
    "gender": "Male",
    "createTime": "<1530606578"
  }
}`))
	ast.Equal(query, `{"selector":{"createTime":{"$lt":1530606578},"gender":"Male"},"limit":20}`, "时间: 小于查询")

	query, _ = ParseQueryString([]byte(`{
  "selector": {
    "gender": "Male",
    "createTime": ">=1530606578"
  }
}`))
	ast.Equal(query, `{"selector":{"createTime":{"$gte":1530606578},"gender":"Male"},"limit":20}`, "时间: 大于等于查询")

	query, _ = ParseQueryString([]byte(`{
  "selector": {
    "gender": "Male",
    "createTime": "<=1530606578"
  }
}`))
	ast.Equal(query, `{"selector":{"createTime":{"$lte":1530606578},"gender":"Male"},"limit":20}`, "时间: 小于等于查询")

	query, _ = ParseQueryString([]byte(`{
  "selector": {
    "gender": "Male",
    "createTime": 1527012450
  }
}`))
	ast.Equal(query, `{"selector":{"createTime":1527012450,"gender":"Male"},"limit":20}`, "时间: 等于查询")

	query, _ = ParseQueryString([]byte(`{
  "selector": {
    "gender": "Male",
    "money": "20..30"
  }
}`))
	ast.Equal(query, `{"selector":{"gender":"Male","money":{"$gte":20,"$lte":30}},"limit":20}`, "时间: 等于查询")

	query, _ = ParseQueryString([]byte(`{
  "selector": {
    "gender": "Male",
    "money": ">20"
  },
  "page": 2
}`))
	ast.Equal(query, `{"selector":{"gender":"Male","money":{"$gt":20}},"limit":20,"skip":40}`, "时间: 等于查询")
}
