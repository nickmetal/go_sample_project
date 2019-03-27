### Build and run


```sh
$ docker build -t apiservice_transport_solver_image .
$ docker run -it --rm --name apiservice_transport_solver -p 0.0.0.0:3000:3000 apiservice_transport_solver_image
```