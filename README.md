# constant-api-service
constant-api-service

### Config database
    - Rename: `conf.json.example` to `conf.json`
    - Configs: 
        ```
            "env": "localhost",
            "port": 8089,
            "db": "root:root@tcp(localhost:3306)/constant_mvp?charset=utf8&parseTime=True&loc=Local",
        ```

### Install go packages
    run command:  ```
        brew install dep
        dep ensure
        go run server.go 
    ```


