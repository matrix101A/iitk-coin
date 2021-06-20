# IITK-coin

My submission for IIT coin
Index.go is the starting point of the app, all response and request follow JSON format. /login and /signup handles only POST requests, whereas /secretpage handles only GET requests

# Example request -

## /signup

Handler - signupHandler.go
POST
{"name":"Abhinav", "rollno":"190031","password":"mypassword"}

## /login

Handler - loginHandler.go
POST
{"rollno":"190031","password":"mypassword"}

## /secretpage

Handler -secretpageHandler.go
GET

## /addcoins

Handler - addCoinsHandler.go
POST -
{"rollno":"190031","coins":"mypassword"}

## /transfercoin

Handler - transferCoinHandler.go
POST -
{"firstrollno":"190031","secondrollno":"190031","amount":10}

## /getrcoin

Handler - getCoinHandler.go
GET -
url query - http://localhost:8080/getcoins?rollno=190031
