package generator

import (
	"github.com/yb19890724/go-study/gen/example2/pkg/option"
)

type Generator interface {
	Run(opt *option.Option) error
}
