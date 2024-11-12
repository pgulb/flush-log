# Flush Log  
Flush Log is Progressive Web App that can be used to track your bowel movements  
it is available at https://pgulb.github.io/flush-log/  
  
## Why  
You can improve your toilet habits, for example reduce time mindlessly scrolling through
Social Media if you see how much time it takes  
FL allows you to track stats like time, number/% of times with phone, mean time, mean rating etc
  
## What is Progressive Web App  
In short it's a web app that allows itself to be 'installed' both on desktop and mobile devices
as an app, option to install appears usually at the right side of searchbar, and on mobile
in the three dot menu
  
PWAs aim to look more like native application than a regular website, they also cache themselves on
user's device to load faster and notify user when an app can be updated  
  
### Initial diagram for how it should work (now slightly modified)  
<img src="./systems.png" alt="app diagram" align="center"/>
  
## development info  
  
flush-log uses Taskfile (https://taskfile.dev/) to automate commands used to
develop and test its components  
Components are:
- PWA (Go + go-app https://go-app.dev/)  
- REST API (Python + FastAPI https://fastapi.tiangolo.com/)  
- MongoDB (At production it's MongoDB Atlas free tier https://www.mongodb.com/docs/atlas/reference/free-shared-limitations/)  
  
### App server vs static website
Go-app can be compiled to executable to run as a server, or compiled to static content  
In production I compile it to static and deploy to Github Pages  
Tests and local development is done with server version (pipeline tests with Docker)  
  
---
### List all available tasks  
```sh
task
```
---
### Prepare venv
I use uv https://github.com/astral-sh/uv to create venv and install dependencies while developing  
Tests are run using venv  
```sh
task init-uv-venv
```
---
### Spin up whole stack in Docker locally and watch logs  
  
```sh
task dev
```
And go to http://localhost:8080 for hot-reloading app  
API is located at http://localhost:6789  
and Mongo at mongodb://localhost:27017  
  
---
### Remove containers  
```sh
task cleanup
```
---
### API tests
unit  
```sh
task test-api-unit
```
mocked  
```sh
task test-api-mock
```
with local mongo  
```sh
task test-api-integration
```
---
### PWA tests
unit  
```sh
task test-pwa-unit
```
integration (with go-rod https://go-rod.github.io/#/)  
```sh
task test-pwa-integration
```
integration with visible browser window  
```sh
task test-pwa-integration-show-window
```
## Github Workflows
- test-api.yml - run tests for API on each push
- test-pwa.yml - build docker image and run tests for PWA on each push
- publish-api-image.yml - build and publish API docker image after tests complete (on main branch)
- deploy-api.yml - deploy API image after build (currently to VPS through ssh) (on main branch)
- deploy-pwa.yml - build static content and deploy to Github Pages after tests complete (on main branch)
