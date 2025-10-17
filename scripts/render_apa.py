#!/usr/bin/python3

import json
from datetime import datetime
import sys

def format_apa_reference(ref_data):
    """Convert a reference from JSON format to APA style"""
    
    # Format authors
    authors = []
    for author in ref_data.get('author', []):
        if author.get('sequence') == 'first':
            # First author: Family, Initials
            given = author.get('given', '').replace('.', '').split()
            initials = '. '.join([name[0] for name in given if name]) + '.'
            if len(given) > 0:
                authors.append(f"{author['family']}, {initials}")
            else:
                authors.append(author['family'])
        else:
            # Additional authors: Initials Family
            given = author.get('given', '').replace('.', '').split()
            initials = '. '.join([name[0] for name in given if name]) + '.'
            if len(given) > 0:
                authors.append(f"{initials} {author['family']}")
            else:
                authors.append(author['family'])
    
    if len(authors) == 1:
        author_string = authors[0]
    elif len(authors) == 2:
        author_string = f"{authors[0]} & {authors[1]}"
    else:
        author_string = ', '.join(authors[:-1]) + ', & ' + authors[-1]
    
    # Format date
    date_parts = ref_data.get('issued', {}).get('date-parts', [[None]])[0]
    if date_parts and date_parts[0]:
        year = date_parts[0]
    else:
        year = "n.d."
    
    # Format title
    title = ref_data.get('title', [''])[0]
    
    # Handle different types
    ref_type = ref_data.get('type', '')
    
    if ref_type == 'journal-article':
        # Journal article format
        journal = ref_data.get('container-title', [''])[0]
        volume = ref_data.get('volume', '')
        issue = ref_data.get('issue', '')
        pages = ref_data.get('page', '')
        
        # Short journal name if available
        short_journal = ref_data.get('short-container-title', [''])[0]
        journal_display = short_journal if short_journal else journal
        
        volume_issue = f"{volume}"
        if issue:
            volume_issue += f"({issue})"
        
        return f"{author_string} ({year}). {title}. {journal_display}, {volume_issue}, {pages}. https://doi.org/{ref_data['doi']}"
    
    elif ref_type == 'book':
        # Book format
        publisher = "Springer"  # Default publisher for these references
        return f"{author_string} ({year}). {title}. {publisher}. https://doi.org/{ref_data['doi']}"
    
    elif ref_type == 'book-chapter':
        # Book chapter format
        book_title = ref_data.get('container-title', [''])[0]
        pages = ref_data.get('page', '')
        publisher = "Springer"  # Default publisher
        
        if len(ref_data.get('container-title', [])) > 1:
            book_title = ref_data['container-title'][-1]  # Use the last container title as book title
        
        return f"{author_string} ({year}). {title}. In {book_title} (pp. {pages}). {publisher}. https://doi.org/{ref_data['doi']}"
    
    else:
        # Default format for other types
        return f"{author_string} ({year}). {title}. https://doi.org/{ref_data['doi']}"

# Load and process your JSON data
def render_apa_references(json_file_path):
    with open(json_file_path, 'r', encoding='utf-8') as file:
        data = json.load(file)
    
    apa_references = []
    
    for doi, ref_data in data.items():
        try:
            apa_ref = format_apa_reference(ref_data)
            apa_references.append(apa_ref)
        except Exception as e:
            print(f"Error processing {doi}: {e}")
    
    # Sort references alphabetically by first author's last name
    apa_references.sort()
    
    return apa_references

# Example usage
if __name__ == "__main__":
    references = render_apa_references(sys.argv[1])
    
    for ref in references:
        print(f"{ref}")
