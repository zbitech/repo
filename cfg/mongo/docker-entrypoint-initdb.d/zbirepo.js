//echo "Creating mongo users..."
//mongo admin --host localhost -u admin -p password --eval "db.createUser({user: 'zbiadmin', pwd: 'password', roles: [{role: 'readWrite', db: 'zbiRepo'}]});"
//echo "Mongo users created."

//db.auth('admin','password');
db = db.getSiblingDB('zbiRepo');
db.createUser({user: 'zbiadmin', pwd: 'password', roles: [{role: 'readWrite', db: 'zbiRepo'}]})

let res = [
    db.createCollection("users"),
    db.users.createIndex({ "userid": 1, "email": 1 }, { unique: true }),

    db.createCollection("passwords"),
    db.passwords.createIndex({"userid": 1}, {unique: true}),

    db.createCollection("user_policy"),
    db.user_policy.createIndex({ "userid": 1 }, { unique: true }),

    db.createCollection("apikeys"),
    db.apikeys.createIndex({ "key": 1, "userid": 1 }, { unique: true }),

    db.createCollection("apikey_policy"),
    db.apikey_policy.createIndex({ "key": 1}, { unique: true }),

    db.createCollection("teams"),
    db.teams.createIndex({ "teamid": 1, "owner": 1 }, { unique: true }),

    db.createCollection("team_members"),
    db.team_members.createIndex({ "teamid": 1, "email": 1, "key": 1}, { unique: true }),

    db.createCollection("projects"),
    db.projects.createIndex({ "name": 1, "owner": 1, "team": 1}, { unique: true }),

    db.createCollection("instances"),
    db.instances.createIndex({ "project": 1, "name": 1, "type": 1, "owner": 1}, { unique: true }),

    db.createCollection("k8s_resources"),
    db.k8s_resources.createIndex({ "project": 1, "instance": 1, "type": 1, "name": 1}, { unique: true }),

    db.createCollection("instance_policy"),
    db.instance_policy.createIndex({ "project": 1, "instance": 1 }, { unique: true }),
]

printjson(res)
