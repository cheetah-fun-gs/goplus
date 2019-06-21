package mgo

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// ErrorBulkParam 参数错误
var (
	ErrorBulkParam = fmt.Errorf("bulk param error")
)

// BulkKeyWildcard 通配Key 特殊用法
const (
	BulkKeyWildcard = "*"
)

// BulkDataType 数据类型
type BulkDataType int

// 数据类型 定义
const (
	BulkDataTypeDefault = iota
	BulkDataTypeInt
)

// BulkUpdateMethod 更新方法
type BulkUpdateMethod int

// 操作方法 定义
const (
	BulkUpdateMethodSet = iota
	BulkUpdateMethodDel
	BulkUpdateMethodInc
)

// BulkItem 存储项
type BulkItem struct {
	Key   string
	Field string
	Value interface{}
}

// BulkUpdateItem 更新项
type BulkUpdateItem struct {
	*BulkItem
	Type   BulkDataType
	Method BulkUpdateMethod
}

func bulkParseResult1(in []*bson.M, out []*BulkItem) {
	for _, item := range out {
		var isMatch bool
		for _, v := range in {
			if isMatch {
				break
			}
			vv := *v
			key := vv["_id"]
			for field, value := range vv {
				if field == "_id" {
					continue
				}
				if isMatch {
					break
				}
				if key.(string) == item.Key && field == item.Field {
					item.Value = value
					isMatch = true
				}
			}
		}
	}
	return
}

func bulkParseResult2(in []*bson.M) (out []*BulkItem) {
	out = []*BulkItem{}
	for _, v := range in {
		vv := *v
		key := vv["_id"]
		for field, value := range vv {
			if field == "_id" {
				continue
			}
			out = append(out, &BulkItem{
				Key:   key.(string),
				Field: field,
				Value: value,
			})
		}
	}
	return
}

// BulkFind 根据指定 key 和 field 获取数据
func BulkFind(mgoDB *mgo.Database, collection string, items []*BulkItem) error {
	keys := []string{}
	selector := bson.M{}
	for _, item := range items {
		if item.Key == "" || item.Key == BulkKeyWildcard || item.Field == "" {
			return ErrorBulkParam
		}

		keys = append(keys, item.Key)
		selector[item.Field] = 1
	}
	if len(keys) == 0 || len(selector) == 0 {
		return nil
	}

	query := bson.M{"_id": bson.M{"$in": keys}}

	result, err := daoBulkFind(mgoDB, collection, query, selector)
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}

	bulkParseResult1(result, items)
	return nil
}

// BulkFindByKeys 根据指定 key 获取数据
func BulkFindByKeys(mgoDB *mgo.Database, collection string, keys []string) ([]*BulkItem, error) {
	if len(keys) == 0 {
		return nil, ErrorBulkParam
	}

	query := bson.M{"_id": bson.M{"$in": keys}}

	result, err := daoBulkFind(mgoDB, collection, query, nil)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []*BulkItem{}, nil
	}

	return bulkParseResult2(result), nil
}

func daoBulkFind(mgoDB *mgo.Database, collection string, query, selector bson.M) ([]*bson.M, error) {
	var err error
	result := []*bson.M{}
	if selector != nil && len(selector) != 0 {
		err = mgoDB.C(collection).Find(query).Select(selector).All(&result)
	} else {
		err = mgoDB.C(collection).Find(query).All(&result)
	}

	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}
	if err == mgo.ErrNotFound {
		return []*bson.M{}, nil
	}
	return result, nil
}

type bulkUpdateParam struct {
	Key    string
	Query  bson.M
	Update bson.M
}

func bulkParseUpdate(in []*BulkUpdateItem) ([]*bulkUpdateParam, error) {
	cache := map[string]*bulkUpdateParam{}

	var isHasWildcardKey bool

	for _, item := range in {
		if item.Key == "" || item.Field == "" {
			return nil, ErrorBulkParam
		}

		if item.Key == BulkKeyWildcard {
			isHasWildcardKey = true
		}

		if isHasWildcardKey && item.Key != BulkKeyWildcard {
			return nil, ErrorBulkParam
		}

		param, ok := cache[item.Key]
		if !ok {
			param = &bulkUpdateParam{
				Key:    item.Key,
				Query:  bson.M{},
				Update: bson.M{},
			}
			cache[item.Key] = param
		}

		if item.Type == BulkDataTypeDefault && item.Method == BulkUpdateMethodInc {
			return nil, ErrorBulkParam
		}
		if item.Method == BulkUpdateMethodInc {
			if _, ok := item.Value.(int); !ok {
				return nil, ErrorBulkParam
			}
		}

		switch item.Method {
		case BulkUpdateMethodSet:
			set, ok := param.Update["$set"]
			if ok {
				set.(bson.M)[item.Field] = item.Value
			} else {
				param.Update["$set"] = bson.M{item.Field: item.Value}
			}
		case BulkUpdateMethodDel:
			unset, ok := param.Update["$unset"]
			if ok {
				unset.(bson.M)[item.Field] = ""
			} else {
				param.Update["$unset"] = bson.M{item.Field: ""}
			}
		case BulkUpdateMethodInc:
			inc, ok := param.Update["$inc"]
			if ok {
				inc.(bson.M)[item.Field] = item.Value
			} else {
				param.Update["$inc"] = bson.M{item.Field: item.Value}
			}
			if item.Value.(int) < 0 {
				field, ok := param.Query[item.Field]
				if ok {
					field.(bson.M)["$gte"] = item.Value.(int) * -1
				} else {
					param.Query[item.Field] = bson.M{"$gte": item.Value.(int) * -1}
				}
			}
		}
	}

	out := []*bulkUpdateParam{}
	for _, param := range cache {
		out = append(out, param)
	}
	return out, nil
}

// BulkUpdateSingleKey 单key更新 保证原子性
func BulkUpdateSingleKey(mgoDB *mgo.Database, collection string, updateItems []*BulkUpdateItem) ([]*BulkItem, error) {
	params, err := bulkParseUpdate(updateItems)
	if err != nil {
		return nil, err
	}

	if len(params) != 1 {
		return nil, ErrorBulkParam
	}

	param := params[0]
	if param.Key == BulkKeyWildcard {
		return nil, ErrorBulkParam
	}

	result, err := daoBulkUpdate(mgoDB, collection, []string{param.Key}, param.Query, param.Update)
	if err != nil {
		return nil, err
	}

	items := []*BulkItem{}
	for _, updateItem := range updateItems {
		items = append(items, updateItem.BulkItem)
	}
	bulkParseResult1(result, items)
	return items, nil
}

// BulkUpdateMultiKey 多key同时更新 保证原子性
func BulkUpdateMultiKey(mgoDB *mgo.Database, collection string, keys []string, updateItems []*BulkUpdateItem) ([]*BulkItem, error) {
	params, err := bulkParseUpdate(updateItems)
	if err != nil {
		return nil, err
	}

	if len(params) != 1 {
		return nil, ErrorBulkParam
	}

	param := params[0]
	if param.Key != BulkKeyWildcard {
		return nil, ErrorBulkParam
	}

	result, err := daoBulkUpdate(mgoDB, collection, keys, param.Query, param.Update)
	if err != nil {
		return nil, err
	}

	items := []*BulkItem{}
	for _, key := range keys {
		for _, updateItem := range updateItems {
			updateItem.Key = key
			items = append(items, updateItem.BulkItem)
		}
	}

	bulkParseResult1(result, items)
	return items, nil
}

// BulkUpdate 多key循环更新 不能保证原子性
func BulkUpdate(mgoDB *mgo.Database, collection string, updateItems []*BulkUpdateItem) ([]*BulkItem, error) {
	params, err := bulkParseUpdate(updateItems)
	if err != nil {
		return nil, err
	}

	result := []*bson.M{}
	for _, param := range params {
		if param.Key == BulkKeyWildcard {
			return nil, ErrorBulkParam
		}

		r, err := daoBulkUpdate(mgoDB, collection, []string{param.Key}, param.Query, param.Update)
		if err != nil {
			continue
		}
		result = append(result, r...)
	}

	items := []*BulkItem{}
	for _, updateItem := range updateItems {
		items = append(items, updateItem.BulkItem)
	}

	bulkParseResult1(result, items)
	return items, nil
}

func daoBulkUpdate(mgoDB *mgo.Database, collection string, keys []string, query, update bson.M) ([]*bson.M, error) {
	if len(keys) == 0 {
		return nil, mgo.ErrNotFound // TODO:
	}
	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	if query != nil && len(query) != 0 {
		if len(keys) == 1 {
			query["_id"] = keys[0]
		} else {
			query["_id"] = bson.M{"$in": keys}
		}
	} else {
		change.Upsert = true
		if len(keys) == 1 {
			query = bson.M{"_id": keys[0]}
		} else {
			query = bson.M{"_id": bson.M{"$in": keys}}
		}
	}

	result := []*bson.M{}
	_, err := mgoDB.C(collection).Find(query).Apply(change, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
