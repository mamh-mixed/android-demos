package util

import "testing"

type TestCommon struct {
	InCommon string
}

type TestStruct struct {
	TestCommon
	InStruct int
}

func TestQuery(t *testing.T) {
	s := TestStruct{
		TestCommon: TestCommon{
			InCommon: "test",
		},
		InStruct: 1,
	}

	buf, err := Query(s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(buf.String())
}
