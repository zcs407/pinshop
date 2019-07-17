package models

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func DbConfigInfo(path string) map[string]string {
	//make map type myMap for dbconfig
	MyMap := make(map[string]string)
	//open configfile return file and err info
	f, err := os.Open(path)
	//will be panic if has some err
	if err != nil {
		panic(err)
	}
	//close file after end
	defer f.Close()
	//new reader buf for file
	r := bufio.NewReader(f)
	for {
		//readeline file and return []byte singol slice to b
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//Tirm space for file line data
		s := strings.ReplaceAll(string(b), " ", "")
		//get index of =
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		//get map key and judge is null or not
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}

		//get map value and judge is null or not
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		MyMap[key] = value
	}
	return MyMap
}
