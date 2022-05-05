package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/golang/protobuf/proto"
	"movieexample.com/gen"
	"movieexample.com/metadata/pkg/model"
)

func main() {
	metadata := &model.Metadata{
		ID:          "123",
		Title:       "The Movie 2",
		Description: "Sequel of the legendary The Movie",
		Director:    "Foo Bars",
	}

	jsonBytes, err := json.Marshal(metadata)
	if err != nil {
		panic(err)
	}

	xmlBytes, err := xml.Marshal(metadata)
	if err != nil {
		panic(err)
	}

	genMetadata := &gen.Metadata{
		Id:          "123",
		Title:       "The Movie 2",
		Description: "Sequel of the legendary The Movie",
		Director:    "Foo Bars",
	}

	protoBytes, err := proto.Marshal(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size:\t%dB\n", len(jsonBytes))
	fmt.Printf("XML size:\t%dB\n", len(xmlBytes))
	fmt.Printf("Proto size:\t%dB\n", len(protoBytes))
}
