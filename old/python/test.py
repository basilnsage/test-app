import bson
import json
import sys


from mongo_connectors import MongoCollection, MongoDatabase, MongoConnector
from mongo_migrations import wishlist_items_create_index


def main(wishlist_items):
    mongo_conn = MongoConnector('localhost').get_client()
    wishlist_db = MongoDatabase(mongo_conn, 'wishlistApp').get_database()
    wishlist_coll = MongoCollection(wishlist_db, 'wishlistItems')

    wishlist_items_create_index(wishlist_coll.get_collection())

    for item in wishlist_items:
        wishlist_coll.insert(item)

    for item in wishlist_items:
        name, desc = item['name'], item['description']
        oid = wishlist_coll.get_by_name(name)['_id']
        oid = bson.objectid.ObjectId(oid)
        wishlist_coll.update_name(oid, name + ' UPDATED')


if __name__ == '__main__':
    item_file = sys.argv[1]
    with open(item_file) as f:
        items = json.load(f)['items']
    main(items)
