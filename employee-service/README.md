# MrAzharuddin/employee-crud/employee-service
employee-service


### REST Server
- To run the server, you can use the following command:
    ```go
    go mod tidy
    go run main.go
    ```
- To find the documentation, go to browser and paste this url:
    [http://localhost:8000/swagger/index.html](http://localhost:8000/swagger/index.html)

- Test is incomplete, but you can run:
    ```go
    go test -v ./...
    ```

- Used sqlite with gorm, so no need to setup any database.

[![Open in DevPod!](https://devpod.sh/assets/open-in-devpod.svg)](https://devpod.sh/open#https://github.com/MrAzharuddin/employee-crud/employee-service)