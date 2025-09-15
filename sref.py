from habanero import Crossref
from pprint import pprint
import sqlite3


class Reference():
    def __init__(self, message):
        self.ref_type = message['type']
        self.doi = message['DOI']
        self.title = message['title']
        self.authors = message['author']
        self.journal = message['container-title']
        self.short_journal = message['short-container-title']  # If empty, use the same than the journal 
        self.pages = message['page']
        self.issue = message['issue']
        self.url = message['URL']


con = sqlite3.connect(":memory:")
cur = con.cursor()

# TODO: Format authors names
# TODO: Check if paper or book


def fetch_paper_info(query):
    # Fetch the doi
    cr = Crossref()
    res = cr.works(query=query, limit=1)

    paper_doi = res["message"]["items"][0]['DOI']
    
    # Here fetch all information for that DOI
    res = cr.works(ids = paper_doi)
    return Reference(res['message'])


  
REF = "references.txt"
with open(REF, "r") as f:
    content = f.readlines()


ref = fetch_paper_info(content[3])

pprint(ref.title)


