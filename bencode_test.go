package bencode

import (
	"reflect"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	cases := []struct {
		in  int
		out string
	}{
		{1, "i1e"},
		{0, "i0e"},
		{42, "i42e"},
		{-42, "i-42e"},
	}

	for idx := 0; idx != len(cases); idx++ {
		if str, err := encodeInt(cases[idx].in); err == nil {
			if str != cases[idx].out {
				t.Fatal(idx, str)
			}
		} else {
			t.Fatal(idx, err)
		}
	}
}
func TestEncodeString(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"spam", "4:spam"},
		{"hello,中国", "12:hello,中国"},
	}

	for idx := 0; idx != len(cases); idx++ {
		if str, err := encodeString(cases[idx].in); err == nil {
			if str != cases[idx].out {
				t.Fatal(idx, str)
			}
		} else {
			t.Fatal(idx, err)
		}
	}
}

func TestEncodeSlice(t *testing.T) {
	cases := []struct {
		in  []interface{}
		out string
	}{
		{[]interface{}{1, "spam"}, "li1e4:spame"},
		{[]interface{}{1, "spam", -1}, "li1e4:spami-1ee"},
		{[]interface{}{1, "spam", []int{1, 2}}, "li1e4:spamli1ei2eee"},
	}

	for idx := 0; idx != len(cases); idx++ {
		v := reflect.Indirect(reflect.ValueOf(cases[idx].in))
		if str, err := encodeSlice(v); err == nil {
			if str != cases[idx].out {
				t.Fatal(idx, str)
			}
		} else {
			t.Fatal(idx, err)
		}
	}
}

func TestEncodeMap(t *testing.T) {
	in := make(map[string]interface{})
	in["q"] = "ping"
	in["id"] = "identify"
	in["t"] = 123
	in["list"] = []string{"abc", "def"}

	out := "d2:id8:identify4:listl3:abc3:defe1:q4:ping1:ti123ee"
	v := reflect.Indirect(reflect.ValueOf(in))

	if str, err := encodeMap(v); err == nil {
		if str != out {
			t.Fatal(str)
		}
	} else {
		t.Fatal(err)
	}
}

func TestEncodeStruct(t *testing.T) {
	in := struct {
		Q    string   `json:"q"`
		Id   string   `json:"id"`
		T    int      `json:"t"`
		List []string `json:"list"`
	}{
		"ping",
		"identify",
		123,
		[]string{"abc", "def"},
	}

	out := "d2:id8:identify4:listl3:abc3:defe1:q4:ping1:ti123ee"
	v := reflect.Indirect(reflect.ValueOf(in))

	if str, err := encodeStruct(v); err == nil {
		if str != out {
			t.Fatal(str)
		}
	} else {
		t.Fatal(err)
	}
}

func TestEncodeTop(t *testing.T) {
	in1 := []struct {
		Q    string   `json:"q"`
		Id   string   `json:"id"`
		T    int      `json:"t"`
		List []string `json:"list"`
	}{
		{
			"ping",
			"identify",
			123,
			[]string{"abc", "def"}},
		{
			"r",
			"who",
			321,
			[]string{"rst", "xyz"}},
	}
	out1 := "ld2:id8:identify4:listl3:abc3:defe1:q4:ping1:ti123eed2:id3:who4:listl3:rst3:xyze1:q1:r1:ti321eee"

	cases := []struct {
		in  interface{}
		out string
	}{
		{42, "i42e"},
		{"hello,中国", "12:hello,中国"},
		{[]interface{}{1, "spam", []int{1, 2}}, "li1e4:spamli1ei2eee"},
		{in1, out1},
	}

	for idx := 0; idx != len(cases); idx++ {
		if str, err := Encode(cases[idx].in); err == nil {
			if str != cases[idx].out {
				t.Fatal(idx, str)
			}
		} else {
			t.Fatal(idx, err)
		}
	}
}
func TestEncodeComplex(t *testing.T) {
	type testdecodest struct {
		Q  string         `json:"q"`
		Id string         `json:"id"`
		T  string         `json:"t"`
		W  int            `json:"who"`
		A  int            `json:"age"`
		L  []interface{}  `json:"lt"`
		M  map[string]int `json:"mp"`
		S  struct {
			When  int    `json:"when"`
			Where string `json:"where"`
		} `json:"embed"`
	}
	in := testdecodest{
		"ping",
		"identify",
		"123",
		42,
		36,
		[]interface{}{1, "items"},
		map[string]int{"peoples": 1, "citys": 2},
		struct {
			When  int    `json:"when"`
			Where string `json:"where"`
		}{3, "where"},
	}
	expect := "d3:agei36e5:embedd4:wheni3e5:where5:wheree2:id8:identify2:ltli1e5:itemse2:mpd5:citysi2e7:peoplesi1ee1:q4:ping1:t3:1233:whoi42ee"

	out, err := Encode(in)
	if err != nil || out != expect {
		t.Fatal(expect, out)
	}
	t.Log(out, err)
}

func TestFindFirstNode(t *testing.T) {
	in := "li1e4:spamli1ei2eee"
	expectid, expectend := bencode_type_list, len(in)-1
	id, end, err := findFirstNode([]byte(in), 0)
	if err != nil {
		t.Fatal(err)
	}
	if expectend != end || expectid != id {
		t.Fatal(id, end)
	}
}

func TestDecode(t *testing.T) {
	in := "12:hello,中国"
	expect := "hello,中国"

	out := ""
	err := Decode([]byte(in), &out)
	if err != nil || expect != out {
		t.Fatal(out, err, expect)
	}

	in1 := "i42e"
	expect1 := uint32(42)

	out1 := uint32(0)
	err = Decode([]byte(in1), &out1)
	if err != nil || expect1 != out1 {
		t.Fatal(out1, err, expect1)
	}

	in2 := "li42ei36ee"
	expect2 := []int{42, 36}

	out2 := []int{}
	err = Decode([]byte(in2), &out2)

	if err != nil {
		t.Fatal(out2, err)
	}

	if reflect.DeepEqual(out2, expect2) == false {
		t.Fatal(expect2, out2)
	}

	in3 := "l12:hello,中国4:spame"
	expect3 := []string{"hello,中国", "spam"}

	out3 := []string{}
	err = Decode([]byte(in3), &out3)

	if err != nil {
		t.Fatal(out3, err)
	}

	if reflect.DeepEqual(out3, expect3) == false {
		t.Fatal(expect3, out3)
	}

	in4 := "d3:bar4:spam3:foo3:abce"
	expect4 := map[string]string{"bar": "spam", "foo": "abc"}

	//	out4 := make(map[string]string)
	out4 := map[string]string{}
	err = Decode([]byte(in4), &out4)

	if err != nil {
		t.Fatal(out4, err)
	}
	t.Log(out4)

	if reflect.DeepEqual(out4, expect4) == false {
		t.Fatal(expect4, out4)
	}

	in5 := "d2:id8:identify1:q4:ping1:t3:123e"
	type testdecodest struct {
		Q  string `json:"q"`
		Id string `json:"id"`
		T  string `json:"t"`
	}
	expect5 := testdecodest{
		"ping",
		"identify",
		"123",
	}
	out5 := testdecodest{}
	err = Decode([]byte(in5), &out5)

	if err != nil {
		t.Fatal(out5, err)
	}
	t.Log(out5)

	if reflect.DeepEqual(out5, expect5) == false {
		t.Fatal(expect5, out5)
	}
}

func TestDecodeCombin(t *testing.T) {
	in2 := "li42ei36e4:spame"
	expect2 := []interface{}{42, 36, "spam"}

	out2 := []interface{}{}
	err := Decode([]byte(in2), &out2)

	if err != nil {
		t.Fatal(out2, err)
	}

	if reflect.DeepEqual(out2, expect2) == false {
		t.Fatal(expect2, out2)
	}
	t.Log(expect2, out2)

	in3 := "l12:hello,中国4:spami42ei36ee"
	expect3 := []interface{}{"hello,中国", "spam", 42, 36}

	out3 := []interface{}{}
	err = Decode([]byte(in3), &out3)

	if err != nil {
		t.Fatal(out3, err)
	}

	if reflect.DeepEqual(out3, expect3) == false {
		t.Fatal(expect3, out3)
	}
	t.Log(expect3, out3)

	in4 := "d3:bar4:spam3:foo3:abc3:whoi42e3:agei36ee"
	expect4 := map[string]interface{}{"bar": "spam", "foo": "abc", "who": 42, "age": 36}

	out4 := map[string]interface{}{}
	err = Decode([]byte(in4), &out4)

	if err != nil {
		t.Fatal(out4, err)
	}

	if reflect.DeepEqual(out4, expect4) == false {
		t.Fatal(expect4, out4)
	}
	t.Log(expect4, out4)

	in5 := "d2:id8:identify1:q4:ping1:t3:1233:whoi42e3:agei36e"
	type testdecodest struct {
		Q  string `json:"q"`
		Id string `json:"id"`
		T  string `json:"t"`
		W  int    `json:"who"`
		A  int    `json:"age"`
	}
	expect5 := testdecodest{
		"ping",
		"identify",
		"123",
		42,
		36,
	}
	out5 := testdecodest{}
	err = Decode([]byte(in5), &out5)

	if err != nil {
		t.Fatal(out5, err)
	}

	if reflect.DeepEqual(out5, expect5) == false {
		t.Fatal(expect5, out5)
	}
	t.Log(expect5, out5)
}

func TestDecodeComplex(t *testing.T) {
	in5 := "d3:agei36e5:embedd4:wheni3e5:where5:wheree2:id8:identify2:ltli1e5:itemse2:mpd5:citysi2e7:peoplesi1ee1:q4:ping1:t3:1233:whoi42ee"
	type testdecodest struct {
		Q  string         `json:"q"`
		Id string         `json:"id"`
		T  string         `json:"t"`
		W  int            `json:"who"`
		A  int            `json:"age"`
		L  []interface{}  `json:"lt"`
		M  map[string]int `json:"mp"`
		S  struct {
			When  int    `json:"when"`
			Where string `json:"where"`
		} `json:"embed"`
	}
	expect5 := testdecodest{
		"ping",
		"identify",
		"123",
		42,
		36,
		[]interface{}{1, "items"},
		map[string]int{"peoples": 1, "citys": 2},
		struct {
			When  int    `json:"when"`
			Where string `json:"where"`
		}{3, "where"},
	}
	out5 := testdecodest{}
	err := Decode([]byte(in5), &out5)

	if err != nil {
		t.Fatal(out5, err)
	}

	if reflect.DeepEqual(out5, expect5) == false {
		t.Fatal(expect5, out5)
	}
	t.Log(expect5, out5)
}
