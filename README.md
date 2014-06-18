citizens
========

Go solution to provide a fast way to consult Costa Rican citizens basic information

## Ubuntu Server Setup

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


## Installing Go

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





