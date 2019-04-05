### Build and run

```sh
$ docker build -t apiservice_transport_solver_image .
$ docker run -it --rm --name apiservice_transport_solver -p 0.0.0.0:3000:3000 apiservice_transport_solver_image
```

#### Result:
1. microservice +
TODO - return json errors
2. Logging +
TODO - added logs for response
3. Unit test - 
4. Web Ui +-
5. Deploy -
TODO - add service to DO

#### Transportation issue results:
+ build basic result
+ check basic result and issue type(open/closed)
+ build cell differences
- build solve loop and find total result
