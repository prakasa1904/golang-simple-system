## Description

Golang Web Application Boilerplate, modular clean code architecture seperate code by service to isolate unique functionality the service. And simplify the changes in the service without crossing the border each folder structure. The idea is to prepare if for each service require to move to stand alone service.

Use command below to run:

0. Before use `make run` or `make run-hot`, execute command `make init` to prepare development environment

```sh
make init
```

1. Open other terminal execute command to run css `make css-dev`

```sh
make css-dev
```

1. Use `make run` or `make run-hot` to run service with Makefile or use go binary command `export ENV=development && go run cmd/webapp/*.go`. 

```sh
make run-hot

# or

make run

# or

export ENV=development && go run cmd/webapp/*.go
```

2. Open [http://localhost:3000](http://localhost:3000/)

## Todo
- [ ] Integrate with tailwind CLI
- [ ] Deploy to [dPanel](https://cloud.terpusat.com/)


## Reference

- [Clean Code Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)