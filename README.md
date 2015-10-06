crashmat
======== 

A simple web API for bug logging into ElasticSearch

`go get github.com/AlexsJones/crashmat`

API
==

```
GET /Upload
GET /Upload/{ApplicationId}
GET /Auth/Github
```

```
POST /Upload 
{
"Applicationid":"01",
"raw":"{}"
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
```
