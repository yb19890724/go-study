package generator

import (
	"fmt"
	"sync"
)

var GeneratorMgr *GeneratorManger

type GeneratorManger struct {
	generatorMap sync.Map
}

func init() {

	GeneratorMgr = &GeneratorManger{
		generatorMap: sync.Map{},
	}

}

// 注册生成器
func Register(name string, gen Generator) (err error) {

	if _, ok := GeneratorMgr.generatorMap.Load(name); ok {
		err = fmt.Errorf("generator %s is exists", name)
	}

	GeneratorMgr.generatorMap.Store(name, gen)

	return
}
