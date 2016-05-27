package ela

import (
	// "fmt"
	"github.com/codegangsta/inject"
	"github.com/gogather/com/log"
	"reflect"
)

type injection struct {
	middlewares     map[reflect.Type]interface{}
	middlewaresName []reflect.Type
	injector        inject.Injector
}

func newInjection() *injection {
	inj := injection{middlewares: make(map[reflect.Type]interface{}), injector: inject.New()}
	return &inj
}

func (inj *injection) injectFuc(fun interface{}) ([]reflect.Value, error) {
	for i := 0; i < len(inj.middlewaresName); i++ {
		element := inj.middlewaresName[i]
		if element.Kind() != reflect.Func {
			inj.injector.Map(inj.middlewares[element])
		}
	}
	return inj.injector.Invoke(fun)
}

func (inj *injection) appendMiddleware(middleware interface{}) {
	t := reflect.TypeOf(middleware)
	inj.middlewares[t] = middleware
	for i := 0; i < len(inj.middlewaresName); i++ {
		element := inj.middlewaresName[i]
		if element == t {
			return
		}
	}
	inj.middlewaresName = append(inj.middlewaresName, t)
}

func (inj *injection) headMiddleware(middleware interface{}) {
	var mids []reflect.Type
	midName := reflect.TypeOf(middleware)
	mids = append(mids, midName)

	inj.middlewares[midName] = middleware

	for i := 0; i < len(inj.middlewaresName); i++ {
		element := inj.middlewaresName[i]
		if element != midName {
			mids = append(mids, element)
		}
	}

	inj.middlewaresName = mids
}

func (inj *injection) parseInjectionObjects() {
	var mids []reflect.Type
	var midw interface{}

	for i := 0; i < len(inj.middlewaresName); i++ {
		mid := inj.middlewaresName[i]
		if mid.Kind() == reflect.Func {
			// get the finnal object as middleware object
			midw = inj.parseInjectionObject(inj.middlewares[mid])
			inj.appendMiddleware(midw)

			if midw != nil {
				mids = append(mids, mid)
			}
		} else {
			inj.appendMiddleware(inj.middlewares[mid])
			mids = append(mids, mid)
		}

	}

	inj.middlewaresName = mids
}

func (inj *injection) parseInjectionObject(mid interface{}) interface{} {
	t := reflect.TypeOf(mid)
	if t.Kind() == reflect.Func {
		result, err := inj.injectFuc(mid)
		if err != nil {
			log.Redf("injection failed: %s :L51\n", err)
			return nil
		} else {
			mid = result[0]
			return inj.parseInjectionObject(mid)
		}
	} else {
		return mid
	}
}
