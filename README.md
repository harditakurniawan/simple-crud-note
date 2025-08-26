# Simple CRUD Note
*This is a Rest API based monolithic application that provides basic CRUD with authentication to create some **Notes***

## Tech Stack
- Framework : Go Fiber
- Database : PostgreSQL & Redis
- Tools : Postman

## How To Install
I will use ***make*** command to install this application so make sure you have already installed ***GNU Make*** on your os.
1. make sure you already made "logs" folder in the root dir (for logging file)
2. create 2 files named "private.key" & "public.key" and save your private & public key in there
3. create a .env file from .env.example and fill in all the variables base on your local environment
4. ```make install```
5. to run dev-mode and watch for any changes : ```make run-watch``` | to stop dev-mode : ```control + Z``` or ```CTRL + Z```
6. to run prod-mod : ```make build && make run-prod``` | to stop prod-mode : ```make stop```