
```

use admin

db.createUser( {
    user: "admin",
    pwd: "admin",
    roles: [ { role: "root", db: "admin" } ]
});

db.createUser( {
    user: "quickpay",
    pwd: "quickpay",
    roles: [ { role: "readWrite", db: "quickpay" } ]
});

```
