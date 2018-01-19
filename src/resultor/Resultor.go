package resultor

import (
	"fmt"
	"net/http"
	"reflect"
	"github.com/pquerna/ffjson/ffjson"
)

func RetChanges(changes int64) string {
	return fmt.Sprintf(`{"ok":%v,"changes":%v}`, true, changes)
}

func RetOk(w http.ResponseWriter, result interface{}) {
	resValue := reflect.ValueOf(result)

	if result == nil {
		fmt.Fprint(w, RetChanges(0))
	}

	var res interface{}
	if resValue.Kind() == reflect.Array || resValue.Kind() == reflect.Slice {
		res = result
	} else {
		res = []interface{}{result}
	}
	resValue = reflect.ValueOf(res)
	bytes, _ := ffjson.Marshal(res)
	fmt.Fprintf(w,`{"ok":%v,"changes":%v,"data":%v}`, true, resValue.Len(), string(bytes))
}

func RetErr(w http.ResponseWriter, err string) {
	fmt.Fprintf(w, `{"ok":%v, "err":{"msg":"%v"}}`, false, err)
}