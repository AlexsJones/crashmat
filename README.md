crashmat
======== 

`Status: In development`

A simple web API for call stack logging over the web then query results based on application ID filter

`go get github.com/AlexsJones/crashmat`

API
==

```
GET /Upload
GET /Upload/{applicationId}
GET /Application
```

```
POST /Upload 
{
"authorization":"Basic",
"username":"Bob",
"password":"Password",
"applicationid":"01",
"raw":"{}"
}

POST /Application 
{
"authorization":"Basic",
"username":"Bob",
"password":"Password",
"applicationid":"01",
}
```
```
DELETE /application/{applicationid}
{
"authorization":"Basic",
"username":"Bob",
"password":"Password",
}
```
Configuration
============

Either populate values in the `conf/app.json` and/or add secret information toenvironmental variables e.g.

```
export CRASHMAT_CLIENTSECRET=""
export CRASHMAT_CLIENTID=""
export CRASHMAT_ELASTICHOSTADDRESS=""
export PORT="8080"
export CRASHMAT_POSTGRESCONNECTION="host=localhost port=5432 user=bob password=bob"
export CRASHMAT_UPDATEFREQ=9000
```
