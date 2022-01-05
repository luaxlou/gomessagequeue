package mqengines

//对于中间件引擎的抽象，随着后续模式的增加，做灵活调整

type MqEngine interface {
	Add(key string, content string,maxLen int64) error
	ReadBlock(key string, group string, count int64, onRead func(contents []string) error)
	Read(key string, group string, count int64, onRead func(contents []string) error)
}
