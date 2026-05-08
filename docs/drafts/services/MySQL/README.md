

```
docker container exec -it my-mysql-container mysql -u root -p
EXIT

SELECT user, host, plugin FROM mysql.user;
ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'new_password';
```

```sh
docker container stop my-single-mysql-container
docker container rm my-single-mysql-container
```
