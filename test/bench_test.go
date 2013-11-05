package test

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"testing"
)

var codeJSON []byte
var codeStruct codeResponse

func codeInit() {
	f, err := os.Open("testdata/code.json.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	codeJSON = data

	/*
		if err := Unmarshal(codeJSON, &codeStruct); err != nil {
			panic("unmarshal code.json: " + err.Error())
		}

		if data, err = Marshal(&codeStruct); err != nil {
			panic("marshal code.json: " + err.Error())
		}

		if !bytes.Equal(data, codeJSON) {
			println("different lengths", len(data), len(codeJSON))
			for i := 0; i < len(data) && i < len(codeJSON); i++ {
				if data[i] != codeJSON[i] {
					println("re-marshal: changed at byte", i)
					println("orig: ", string(codeJSON[i-10:i+10]))
					println("new: ", string(data[i-10:i+10]))
					break
				}
			}
			panic("re-marshal code.json: different result")
		}
	*/
}

func BenchmarkCodeEncoder(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	enc := NewcodeResponseJSONEncoder(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		if err := enc.Encode(&codeStruct); err != nil {
			b.Fatal("Encode:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}
