// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

type Meta struct {
	Corp       string
	EndPoint   string
	Model      string
	APIVersion string
	APIKey     string
	SecretKey  string

	ChoosableThinking bool //是否可选思考模式
	EnabledThinking   bool //是否开启思考模式

	Region string //for volcengine and tencent

	APPID string //for spark
}

type Adaptor struct {
	meta Meta
}

func (a *Adaptor) Init(meta Meta) {
	a.meta = meta
}
