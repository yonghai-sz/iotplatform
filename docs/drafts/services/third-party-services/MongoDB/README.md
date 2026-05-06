

```sh
docker container stop my-single-mongo-container
docker container rm my-single-mongo-container
```
`docker container exec -it my-mongo-container mongosh --username root`    

`docker container exec -it my-mongo-container mongosh admin`  在 admin 数据库下    



  
`mongosh --host mongodb0.example.com --port 28015`  
`mongosh "mongodb://mongodb0.example.com:28015"`  
`mongosh "mongodb://mongodb0.example.com:28015" --username alice --authenticationDatabase admin`   

to connect to a database called `qa` on localhost:  `mongosh "mongodb://localhost:27017/qa"`   

To verify your current database connection, use the `db.getMongo()` method.    


  


查看当前数据库:        `db`   
查看当前的所有的数据库: `show dbs`   

查看数据库下所有集合:   
`use database`   
`show collections`   

退出容器命令： `exit`      









`df -h`    
`du -sh /var/lib/docker/volumes/my_mongodb-data`    
      

`db.locations.distinct( "deviceID",    { createTime: { $gt: 1702656000 }                                       }).length`
`db.locations.distinct( "deviceID",    { createTime: { $gt: 1702656000 },    upTime: { $lt: 1672635109000 }    })`
 

`db.locations.find({      deviceID: "99911500888823319999"     }).limit(1)`
`db.locations.find({      latitude:"NaN"                       }).sort({ field: -1 }).limit(1)`



`db.locations.countDocuments({       createTime: { $lt: 1725430920 },          level: { $exists: false }        })`
`db.locations.countDocuments({       createTime: { $lt: 1725430920 },          level: 0                         })`  
`db.locations.countDocuments({       createTime: { $gt: 1702656000 },          upTime: { $lt: 1672635109000 }   })` 



```
db.locations.deleteMany({        createTime: { $lt: 1725430920 },          level: 0     })
db.users.deleteOne({             _id: ObjectId("6629f22ac58c7bf9587177d4")  });
db.example.drop()
```