citizens
========

Go solution to provide a fast way to consult Costa Rican citizens basic information

## Ubuntu Server Setup

Install MongoDB following this instructions: [Installing MongoDB on Ubuntu](https://www.digitalocean.com/community/tutorials/how-to-install-mongodb-on-ubuntu-12-04)

Run the mongo shell:

`mongo`

Add admin user to mongo:

`use admin`

`db.addUser( { user: "adminuser",
              pwd: "adminpassword",
              roles: [ "userAdminAnyDatabase" ] } )`

Edit mongodb.conf and add auth=true
(http://stackoverflow.com/questions/6235808/how-can-i-restart-mongodb-with-auth-option-in-ubuntu-10-04)

Restart mongodb service

`sudo service mongodb restart`

Run the MongoDB shell

`mongo



