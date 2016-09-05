package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strings"
	"time"
)

const (
	directoryName             string = "xmlFiles"
	xmlNumerOfElementsPerFile int    = 5
)

type subDocument struct {
	Name    string `xml:"name"`
	Surname string `xml:"surname"`
	Age     int    `xml:"age"`
}

type document struct {
	XMLName      xml.Name      `xml:"document"`
	SubDocuments []subDocument `xml:"subDocument"`
}

func main() {
	fileName, err := getFirstXMLFileFound()
	if err != nil {
		writeToLogFile(err.Error())
		return
	}
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		writeToLogFile(err.Error())
		return
	}
	var xmlData document
	err = xml.Unmarshal(b, &xmlData)
	if err != nil {
		writeToLogFile(err.Error())
		return
	}
	err = createSplittedXMLFiles(&xmlData)
	if err != nil {
		writeToLogFile(err.Error())
		return
	}
	writeToLogFile("Ndarja e XML u krye me sukses")
}

func writeToLogFile(message string) {
	content := fmt.Sprintf("%s: %s\n", time.Now(), message)

	f, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}

func getFirstXMLFileFound() (string, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".xml") || strings.Contains(file.Name(), ".XML") {
			return file.Name(), nil
		}
	}
	return "", errors.New("Nuk u gjet asnje file xml ne direktorine ku ekzekutohet programi")
}

func createSplittedXMLFiles(doc *document) error {
	exists := func(name string) bool {
		if _, err := os.Stat(name); err != nil {
			if os.IsNotExist(err) {
				return false
			}
		}
		return true
	}
	if exists(directoryName) {
		err := os.RemoveAll(directoryName)
		if err != nil {
			return err
		}
	}
	err := os.Mkdir(directoryName, 0777)
	if err != nil {
		return err
	}
	div := float64(len(doc.SubDocuments)) / float64(xmlNumerOfElementsPerFile)
	nrFiles := int(math.Ceil(div))
	lastFileOffset := len(doc.SubDocuments) % xmlNumerOfElementsPerFile
	currOffset := 0
	for i := 0; i < nrFiles; i++ {
		if i == nrFiles-1 {
			currOffset = lastFileOffset
		}
		fileContent := document{
			XMLName:      doc.XMLName,
			SubDocuments: doc.SubDocuments[i*xmlNumerOfElementsPerFile : (i+1)*xmlNumerOfElementsPerFile-currOffset],
		}
		cont, err := xml.MarshalIndent(fileContent, "", "   ")
		if err != nil {
			return err
		}
		splitFileName := path.Join(directoryName, fmt.Sprintf("splitFile%d.xml", i))
		contentStr := xml.Header + string(cont)
		ioutil.WriteFile(splitFileName, []byte(contentStr), 0644)
	}
	return nil
}
