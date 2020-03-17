import pymongo


def wishlist_items_create_index(collection):
    collection.create_index([('name', pymongo.ASCENDING)], unique=True)
