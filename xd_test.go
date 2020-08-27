package main

import "testing"

func TestHexdata(t *testing.T) {
	reply := "hexdata(%q, 16, %v) is\n%q\nwant\n%q"

	cases := []struct {
		nogo uint8		/* no of grouped octets */
		s string
	}{
		{1, "./0123456789:;<="},
		{2, "./0123456789:;<="},
		{2, "10"},
		{4, "./0123456789:;<="},
	}

	want := []string {
		"2e 2f 30 31 32 33 34 35  36 37 38 39 3a 3b 3c 3d",
		"2e2f 3031 3233 3435  3637 3839 3a3b 3c3d",
		"3130                                    ",
		"2e2f3031 32333435  36373839 3a3b3c3d",
	}

	for i, in := range cases {
		got := hexdata([]byte(in.s), 16, in.nogo)
		if string(got) != want[i] {
			t.Errorf(reply, in.s, in.nogo, got, want[i])
		}
	}
}

func TestHexString(t *testing.T) {
	var off uint32 = 0
	reply := "hex_string(%q) is\n%q\nwant\n%q"

	cases := []string {
		"./0123456789:;<=",
		"fuck.",
	}

	want := []string {
"00000000: 2e2f 3031 3233 3435  3637 3839 3a3b 3c3d  ./0123456789:;<=\n",
"00000010: 6675 636b 2e                              fuck.\n",
	}

	for i, in := range cases {
		got := hex_string([]byte(in), off)
		off += uint32(len(in))

		if got != want[i] {
			t.Errorf(reply, in, got, want[i])
		}
	}
}

