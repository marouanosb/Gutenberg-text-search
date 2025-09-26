import requests
import json

url = "https://gutendex.com/books/"
url_text = "https://gutenberg.org/cache/epub/{}/pg{}.txt"
response = requests.get(url)
books = response.json()['results']

def fetch_text(id):

    r = requests.get(url_text.format(id, id))
    r.encoding = 'utf-8-sig'
    text = r.text.replace('\r\n', '\n').strip()

    return text

# keep id and title
books = {
        book['id']:
            {
            "title": book['title'],
            "authors": [x['name'] for x in book['authors']],
            "text": fetch_text(book['id'])
            }
        for book in books
        }

# save json
with open('./resources/books.json', 'w', encoding='utf-8') as f:
    json.dump(books, f, indent=0)