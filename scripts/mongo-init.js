var admin = db.getSiblingDB('admin');
admin.auth('admin-user', 'admin-password');

var db = db.getSiblingDB('combotest');

db.createUser({
        user: 'test-user',
        pwd: 'test-password',
        roles : [
            {
                role: 'readWrite',
                db: 'combotest',
            },
        ],
    });


db.createCollection("users");

db.createCollection("events");

db.users.insertOne({ 
        "_id" : ObjectId("600893550b1d7baabe1e01a4"),
        "role" : "admin", 
        "login" : "admin", 
        "password" : "9bc7aa55f08fdad935c3f8362d3f48bcf70eb280", 
        "confirmed" : true 
    });