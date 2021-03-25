## Go-based clean architecture

This is my personal perspective towards clean architecture application in Golang.
Feedback are most welcome. 

## Purpose

While I do not have any intention to make a framework from this and as it is a very opinionated structure, this would be useful for me to as starting point of a project, 
but if you find this structure useful to you and you have some improvement in mind, just let me know by opening issues.

Also, I'm curious of what kind of complexity I need to face if I want to make a monolithic code for microservices.

## What's inside
- Domains: contains repository and use case interfaces as well as base entity.
- App: Bare code and implementation of domains. It provides basic implementation and validation of the entity. It also connects to the defined persistence storage. It **doesn't** contain any front-facing layer (http api, grpc, etc)
- Api: Exposing the app through Http. The responsibility is basically transforming request to app and response from app if needed. 
- Worker: Exposing the app through background worker. Similar to Api, but only works in background app fashion. 

## TODO
- [ ] Add more tests for repository and handler
- [ ] Expose some entity via microservice
- [ ] Put it in docker (docker-compose as well)