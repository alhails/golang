package struct__test

import (
	"github.com/golang/protobuf/jsonpb"
	_struct "github.com/searKing/golang/thirdparty/github.com/golang/protobuf/ptypes/struct"
	"strings"
	"testing"
)

type Human struct {
	Name      string
	Friends   []string
	Strangers []Human
}

type ToProtoStructTests struct {
	input  Human
	output string
}

var (
	toProtoStructTests = []ToProtoStructTests{{
		input: Human{
			Name:    "Alice",
			Friends: []string{"Bob", "Carol", "Dave"},
			Strangers: []Human{{
				Name:    "Eve",
				Friends: []string{"Oscar"},
				Strangers: []Human{{
					Name:    "Isaac",
					Friends: []string{"Justin", "Trent", "Pat", "Victor", "Walter"},
				}},
			}},
		},
		output: `{
  "Friends": [
     "Bob",
     "Carol",
     "Dave"
    ],
  "Name": "Alice",
  "Strangers": [
     {
        "Friends": [
           "Oscar"
          ],
        "Name": "Eve",
        "Strangers": [
           {
              "Friends": [
                 "Justin",
                 "Trent",
                 "Pat",
                 "Victor",
                 "Walter"
                ],
              "Name": "Isaac",
              "Strangers": null
             }
          ]
       }
    ]
 }`,
	},
	}
)

func TestToProtoStruct(t *testing.T) {
	for m, test := range toProtoStructTests {
		humanStructpb, err := _struct.ToProtoStruct(test.input)
		if err != nil {
			t.Errorf("#%d: ToProtoStruct(%+v): got: _, %v exp: _, nil", m, test.input, err)
		}

		marshal := jsonpb.Marshaler{EmitDefaults: false, Indent: " ", OrigName: true}
		humanStr, err := marshal.MarshalToString(humanStructpb)

		if err != nil {
			t.Errorf("#%d: json.Marshal(%+v): got: _, %v exp: _, nil", m, test.input, err)
		}

		if strings.Compare(humanStr, test.output) != 0 {
			t.Errorf("#%d: json.Marshal(%+v): got: %v exp: %v", m, test.input, humanStr, test.output)
		}
	}
}
