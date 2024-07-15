package cols

import (
	"reflect"
)

func Make(data interface{}) []map[string]interface{} {
	//v := reflect.ValueOf(data)
	//t := v.Type()
	//
	//for i := 0; i < v.NumField(); i++ {
	//	field := v.Field(i)
	//	tag := t.Field(i).Tag.Get("validate")
	//	// 解析标签中的函数名称和参数
	//	funcName, funcParams := parseTag(tag)
	//
	//	// 根据函数名称调用相应的函数
	//	if funcName != "" {
	//		funcValue := reflect.ValueOf(validateFuncs[funcName])
	//		if funcValue.IsValid() {
	//			params := make([]reflect.Value, 0)
	//			params = append(params, field)
	//			params = append(params, funcParams...)
	//			result := funcValue.Call(params)
	//			// 处理函数返回值
	//			if len(result) > 0 {
	//				if !result[0].Bool() {
	//					fmt.Printf("Validation failed for field %s\n", t.Field(i).Name)
	//				}
	//			}
	//		}
	//	}
	//}
	return nil
}

// 解析标签中的函数名称和参数
func parseTag(tag string) (string, []reflect.Value) {
	// 在实际应用中，您可能需要根据标签的格式进行更复杂的解析
	// 这里只是一个简单的示例
	if tag != "" {
		funcName := tag
		funcParams := make([]reflect.Value, 0)
		return funcName, funcParams
	}

	return "", nil
}
