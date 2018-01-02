package examples

import (
	"testing"
	"github.com/chemikadze/asn1go/internal/utils"
	"os"
	"encoding/asn1"
	"fmt"
)

var (
	_ = os.Args
)

//go:generate asn1go -package examples rfc4120.asn1 rfc4120_generated.go

func TestMessagesDeclared(t *testing.T) {
	var (
		_ KDC_REQ
		_ KDC_REP
		_ AS_REQ
		_ AS_REP
		_ AP_REQ
		_ AP_REP
	)
}


// TODO: Failed to parse: asn1: structure error: tags don't match (16 vs {class:0 tag:3 length:5 isCompound:false}) {optional:false explicit:true application:false defaultValue:<nil> tag:0x820306790 stringType:0 timeType:0 set:false omitEmpty:false} KDCOptions @4
func disabled_TestKdcReq(t *testing.T) {
	msgBytes := utils.ParseWiresharkHex(`
0000   30 81 aa a1 03 02 01 05 a2 03 02 01 0a a3 0e 30
0010   0c 30 0a a1 04 02 02 00 95 a2 02 04 00 a4 81 8d
0020   30 81 8a a0 07 03 05 00 00 00 00 10 a1 17 30 15
0030   a0 03 02 01 01 a1 0e 30 0c 1b 0a 63 68 65 6d 69
0040   6b 61 64 7a 65 a2 10 1b 0e 41 54 48 45 4e 41 2e
0050   4d 49 54 2e 45 44 55 a3 23 30 21 a0 03 02 01 02
0060   a1 1a 30 18 1b 06 6b 72 62 74 67 74 1b 0e 41 54
0070   48 45 4e 41 2e 4d 49 54 2e 45 44 55 a5 11 18 0f
0080   32 30 31 38 30 31 30 33 30 36 30 34 30 37 5a a7
0090   06 02 04 64 21 bb 89 a8 14 30 12 02 01 12 02 01
00a0   11 02 01 10 02 01 17 02 01 19 02 01 1a
`)

	var req AS_REQ
	rest, err := asn1.Unmarshal(msgBytes, &req)
	if err != nil {
		t.Errorf("Failed to parse: %v", err.Error())
	}
	if len(rest) != 0 {
		t.Errorf("Expected no trailing data, got %v bytes", len(rest))
	}
}

func TestKrbError(t *testing.T) {
	msgBytes := utils.ParseWiresharkHex(`
0000   30 81 b2 a0 03 02 01 05 a1 03 02 01 1e a2 11 18
0010   0f 32 30 32 33 30 33 32 37 31 35 35 31 33 37 5a
0020   a4 11 18 0f 32 30 31 38 30 31 30 32 30 36 30 34
0030   30 37 5a a5 05 02 03 04 88 a8 a6 03 02 01 06 a7
0040   10 1b 0e 41 54 48 45 4e 41 2e 4d 49 54 2e 45 44
0050   55 a8 17 30 15 a0 03 02 01 01 a1 0e 30 0c 1b 0a
0060   63 68 65 6d 69 6b 61 64 7a 65 a9 10 1b 0e 41 54
0070   48 45 4e 41 2e 4d 49 54 2e 45 44 55 aa 23 30 21
0080   a0 03 02 01 02 a1 1a 30 18 1b 06 6b 72 62 74 67
0090   74 1b 0e 41 54 48 45 4e 41 2e 4d 49 54 2e 45 44
00a0   55 ab 12 1b 10 43 4c 49 45 4e 54 5f 4e 4f 54 5f
00b0   46 4f 55 4e 44
`)
	expected := KRB_ERROR{
		Pvno: 5,
		Msg_type: 30,  // krb-error
		Ctime: utils.ParseWiresharkTime("2023-03-27 15:51:37"),
		Stime: utils.ParseWiresharkTime("2018-01-02 06:04:07"),
		Susec: 297128,
		Error_code: 6, // principal unknown
		Crealm: "ATHENA.MIT.EDU",
		Cname: PrincipalName{1, []KerberosString{"chemikadze"}},
		Realm: "ATHENA.MIT.EDU",
		Sname: PrincipalName{2, []KerberosString{"krbtgt", "ATHENA.MIT.EDU"}},
		E_text: "CLIENT_NOT_FOUND",
	}

	// verify it can be parsed
	var parsed KRB_ERROR
	rest, err := asn1.Unmarshal(msgBytes, &parsed)
	if err != nil {
		t.Errorf("Failed to parse: %v", err.Error())
	}
	if len(rest) != 0 {
		t.Errorf("Expected no trailing data, got %v bytes", len(rest))
	}
	if es, ps := fmt.Sprintf("%+v", expected), fmt.Sprintf("%+v", parsed); es != ps {
		t.Errorf("Repr mismatch:\n exp: %v\n got: %v", es, ps)
	}

	// verify that it can be generated and serialization is reversible
	generatedBytes, err := asn1.Marshal(expected)
	if err != nil {
		t.Fatal("Failed to marshall message")
	}
	_, err = asn1.Unmarshal(generatedBytes, &parsed)
	if err != nil {
		t.Fatal("Failed to unmarshall message")
	}
	if es, ps := fmt.Sprintf("%+v", expected), fmt.Sprintf("%+v", parsed); es != ps {
		t.Errorf("Repr mismatch:\n exp: %v\n got: %v", es, ps)
	}
}