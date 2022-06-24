package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Arguments map[string]string

type myItem struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
type myData []myItem

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}()

	var fileContent []byte
	fileContent, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var data myData
	if len(fileContent) != 0 {
		err = json.Unmarshal(fileContent, &data)
		if err != nil {
			panic(err)
		}
	}
	defer func() {

	}()

	switch args["operation"] {
	case "add":
		if args["item"] == "" {
			return errors.New("-item flag has to be specified")
		}
		var newItem myItem
		err = json.Unmarshal([]byte(args["item"]), &newItem)
		if err != nil {
			panic(err)
		}
		for _, storedItem := range data {
			if storedItem.ID == newItem.ID {
				_, err = writer.Write([]byte(fmt.Sprintf("Item with id %v already exists", newItem.ID)))
				if err != nil {
					panic(err)
				}
				return nil
			}
		}

		data = append(data, newItem)

	case "list":
		_, err = writer.Write(fileContent)
		if err != nil {
			panic(err)
		}

	case "findById":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}

		var out []byte
		for _, storedItem := range data {
			if storedItem.ID == id {
				out, err = json.Marshal(storedItem)
				break
			}
		}

		_, err = writer.Write(out)
		if err != nil {
			panic(err)
		}

	case "remove":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		idExists := false
		for index, storedItem := range data {
			if storedItem.ID == id {
				data = append(data[:index], data[index+1:]...)
				idExists = true
				break
			}
		}
		if !idExists {
			_, err = writer.Write([]byte(fmt.Sprintf("Item with id %v not found", id)))
			if err != nil {
				panic(err)
			}
		}

	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}

	fileContent, err = json.Marshal(data)

	err = ioutil.WriteFile(args["fileName"], fileContent, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func parseArgs() Arguments {
	var idFlag = flag.String("id", "", "")
	var itemFlag = flag.String("item", "", "")
	var operationFlag = flag.String("operation", "", "")
	var fileNameFlag = flag.String("fileName", "", "")
	flag.Parse()

	var args Arguments = make(Arguments, 4)
	args["id"] = *idFlag
	args["item"] = *itemFlag
	args["operation"] = *operationFlag
	args["fileName"] = *fileNameFlag

	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
