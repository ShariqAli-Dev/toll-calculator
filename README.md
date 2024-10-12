# microservice

obu generates lat,long
reciever recieves lat,long - puts in a kafka que

### how to start

- `make reciever` (needs to be started first b/c obu connects to its ws)
- `make obu`
- `docker compose -up` to start kafka
