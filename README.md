# Hitotose

## Memo

``` sh
docker compose up -d

sudo docker exec -it 104 bash
sudo docker cp 202408 104:/202408
sudo docker exec -i 104 /usr/bin/mongorestore --db hitotose /202408/hitotose

sudo docker exec 104 mongodump --db hitotose --out /mongodump/20241125
sudo docker cp 104:/mongodump/20241125 ~/mongodump/20241125

# mongosh
db.game.updateMany({}, { $unset: { rating: "" } })
```
