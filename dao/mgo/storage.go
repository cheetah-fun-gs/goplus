package mgo

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// ErrorStorageParam 参数错误
var (
	ErrorStorageParam = fmt.Errorf("Storage param error")
)

// StorageKeyWildcard 通配Key 仅有 StorageUpdateMultiKey 才会使用
const (
	StorageKeyWildcard = "*"
)

// StorageUpdateMethod 更新方法
type StorageUpdateMethod int

// 操作方法 定义
const (
	StorageUpdateMethodGet = iota // 虽然不想修改 但是想返回最新的值
	StorageUpdateMethodSet
	StorageUpdateMethodDel
	StorageUpdateMethodInc
)

// StorageItem 存储项
type StorageItem struct {
	Key   string
	Field string
	Value interface{} // 为 nil 表示不存在该 field
}

// StorageUpdateItem 更新项
type StorageUpdateItem struct {
	Key    string
	Field  string
	Value  interface{}
	Method StorageUpdateMethod
}

func storageParseResult1(in []*bson.M, out []*StorageItem) {
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

func storageParseResult2(in []*bson.M) (out []*StorageItem) {
	out = []*StorageItem{}
	for _, v := range in {
		vv := *v
		key := vv["_id"]
		for field, value := range vv {
			if field == "_id" {
				continue
			}
			out = append(out, &StorageItem{
				Key:   key.(string),
				Field: field,
				Value: value,
			})
		}
	}
	return
}

// StorageFind 根据指定 key 和 field 获取数据
func StorageFind(mgoDB *mgo.Database, collection string, items []*StorageItem) error {
	keys := []string{}
	selector := bson.M{}
	for _, item := range items {
		if item.Key == "" || item.Key == StorageKeyWildcard || item.Field == "" {
			return ErrorStorageParam
		}

		keys = append(keys, item.Key)
		selector[item.Field] = 1
	}
	if len(keys) == 0 || len(selector) == 0 {
		return nil
	}

	query := bson.M{"_id": bson.M{"$in": keys}}

	result, err := daoStorageFind(mgoDB, collection, query, selector)
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}

	storageParseResult1(result, items)
	return nil
}

// StorageFindByKeys 根据指定 key 获取数据
func StorageFindByKeys(mgoDB *mgo.Database, collection string, keys []string) ([]*StorageItem, error) {
	if len(keys) == 0 {
		return nil, ErrorStorageParam
	}

	query := bson.M{"_id": bson.M{"$in": keys}}

	result, err := daoStorageFind(mgoDB, collection, query, nil)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []*StorageItem{}, nil
	}

	return storageParseResult2(result), nil
}

func daoStorageFind(mgoDB *mgo.Database, collection string, query, selector bson.M) ([]*bson.M, error) {
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

type storageUpdateByKey struct {
	Key    string
	Query  bson.M
	Update bson.M
}

func storageParseUpdateParams(in []*StorageUpdateItem) ([]*storageUpdateByKey, error) {
	m := map[string]*storageUpdateByKey{}

	var isHasWildcardKey bool

	for _, item := range in {
		if item.Key == "" || item.Field == "" {
			return nil, ErrorStorageParam
		}

		if item.Key == StorageKeyWildcard {
			isHasWildcardKey = true
		}

		if isHasWildcardKey && item.Key != StorageKeyWildcard {
			return nil, ErrorStorageParam
		}

		param, ok := m[item.Key]
		if !ok {
			param = &storageUpdateByKey{
				Key:    item.Key,
				Query:  bson.M{},
				Update: bson.M{},
			}
			m[item.Key] = param
		}

		if item.Method == StorageUpdateMethodInc {
			if _, ok := item.Value.(int); !ok {
				return nil, ErrorStorageParam
			}
		}

		switch item.Method {
		case StorageUpdateMethodSet:
			set, ok := param.Update["$set"]
			if ok {
				set.(bson.M)[item.Field] = item.Value
			} else {
				param.Update["$set"] = bson.M{item.Field: item.Value}
			}
		case StorageUpdateMethodDel:
			unset, ok := param.Update["$unset"]
			if ok {
				unset.(bson.M)[item.Field] = ""
			} else {
				param.Update["$unset"] = bson.M{item.Field: ""}
			}
		case StorageUpdateMethodInc:
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

	out := []*storageUpdateByKey{}
	for _, param := range m {
		out = append(out, param)
	}
	return out, nil
}

// StorageUpdateSingleKey 单key更新 保证原子性
func StorageUpdateSingleKey(mgoDB *mgo.Database, collection string, updateItems []*StorageUpdateItem) ([]*StorageItem, error) {
	params, err := storageParseUpdateParams(updateItems)
	if err != nil {
		return nil, err
	}

	if len(params) != 1 {
		return nil, ErrorStorageParam
	}

	param := params[0]
	if param.Key == StorageKeyWildcard {
		return nil, ErrorStorageParam
	}

	result, err := daoStorageUpdate(mgoDB, collection, []string{param.Key}, param.Query, param.Update)
	if err != nil {
		return nil, err
	}

	items := []*StorageItem{}
	for _, updateItem := range updateItems {
		items = append(items, &StorageItem{
			Key:   updateItem.Key,
			Field: updateItem.Field,
		})
	}
	storageParseResult1(result, items)
	return items, nil
}

// StorageUpdateMultiKey 多key同时更新 保证原子性
func StorageUpdateMultiKey(mgoDB *mgo.Database, collection string, keys []string, updateItems []*StorageUpdateItem) ([]*StorageItem, error) {
	params, err := storageParseUpdateParams(updateItems)
	if err != nil {
		return nil, err
	}

	if len(params) != 1 {
		return nil, ErrorStorageParam
	}

	param := params[0]
	if param.Key != StorageKeyWildcard {
		return nil, ErrorStorageParam
	}

	result, err := daoStorageUpdate(mgoDB, collection, keys, param.Query, param.Update)
	if err != nil {
		return nil, err
	}

	items := []*StorageItem{}
	for _, key := range keys {
		for _, updateItem := range updateItems {
			updateItem.Key = key
			items = append(items, &StorageItem{
				Key:   updateItem.Key,
				Field: updateItem.Field,
			})
		}
	}
	storageParseResult1(result, items)
	return items, nil
}

// StorageUpdate 多key循环更新 不能保证原子性 可以反复使用 failure 重试
func StorageUpdate(mgoDB *mgo.Database, collection string, updateItems []*StorageUpdateItem) (
	returnNew []*StorageItem, success, failure []*StorageUpdateItem, err error) {
	params, err := storageParseUpdateParams(updateItems)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, param := range params {
		if param.Key == StorageKeyWildcard {
			return nil, nil, nil, ErrorStorageParam
		}
	}

	successKeys := map[string]bool{}
	result := []*bson.M{}
	for _, param := range params {
		r, err := daoStorageUpdate(mgoDB, collection, []string{param.Key}, param.Query, param.Update)
		if err != nil {
			continue
		}
		successKeys[param.Key] = true
		result = append(result, r...)
	}

	success, failure = []*StorageUpdateItem{}, []*StorageUpdateItem{}
	returnNew = []*StorageItem{}
	for _, updateItem := range updateItems {
		if _, ok := successKeys[updateItem.Key]; ok {
			returnNew = append(returnNew, &StorageItem{
				Key:   updateItem.Key,
				Field: updateItem.Field,
			})
			success = append(success, updateItem)
			continue
		}
		failure = append(failure, updateItem)
	}

	storageParseResult1(result, returnNew)
	return
}

func daoStorageUpdate(mgoDB *mgo.Database, collection string, keys []string, query, update bson.M) ([]*bson.M, error) {
	if len(keys) == 0 {
		return nil, ErrorStorageParam
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

	var err error
	var result []*bson.M

	if len(keys) == 1 {
		r := bson.M{}
		_, err = mgoDB.C(collection).Find(query).Apply(change, &r)
		result = append(result, &r)
	} else {
		r := []*bson.M{}
		_, err = mgoDB.C(collection).Find(query).Apply(change, &r)
		result = r
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}
