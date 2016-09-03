# Nuvi-Project
Nuvi Interview Code Project 


Application name:  nuvi-project.go

Language and Database : Golang and Redis

Purpose: Application downloads the contents from the URL: http://bitly.com/nuvi-plz. This URL is an http folder containing a list of zip files. This Script Creates a new work directory from the  current work directory, downloads all the zip files to the work directory and unzip all the zip files and extracts all the contents from each zip file to corresponding Zip Folder. Each zip folder contains a bunch of xml files. Each xml file contains 1 news report with xml tags as <document>,<type>,<forum>.. etc and the xml contents as key -value pair are stored
to Redis Database as Set.


Dependency: import these go packages :
	   archive/zip
	   golang.org/x/net/html
	   golang.org/x/net/html/atom
	   encoding/xml
	   github.com/astaxie/goredis



To Compile:

go install <path where the nuvi-project.go is placed>

for eg,
go install github.com/sumithas2/nuvi

To Run:

For eg,
go run nuvi-project.go 


Sample Console Output:

home@home:~/work/src/github.com/sumithas2/nuvi$ go run nuvi-project.go
Given URL Redirected to: http://feed.omgili.com/5Rh5AMTrc4Pv/mainstream/posts/
Downloading http://feed.omgili.com/5Rh5AMTrc4Pv/mainstream/posts/1472662644037.zip to 1472662644037.zip
10344980 bytes downloaded.
Downloading http://feed.omgili.com/5Rh5AMTrc4Pv/mainstream/posts/1472662818681.zip to 1472662818681.zip
10269573 bytes downloaded.

Processing each downloaded zip file
Extracting zip file 1472662644037.zip to 1472662644037
Extracting Completed for  1472662644037
Parsing each XML File inside  1472662644037
Extracting zip file 1472662818681.zip to 1472662818681
Extracting Completed for  1472662818681
Parsing each XML File inside  1472662818681

Improvements:

-Need to add log  
-Use Redis Hash type
-We can add new functionality to show the news-xml contents in web page. 
