package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-macaron/inject"
	ginsession "github.com/go-session/gin-session"
	"log"
	"reflect"
)

var globalInjection = inject.New()

// gin下的依赖注入，由E99P1ant实现
func Warp(deps []interface{}, handlers ...interface{}) func(c *gin.Context) {
	inj := inject.New()
	inj.SetParent(globalInjection)

	return func(c *gin.Context) {
		inj.Map(c)
		inj.Map(ginsession.FromContext(c))
		for _, dep := range deps {
			inj.Map(dep)
		}
		//传参调用
		for _, handler := range handlers {
			val, err := inj.Invoke(handler)
			if err != nil {
				log.Fatalf("Failed to invoke: %v", err)
			}
			if len(val) != 0 {
				switch val[0].Interface().(type) {
				case interface{}:
					inj.Map(reflect.ValueOf(val[0].Interface()).Interface())
				}
			}
		}
	}
}
