citizens
========

Go solution to provide a fast way to consult Costa Rican citizens basic information

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

db.addUser( { user: "adminuser",
              pwd: "adminpassword",
              roles: [ "userAdminAnyDatabase" ] } )
```

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

Create a user for necesitotaxi

```javascript
use admin
db.auth('adminuser', 'adminpassword')
use people
db.addUser( { user: "dbuser",
              pwd: "dbpass",
              roles: [ "readWrite", "dbAdmin" ]
            } )
db.auth('dbuser', 'dbpass')
```

Create a dummy document to save the new database permanently

```javascript
db.users.save( {username:"glopez"} )
```

There are other security measures you can follow regarding the MongoDB server:

https://www.digitalocean.com/community/tutorials/how-to-securely-configure-a-production-mongodb-server

### Installing Go

```console
cd ~
wget http://golang.org/dl/go1.2.2.linux-amd64.tar.gz
tar -zxvf go1.2.2.linux-amd64.tar.gz
rm go1.2.2.linux-amd64.tar.gz
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

```console
sudo nano ~/.bash_profile
```

Add the following lines to the end of the file

```console
export GOPATH=/root/packages/
export CZ_DB_ADDRESS="localhost"
export CZ_DB_USER="dbuser"
export CZ_DB_PASS="dbpass"
```

Save the file and return to the command line.  Run the new configuration

```console
source ~/.bash_profile
```

## Running Server

Install package dependencies

```console
go get labix.org/v2/mgo
go get github.com/go-martini/martini
go get github.com/greivinlopez/skue
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

