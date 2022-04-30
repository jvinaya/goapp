# mini-aspire-API

loan application API
-How you can run the app:
I have dockerized the application and added docker compose file.
command to run (after installing docker into your machine)
"docker-compose up"
-for api documentation I have added postman collection in project

-I have added Makefile so that other developer can go with ease and use to run
-unit test
-Migrate db
-db related action (creating ,droping db)
-sqlc
-others

---

Packages I have used

- viper :
  for loading configuration

- Gin :
  Web FrameWork

- Sqlc :
  to generate db through sql schema and query

- Paseto And Jwt :
  For Authentication (I have implemented both you can use any of them. Paseto set as Default)

- pq:
  postgres support
