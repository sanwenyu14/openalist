package template

import (
	"github.com/AlliotTech/openalist/internal/driver"
	"github.com/AlliotTech/openalist/internal/op"
)

type Addition struct {
	// Usually one of two
	driver.RootPath
	driver.RootID
	// define other
	Field string `json:"field" type:"select" required:"true" options:"a,b,c" default:"a"`
}

var config = driver.Config{
	Name:              "Template",
	LocalSort:         false,
	OnlyLocal:         false,
	OnlyProxy:         false,
	NoCache:           false,
	NoUpload:          false,
	NeedMs:            false,
	DefaultRoot:       "root, / or other",
	CheckStatus:       false,
	Alert:             "",
	NoOverwriteUpload: false,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &Template{}
	})
}
