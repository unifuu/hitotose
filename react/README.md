# Hitotose-React

``` sh
# Create Project
npx create-react-app hitotose --template typescript

# Install MUI
npm install @mui/material @emotion/react @emotion/styled

# Install Materials Icons
npm install @mui/icons-material

# Install MUI Lab
npm install @mui/lab

# Install React Bootstrap Icons
npm i react-bootstrap-icons

# Install React Router
npm install react-router-dom

# Install JS-COOKIE
npm i --save-dev @types/js-cookie

npm install -g serve

serve -s build -l 5000
```

## Temp

``` mongosh
db.game.updateMany({}, { $rename: { "developer_id": "developer" } });
db.game.updateMany({}, { $rename: { "publisher_id": "publisher" } });
db.game.updateMany({}, { $rename: { "how_long_to_beat": "time_to_beat" } });

# db.game.updateMany({}, { $unset: { developer: "" } });
# db.game.updateMany({}, { $unset: { publisher: "" } });

db.game.find({}).forEach(function(doc) {
  if (doc.developer instanceof ObjectId) {
    db.game.updateOne(
      { _id: doc._id },
      { $set: { developer: doc.developer.toString() } }
    );
  }
});

db.game.find({}).forEach(function(doc) {
  if (doc.publisher instanceof ObjectId) {
    db.game.updateOne(
      { _id: doc._id },
      { $set: { publisher: doc.publisher.toString() } }
    );
  }
});
```
