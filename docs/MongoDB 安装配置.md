
```

apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
echo "deb http://repo.mongodb.org/apt/ubuntu "$(lsb_release -sc)"/mongodb-org/3.0 multiverse" \
 | sudo tee  /etc/apt/sources.list.d/mongodb-org-3.0.list
apt-get update
apt-get install -y mongodb-org

# rm /etc/mongod.conf
rm /etc/init/mongod.conf
rm /etc/init.d/mongod
rm -rf /var/log/mongodb

mkdir -p /opt/mongo/{rs11,rs12,rs13,arb1}
cd /opt/mongo

mongod --port=27017 --fork --dbpath=/opt/mongo/rs11 --logpath=/opt/mongo/rs11/mongod.log &
mongod --port=27018 --fork --dbpath=/opt/mongo/rs12 --logpath=/opt/mongo/rs12/mongod.log &
mongod --port=27019 --fork --dbpath=/opt/mongo/rs13 --logpath=/opt/mongo/rs13/mongod.log &
mongod --port=30000 --fork --dbpath=/opt/mongo/arb1 --logpath=/opt/mongo/arb1/mongod.log &

mongo --port=27017
mongo --port=27018
mongo --port=27019
mongo --port=30000

use admin
db.createUser( {
    user: "admin",
    pwd: "admin",
    roles: [ { role: "root", db: "admin" } ]
});

use quickpay
db.createUser( {
    user: "quickpay",
    pwd: "quickpay",
    roles: [ { role: "readWrite", db: "quickpay" } ]
});

exit


mongod --shutdown --dbpath=/opt/mongo/rs11
mongod --shutdown --dbpath=/opt/mongo/rs12
mongod --shutdown --dbpath=/opt/mongo/rs13
mongod --shutdown --dbpath=/opt/mongo/arb1


openssl rand -base64 741 > rs1.key
chmod 600 rs1.key

rsync -v root@10.171.199.158:/opt/mongo/rs1.key .

mongod --port=27017 --auth --fork --dbpath=/opt/mongo/rs11 --logpath=/opt/mongo/rs11/mongod.log \
 --replSet=rs1  --keyFile=/opt/mongo/rs1.key
mongod --port=27018 --auth --fork --dbpath=/opt/mongo/rs12 --logpath=/opt/mongo/rs12/mongod.log \
 --replSet=rs1  --keyFile=/opt/mongo/rs1.key
mongod --port=27019 --auth --fork --dbpath=/opt/mongo/rs13 --logpath=/opt/mongo/rs13/mongod.log \
 --replSet=rs1  --keyFile=/opt/mongo/rs1.key
mongod --port=30000 --auth --fork --dbpath=/opt/mongo/arb1 --logpath=/opt/mongo/arb1/mongod.log \
 --replSet=rs1  --keyFile=/opt/mongo/rs1.key


mongo --port=27017
use admin
db.auth('admin','admin')
rs.initiate()
#rs.add('mgo1.set.shou.money:27017')
rs.add('mgo1.set.shou.money:27018')
rs.add('mgo2.set.shou.money:27017')
rs.add('mgo2.set.shou.money:27018')
rs.addArb('mgo2.set.shou.money:30000')
rs.status()

db.bindingInfo.createIndex({ bindingId : 1, merId : 1 },{ unique: true });
db.bindingMap.createIndex({ bindingId : 1, merId : 1 },{ unique: true });
db.trans.createIndex({ orderNum : 1, merId : 1 },{ unique: true });
db.trans.createIndex({ transType : 1, refundOrderNum : 1, merId : 1, transStatus : 1 });
db.transSett.createIndex({ orderNum : 1, merId : 1 },{ unique: true });
db.merchant.createIndex({merId : 1 },{ unique: true });
db.cardBin.createIndex({bin : 1, cardLen : 1},{ unique: true });
db.cfcaBankMap.createIndex({insCode : 1 },{ unique: true });
db.chanMer.createIndex({ chanMerId : 1, chanCode : 1 },{ unique: true });
db.routerPolicy.createIndex({ merId : 1, cardBrand : 1 },{ unique: true });
db.respCode.createIndex({ respCode : 1},{ unique: true });
db.respCode.ensureIndex({"cfca.code":1});
db.respCode.ensureIndex({"cil.code":1});

```
