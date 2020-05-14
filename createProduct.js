
username = process.env("MONGO_USERNAME");
pwd = process.env("MONGO_USER_PASSWORD");
uri = `mongodb+srv://${username}:${pwd}@cluster0-yyofy.gcp.mongodb.net/test?retryWrites=true&w=majority`;
conn = new Mongo(uri);
db = conn.getDB("demo");
collection = db.stock;

var docToInsert = {
    name: "pineapple",
    quantity: 10
};

function sleepFor(sleepDuration) {
    var now = new Date().getTime();
    while (new Date().getTime() < now + sleepDuration) {
        /* do nothing */
    }
}

function create() {
    sleepFor(1000);
    print("inserting doc...");
    docToInsert.quantity = 10 + Math.floor(Math.random() * 10);
    res = collection.insert(docToInsert);
    print(res)
}

while (true) {
    create();
}