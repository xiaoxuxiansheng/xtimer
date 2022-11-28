package conf

type SliceConf struct {
	BucketsNum int `yaml:"bucketsNum"`
}

type SliceConfProvider struct {
	conf *SliceConf
}

func NewSliceConfProvider(conf *SliceConf) *SliceConfProvider {
	return &SliceConfProvider{conf: conf}
}

func (s *SliceConfProvider) Get() *SliceConf {
	return s.conf
}
