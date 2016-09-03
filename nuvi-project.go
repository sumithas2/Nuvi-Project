/*
Script Name: nuvi-project.go

Purpose: Application downloads the contents from the URL: http://bitly.com/nuvi-plz. This URL is an http folder containing a list of zip files. This Script Creates a new work directory from the  current work directory, downloads all the zip files to the work directory and unzip all the zip files and extracts all the contents from each zip file to corresponding Zip Folder. Each zip folder contains a bunch of xml files. Each xml file contains 1 news report with xml tags as <document>,<type>,<forum>.. etc and the xml contents as key -value pair are stored
to Redis Database as Set.


Dependency: import these go packages :
	   archive/zip
	   golang.org/x/net/html
	   golang.org/x/net/html/atom
	   encoding/xml
	   github.com/astaxie/goredis


Author : Sumitha S

Date: 01 Sep 2016

*/

package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/goredis"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/* XML structure, add more xml tags Here */

type NewsXML struct {
	Type  string `xml:"type"`
	Forum string `xml:"forum"`
}

/* Global string constant directory name for all the work to be done */

const targetDir string = "work"

/* getContentfromURL: gets the content from the input URL, searches and gets href tags from the response and
calls the other function to download the contents from the href link ending with ".zip"
*/

func getContentfromURL(url string) {

	err := os.Mkdir("."+string(filepath.Separator)+targetDir, 0777)
	if err != nil {
		fmt.Println("Error while creating TargetDir", "-", err)
	}

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
	}

	//To get the redirected link

	finalURL := response.Request.URL.String()

	fmt.Printf("Given URL Redirected to: %v\n", finalURL)

	defer response.Body.Close()

	doc := html.NewTokenizer(response.Body)
	for tokenType := doc.Next(); tokenType != html.ErrorToken; {
		token := doc.Token()
		if tokenType == html.StartTagToken {
			if token.DataAtom != atom.A {
				tokenType = doc.Next()
				continue
			}
			// Extract url searching for href

			for _, attr := range token.Attr {
				if attr.Key == "href" {

					// Concatenates the redirected URL with href value

					hrefURL := string(finalURL + attr.Val)

					//Calls func downloadFile if the href link ends with .zip

					if strings.Contains(hrefURL, ".zip") {

						downloadFile(hrefURL)

					}

				}
			}

		}
		tokenType = doc.Next()

	}
	response.Body.Close()

}

/* downloadFile: downloads the links from the input:url with the name as last string
from input url  into the target work directory
*/

func downloadFile(url string) {

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	targetDirPath := string("./" + targetDir)

	extractedFilePath := filepath.Join(targetDirPath, fileName)

	output, err := os.OpenFile(extractedFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("Error in func downloadFile while creating", fileName, "-", err)
		return
	}

	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in func downloadFile while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error in func downloadFile while Copying the response.body to Output file", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")

}

/* processData: redirects to  target work directory and loops through all the zip files and
 uncompress all the zip files and extracts the contents to the same named Zip directory,
parses each xml file and stores to the database.
Calls two functions unzipFile ,parseXmlAndInsertToDB.

*/

func processData(targetDir string) {

	fmt.Println("Processing each downloaded zip file")

	zipFiles, err := ioutil.ReadDir(targetDir)
	if err != nil {
		fmt.Println("Error Reading Work dir:", err)
	}
	for _, eachZipFile := range zipFiles {

		zipFileName := string(eachZipFile.Name())

		unzipDirName := strings.Trim(zipFileName, ".zip")
		fmt.Println("Extracting zip file", zipFileName, "to", unzipDirName)

		unzipFile(zipFileName, unzipDirName)

		fmt.Println("Extracting Completed for ", unzipDirName)

		fmt.Println("Parsing each XML File inside ", unzipDirName)

		parseXmlAndInsertToDB(unzipDirName)

	}

}

/*Uncompress the zip files and extracts the contents to the same named zip folder */

func unzipFile(archive, target string) error {

	workDirPath := targetDir + "/"

	reader, err := zip.OpenReader(workDirPath + archive)
	if err != nil {
		fmt.Println("Error in func unzipFile Reading Zip File:", err)

	}

	defer reader.Close()

	if err := os.MkdirAll(workDirPath+target, 0755); err != nil {
		fmt.Println("Error in func unzipFile creating zip extract dir:", err)

	}

	for _, file := range reader.File {
		path := filepath.Join(workDirPath+target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}

	}

	return nil
}

/* parseXmlAndInsertToDB: to parse xml files and inserts the contents to Database,
calls another function:insertToRedisDB */

func parseXmlAndInsertToDB(unzipDirName string) error {

	xmlDirPath := targetDir + "/" + unzipDirName

	xmlFiles, err := ioutil.ReadDir(xmlDirPath)
	if err != nil {
		fmt.Println("Error in func parseXmlAndInsertToDB Reading Work dir:", err)
	}
	for _, eachXmlFile := range xmlFiles {

		if err != nil {
			fmt.Println("Error in func parseXmlAndInsertToDB Reading each XML File:", err)
		}

		eachXmlFileName := string(xmlDirPath + "/" + eachXmlFile.Name())
		parseXmlFile, err := os.Open(eachXmlFileName)
		if err != nil {
			return err
		}
		defer parseXmlFile.Close()

		var newsdata NewsXML

		b, _ := ioutil.ReadAll(parseXmlFile)

		xml.Unmarshal(b, &newsdata)

		newsContent := fmt.Sprintf("DocumentName:%s type: %s, forum: %s", unzipDirName+"/"+eachXmlFile.Name(), newsdata.Type, newsdata.Forum)

		//Connecting and Storing the Parsed XMl Contents To DB

		//fmt.Println("Inserting the News Contents to Redis DB")

		insertToRedisDB(newsContent)

	}

	return nil

}

/* insertToRedisDB: Use the Redis Sets to store the news xml to Redis Database. */

func insertToRedisDB(newsxml string) {

	var client goredis.Client

	client.Addr = "127.0.0.1:6379" // Set the default port in Redis

	newsXmlList := []string{newsxml}
	for _, news := range newsXmlList {
		client.Sadd("NEWS-XML", []byte(news)) // Inserting to RedisDB using Redis Sets Sadd
	}

	/*Commented: To Print the contents from the Redis Sets

	  NewsXMlContents,_ := client.Smembers("NEWS-XML")
	  for i, v := range NewsXMlContents {
	      println(i,":",string(v))
	  }
	*/

}

func clearWorkDir(targetDir string) {
	fmt.Println("Clearing the temp work Dir")
	os.RemoveAll(targetDir)

}

/*Main function starts Here */

func main() {

	mainurl := "http://bitly.com/nuvi-plz"
	getContentfromURL(mainurl)
	processData(targetDir)
	clearWorkDir(targetDir)
	fmt.Println("Exit")

}
