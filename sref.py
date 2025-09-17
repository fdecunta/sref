import argparse
import os
import sqlite3

from habanero import Crossref

# TODO: Format authors names
# TODO: Check if paper or book


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


def init_db(content):
    schema = """
    CREATE TABLE IF NOT EXISTS refs (
        input_ref TEXT NOT NULL UNIQUE,
        crossref_output TEXT,
        fmt_ref TEXT
    )
    """

    con = sqlite3.connect(":memory:")
    cur = con.cursor()
    cur.execute(schema)

    cur.executemany(
        "INSERT INTO refs(input_ref) VALUES(?)",
        [(r.strip(),) for r in content]
    )
    con.commit()   

    return con


def fetch_paper_info(query):
    print("Looking for reference at Crossref...")
    # Fetch the doi
    cr = Crossref()
    res = cr.works(query=query, limit=1)

    paper_doi = res["message"]["items"][0]['DOI']
    
    # Here fetch all information for that DOI
    res = cr.works(ids = paper_doi)
    return Reference(res['message'])


def read_file(file):
    if not os.path.isfile(file):
        print(f"Error: not a file. {file}")
        exit(1)
    
    with open(file, "r") as f:
        return f.readlines()



if __name__ == "__main__":
    parser = argparse.ArgumentParser()

    parser.add_argument("file")

    args = parser.parse_args()

    content = read_file(args.file)
    print(content) 

#con = init_db(content)
#print(con.execute("SELECT * FROM refs LIMIT 1").fetchone()[0])
#con.close()

#ref = fetch_paper_info(content[0])
#print(ref.authors)



