from bson.objectid import ObjectId

def convert(doc):
    if isinstance(doc, dict):
        for key, value in doc.items():
            if isinstance(value, ObjectId):
                doc[key] = str(value)
            elif isinstance(value, dict):
                convert(value)
            elif isinstance(value, list):
                for item in value:
                    convert(item)
    elif isinstance(doc, list):
        for item in doc:
            convert(item)
    return doc