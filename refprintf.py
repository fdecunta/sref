#!/usr/bin/python3

import json
import os
import sys

#   """
#   %A    authors. The last one separated with '&'.
#   %A10  only first ten authors
#   %Y    year
#   %T    title
#   %J    journal name
#   %j    abbreviated journal name
#   %V    volumen
#   %I    issue
#   %P    pages
#   %D    DOI
#   %U    URL
#   
#   For example:
#   %A.(%Y).%T.%V(%I)%P.%D 
#   """

class Reference():
    def __init__(self, message):
        self.ref_type = message['type']
        self.doi = message['DOI']
        self.title = message['title']
        self.authors = message['author']
        self.year = message['created']['date-parts'][0][0]
        self.journal = message['container-title']
        self.short_journal = message['short-container-title']  # If empty, use the same than the journal 
        self.pages = message['page']
        self.issue = message['issue']
        self.url = message['URL']


def abbreviate_names(given_names):
    short_name = ""
    for name in given_names.split():
        name = name.replace(".", "")
        short_name = short_name + name[0] + ". "

    short_name = short_name.strip()    # Remove trailing space
    return short_name
    

def format_authors(authors):
    authors_str = ""

    if len(authors) == 1:
        family_name = authors[0]['family']
        initials = abbreviate_names(authors[0]['given'])
        return f"{family_name}, {initials}."
        

    is_last = False
    for author in authors:
        if author == authors[-1]:
            is_last = True

        family_name = author['family']
        initials = abbreviate_names(author['given'])

        if is_last:
            # If is the last, remove the previous comma and add an ampersand
            # TODO: This will break for one author
            authors_str = authors_str[:-2] + " & " + family_name + " " + initials
        else:
            authors_str = authors_str + family_name + " " + initials + ", "

    return authors_str
    

def format_ref(ref):
    authors_str = format_authors(ref.authors)

    ref_str = f"{authors_str}({ref.year}). {ref.title})"
    print(ref_str)


if __name__ == "__main__":
    file = sys.argv[1]
    
    if not os.path.isfile(file):
        print("File is not a file")
        sys.exit(1)
    
    with open(file, "r") as f:
        ref_json = json.load(f)
    
    ref = Reference(ref_json)
    format_ref(ref)


