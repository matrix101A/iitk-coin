# IITK-coin

My submission for IIT coin
Index.go is the starting point of the app, all response and request follow JSON format. /login and /signup handles only POST requests, whereas /secretpage handles only GET requests

# Example Query syntax for different endpoints -

## /signup

Handler - signupHandler.go <br>
Request type - POST <br>
Requires Authorization - NO<br>
JSON Format - {"name":"Abhinav", "rollno":"190031","password":"mypassword",
"account_type":"member"}

#### Account type decides the type of user - It can be "member" , "CTM"(For core team ) or "admin"(For gen sec )

note - enter account_type value correctly(without typo or whitespaces)

## /login

Handler - loginHandler.go<br>
Request type - POST <br>
Requires Authorization - NO<br>
JSON Format - {"rollno":"190031","password":"mypassword"}

## /secretpage

Handler -secretpageHandler.go<br>
Requires Authorization - YES<br>
Request type - GET <br>

## /addcoins

Handler - addCoinsHandler.go<br>
Request type - POST <br>
Requires Authorization - YES<br>
JSON Format - {"rollno":"190031","coins":"200","remarks":"Why the coins are bein rewarded "}

### Only CTM and admin can use this endpoint to add coins to other members and CTM(only admins can add to CTM )

## /transfercoin

Handler - transferCoinHandler.go<br>
Request type - POST <br>
Requires Authorization - YES<br>
JSON Format- {"rollno":"190031","amount":10}

### Transfers coins from account of logged in user to user(rollno) and deducts appropiate tax . Amount need to be of double type.

## /getcoins

Handler - getCoinsHandler.go<br>
Requires Authorization - YES
Request type - GET <br>

### Gives number of coins for current user

## /redeem

Handler - redeeemCoinsHandler.go<br>
Request type - POST <br>
Requires Authorization - YES<br>
JSON Format - {"itemid":1}

### Redeems the item with given item id and deducts the equivalet coins. itemid needs to be an integer. Can only redeem items after participating in certain events (set in env variables)

## /additems

Handler - addItemsHandler.go<br>
Request type - POST <br>
Requires Authorization - YES<br>
JSON Format - {"itemid":1,"cost":"1000","number":1}

### Adds new (or increases the quantity of ) items for redeem. itemid and number(number of items) need to be integers. Can only be accessed by CTM or admins

## All information related to maximum coins per account , min events to redeem items and secret key are set in the .env file
