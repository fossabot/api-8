package testhelper

import "encoding/json"

func (s *Suite) RespBodyEqual(body []byte, other interface{}) {
	p, err := json.Marshal(other)
	s.Nil(err)
	s.Equal(string(body), string(p))
}
