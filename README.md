![citizens-logo](https://raw.githubusercontent.com/greivinlopez/citizens/master/citizenslogosmall.png)

A [Go](http://golang.org/) solution to provide a fast way to consult Costa Rican citizens basic information

## API Documentation

The Citizens API uses the [REST](http://en.wikipedia.org/wiki/Representational_state_transfer) architectural style.  The API is organized into resources and each resource should be reachable using an unique [URI](http://en.wikipedia.org/wiki/Uniform_resource_identifier).  Different actions occurs depending on the HTTP verb you use to access the resource.  Documentation related to the API can be found here: [Citizens REST API documentation](https://github.com/greivinlopez/citizens/blob/master/CitizensAPIDocumentation.pdf?raw=true)

## Ubuntu Server Setup

### Install MongoDB

Install MongoDB following this instructions: [Installing MongoDB on Ubuntu](https://www.digitalocean.com/community/tutorials/how-to-install-mongodb-on-ubuntu-12-04)

Run the mongo shell:

```console
mongo
```

Add admin user to mongo:

```javascript
use admin

db.createUser(
  {
    user: "<adminuser>",
    pwd: "<adminpassword>",
    roles:
    [
      {
        role: "userAdminAnyDatabase",
        db: "admin"
      }
    ]
  }
)
```
Substitute the values between the <> with the ones you want to use for MongoDB admin user

Edit mongodb.conf and add auth=true

http://stackoverflow.com/questions/6235808/how-can-i-restart-mongodb-with-auth-option-in-ubuntu-10-04

Restart mongodb service

```console
sudo service mongodb restart
```

Run the MongoDB shell

```console
mongo
```

Create a user for people database

```javascript
use admin
db.auth('<adminuser>', '<adminpassword>')
use people
db.createUser(
	{
	    user: "<dbuser>",
	    pwd: "<dbpassword>",
	    roles : [
			{
				role : "readWrite",
				db : "people"
			},
			{
				role : "dbAdmin",
				db : "people"
			}
		]
	}
)
db.auth('<dbuser>', '<dbpassword>')
```

Again remember to substitute the values between <> with your own credentials.

Create a dummy document to save the new database permanently

```javascript
db.users.save( {username:"glopez"} )
```

There are other security measures you can follow regarding the MongoDB server:

https://www.digitalocean.com/community/tutorials/how-to-securely-configure-a-production-mongodb-server

### Installing Go

```console
cd ~
wget https://storage.googleapis.com/golang/go1.3.1.linux-amd64.tar.gz
tar -zxvf go1.3.1.linux-amd64.tar.gz
rm go1.3.1.linux-amd64.tar.gz
```

Binary distributions assume they will be installed at /usr/local/go Otherwise, you must set the GOROOT environment variable

```console
sudo mv ~/go /usr/local
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile
```

### Install Dependencies

```console
apt-get install git
apt-get install bzr
```

### Server Configuration

The server application follows the recomendation of the twelve-factor app regarding [configuration](http://12factor.net/config). Store the configuration in environment variables. To set the variables use the "export" command. For instance:

```console
sudo nano ~/.bash_profile
```

Add the following lines to the end of the file

```console
export GOPATH=/root/packages/
export CZ_DB_ADDRESS="localhost"
export CZ_DB_USER="<dbuser>"
export CZ_DB_PASS="<dbpassword>"
export CZ_API_KEY="<apikey>"
```

Save the file and return to the command line.  Run the new configuration

```console
source ~/.bash_profile
```

The values you use for *dbuser* and *dbpassword* must be the same you use when creating the database user in the previous steps.  The value for *apikey* will hold the API key the clients will use to make requests to the server. Ensure to use a long key value that includes numbers, letters, symbols and to keep it secured.

## Running Server

Install package dependencies

```console
go get -u gopkg.in/mgo.v2
go get -u gopkg.in/martini.v1
go get -u gopkg.in/greivinlopez/skue.v2
```

Download source code

```console
mkdir citizens
cd citizens
git init
git remote add --track master origin https://github.com/greivinlopez/citizens.git
git pull
```

### Fill database

Download the complete "Padron Electoral":

```console
cd filldb
curl -O http://www.tse.go.cr/zip/padron/padron_completo.zip
```

Install zip and unzip facilities to the server:

```console
sudo apt-get install zip unzip
```

Unzip the archived file:

```console
unzip padron_completo.zip
```

Compile and fill database

```console
go build filldb.go
./filldb
```

The filldb command will extract the data from the "padron electoral" files and create the documents on the database. The process will run for 5-10 minutes.

### Compile and run the Server

```console
cd ..
cd server
go build server.go
./server
```
[![License](http://img.shields.io/:license-mit-blue.svg)](http://opensource.org/licenses/MIT)

[![baby-gopher](https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)
