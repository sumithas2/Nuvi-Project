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


Author : Sumitha S

Date: 01 Sep 2016

