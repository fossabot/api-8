package types

type ApiTest struct {
	Desc       string
	Request    interface{}
	ExRespCode int
	ExRespBody interface{}
}
