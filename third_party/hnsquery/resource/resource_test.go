// Adapted from https://github.com/mslipper/handshake

package resource

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestResource_EncodeDecode(t *testing.T) {
	tests := []struct {
		name   string
		infile string
	}{
		{
			"proofofconcept name update in block 8578",
			"proofofconcept_4b131b575145a6d0b44654241e89c02a9e316e1acac0ad53edd4d8bb7af3ce8f.bin",
		},
		{
			"ix name update in block 8293",
			"ix_d74fad18bda1e83d2405aac3ea260513bec2cc71e39ce87d24b76fbe6a911c9c.bin",
		},
		{
			"lifelong name update in block 10162",
			"lifelong_1d0f8de2757488cbd59bea7b8f7c7ad5aa9ebd6459631e801a041062338a8630.bin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expData, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", tt.infile))
			if err != nil {
				t.Fatalf("got error: %v", err)
			}

			resource := new(Resource)
			if err := resource.Decode(bytes.NewReader(expData)); err != nil {
				t.Fatal(err)
			}

			actData := new(bytes.Buffer)
			if err := resource.Encode(actData); err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expData, actData.Bytes()) {
				t.Fatal("not equal")
			}
		})
	}
}
