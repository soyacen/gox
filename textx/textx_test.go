package textx

import (
	"testing"
)

type codec struct {
	name     string
	utf8In   string
	toUtf8   func([]byte) ([]byte, error)
	fromUtf8 func([]byte) ([]byte, error)
}

func TestCodecRoundtrips(t *testing.T) {
	cases := []codec{
		{"GBK", "你好,世界", GBKToUtf8, Utf8ToGBK},
		{"GB18030", "你好,世界", GB18030ToUtf8, Utf8ToGB18030},
		{"HZGB2312", "你好,世界", HZGB2312ToUtf8, Utf8ToHZGB2312},
		{"EUCJP", "こんにちは世界", EUCJPToUtf8, Utf8ToEUCJP},
		{"ISO2022JP", "こんにちは世界", ISO2022JPToUtf8, Utf8ToISO2022JP},
		{"ShiftJIS", "こんにちは世界", ShiftJISToUtf8, Utf8ToShiftJIS},
		{"EUCKR", "안녕하세요 세계", EUCKRToUtf8, Utf8ToEUCKR},
		{"Big5", "你好，世界", Big5ToUtf8, Utf8ToBig5},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := tc.fromUtf8([]byte(tc.utf8In))
			if err != nil {
				t.Fatalf("Utf8To%s err: %v", tc.name, err)
			}
			if len(encoded) == 0 {
				t.Fatalf("Utf8To%s produced empty output", tc.name)
			}

			decoded, err := tc.toUtf8(encoded)
			if err != nil {
				t.Fatalf("%sToUtf8 err: %v", tc.name, err)
			}
			if string(decoded) != tc.utf8In {
				t.Fatalf("%s roundtrip mismatch: got %q want %q", tc.name, decoded, tc.utf8In)
			}
		})
	}
}

func TestCodecEmptyInput(t *testing.T) {
	funcs := map[string]func([]byte) ([]byte, error){
		"GBKToUtf8":       GBKToUtf8,
		"Utf8ToGBK":       Utf8ToGBK,
		"GB18030ToUtf8":   GB18030ToUtf8,
		"Utf8ToGB18030":   Utf8ToGB18030,
		"HZGB2312ToUtf8":  HZGB2312ToUtf8,
		"Utf8ToHZGB2312":  Utf8ToHZGB2312,
		"EUCJPToUtf8":     EUCJPToUtf8,
		"Utf8ToEUCJP":     Utf8ToEUCJP,
		"ISO2022JPToUtf8": ISO2022JPToUtf8,
		"Utf8ToISO2022JP": Utf8ToISO2022JP,
		"ShiftJISToUtf8":  ShiftJISToUtf8,
		"Utf8ToShiftJIS":  Utf8ToShiftJIS,
		"EUCKRToUtf8":     EUCKRToUtf8,
		"Utf8ToEUCKR":     Utf8ToEUCKR,
		"Big5ToUtf8":      Big5ToUtf8,
		"Utf8ToBig5":      Utf8ToBig5,
	}

	for name, fn := range funcs {
		name, fn := name, fn
		t.Run(name, func(t *testing.T) {
			out, err := fn(nil)
			if err != nil {
				t.Fatalf("%s(nil) err: %v", name, err)
			}
			if len(out) != 0 {
				t.Fatalf("%s(nil) got %v, want empty", name, out)
			}

			out, err = fn([]byte{})
			if err != nil {
				t.Fatalf("%s([]) err: %v", name, err)
			}
			if len(out) != 0 {
				t.Fatalf("%s([]) got %v, want empty", name, out)
			}
		})
	}
}

func TestCodecASCIIPassthrough(t *testing.T) {
	ascii := []byte("hello world 123")
	codecs := []codec{
		{"GBK", "", GBKToUtf8, Utf8ToGBK},
		{"GB18030", "", GB18030ToUtf8, Utf8ToGB18030},
		{"EUCJP", "", EUCJPToUtf8, Utf8ToEUCJP},
		{"ShiftJIS", "", ShiftJISToUtf8, Utf8ToShiftJIS},
		{"EUCKR", "", EUCKRToUtf8, Utf8ToEUCKR},
		{"Big5", "", Big5ToUtf8, Utf8ToBig5},
	}

	for _, c := range codecs {
		c := c
		t.Run(c.name, func(t *testing.T) {
			encoded, err := c.fromUtf8(ascii)
			if err != nil {
				t.Fatalf("Utf8To%s err: %v", c.name, err)
			}
			decoded, err := c.toUtf8(encoded)
			if err != nil {
				t.Fatalf("%sToUtf8 err: %v", c.name, err)
			}
			if string(decoded) != string(ascii) {
				t.Fatalf("%s ascii passthrough mismatch: got %q want %q", c.name, decoded, ascii)
			}
		})
	}
}
