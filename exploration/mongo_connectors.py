import urllib


from pymongo import MongoClient


class MongoConnector:

    def __init__(self, host, user='', password=''):
        self.mongourl = 'mongodb://'
        if user != '':
            esc_user = urllib.parse.quote(user)
            self.mongourl += esc_user
            if password != '':
                esc_pass = urllib.parse.quote(password)
                self.mongourl += f':{esc_pass}'
            self.mongourl += '@'
        self.mongourl += host
        self.client = None

    def _create_client(self):
        self.client = MongoClient(self.mongourl)

    def get_client(self):
        if self.client is None:
            self._create_client()
        return self.client


class MongoDatabase:

    def __init__(self, client, database_name: str):
        self.name = database_name
        self.database = client[database_name]

    def get_database(self):
        return self.database


class MongoCollection:

    def __init__(self, database, collection: str):
        self.name = collection
        self.collection = database[collection]

    def insert(self, doc):
        if '_id' in doc:
            raise ValueError('cannot insert document with custom _id')
        self.collection.insert_one(doc)

    def get_by_name_desc(self, name, desc=''):
        self.collection.find_one({
            'name': name,
            'description': desc
        })

    def update_name(self, oid, new_name):
        pass
