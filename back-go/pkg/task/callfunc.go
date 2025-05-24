package task

import (
	"gin-web/pkg"
	"reflect"
)

//参数explain
//fn：传入的函数，可以是任何函数类型,它的类型是 interface{}，表示可以是任意类型的函数。
//args：传入的参数，使用了可变参数（...interface{}），表示可以传递任意个参数，类型是 interface{}。
//返回值explain

func CallFunc(fn interface{}, args ...interface{}) (result []interface{}, err error) {
	//捕获异常
	defer func() { err = pkg.Recover() }()
	fnType := reflect.TypeOf(fn)   //获取传入的函数的类型信息
	fnValue := reflect.ValueOf(fn) //获取传日函数的值信息
	numIn := fnType.NumIn()        //获取该函数的如参数量(函数参数个数)

	var out []reflect.Value
	if numIn == 0 {
		out = fnValue.Call(nil) //直接调用fnValue的函数(直接调用无参函数)
	} else {
		argsLength := len(args) //传入参数的数量
		argumentIn := numIn     //获取函数的参数数量
		//IsVariadic判定规则(如果最后一个参数是...T返回true)
		if fnType.IsVariadic() { //判断函数是否是可变参数函数
			argumentIn-- //如果是可变参数函数,固定参数的数量要减1
		}
		if argsLength < argumentIn {
			panic("callFunc: CallFunc with too few input arguments") //参数过少
		}
		if !fnType.IsVariadic() && argsLength > argumentIn {
			panic("callFunc: CallFunc with too many input arguments") //参数过多
		}
		in := make([]reflect.Value, numIn) // 创建一个 reflect.Value 切片，长度为 numIn
		for i := 0; i < argumentIn; i++ {
			if args[i] == nil {
				in[i] = reflect.Zero(fnType.In(i)) // 若参数为 nil，则赋值为该参数类型的零值
			} else {
				in[i] = reflect.ValueOf(args[i]) // 否则，将参数转换为 reflect.Value
			}
		}
		if fnType.IsVariadic() {
			m := argsLength - argumentIn                         //计算可变参数个数
			slice := reflect.MakeSlice(fnType.In(numIn-1), m, m) //创造可变参数的slice
			in[numIn-1] = slice
			for i := 0; i < m; i++ {
				x := args[argumentIn+i]
				if x != nil {
					slice.Index(i).Set(reflect.ValueOf(x)) //赋值
				}
			}
			out = fnValue.CallSlice(in) //调用带可变参数的函数，等价于 fn(args...)
		} else {
			out = fnValue.Call(in) //调用普通函数（非可变参数），等价于 fn(a, b)
		}

	}
	//将通过 reflect.Value.Call() 或 reflect.Value.CallSlice() 调用函数后返回的结果（out）转换为 []interface{} 类型
	if out != nil && len(out) > 0 {
		result = make([]interface{}, len(out))
		for i, v := range out {
			result[i] = v.Interface()
		}
	}
	return
}
