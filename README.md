# IITK-coin

My submission for IIT coin
Index.go is the starting point of the app, all response and request follow JSON format. /login and /signup handles only POST requests, whereas /secretpage handles only GET requests

# Example Query syntax for different endpoints  -

## /signup

Handler - signupHandler.go <br>
Request type - POST <br>
JSON - {"name":"Abhinav", "rollno":"190031","password":"mypassword"}

## /login

Handler - loginHandler.go<br>
Request type - POST <br>
JSON - {"rollno":"190031","password":"mypassword"}

## /secretpage

Handler -secretpageHandler.go<br>
Request type - GET <br>

## /addcoins

Handler - addCoinsHandler.go<br>
Request type - POST <br>
JSON - {"rollno":"190031","coins":"mypassword"}

## /transfercoin

Handler - transferCoinHandler.go<br>
Request type - POST <br>
JSON - {"firstrollno":"190031","secondrollno":"190031","amount":10}

## /getrcoin

Handler - getCoinHandler.go<br>
Request type - GET <br>
eg url query - http://localhost:8080/getcoins?rollno=190031
