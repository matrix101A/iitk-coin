# IITK-coin
My submission for IIT coin
Index.go is the  starting point of the app, all response and request follow JSON format. /login and /signup handles only POST requests, whereas /secretpage handles only GET requests 

# Example request -
## /signup 
POST
{"name":"Abhinav", "rollno":"190031","password":"mypassword"}

## /login 
POST
{"rollno":"190031","password":"mypassword"}

## /secretpage
GET
