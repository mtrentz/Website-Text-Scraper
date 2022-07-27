# Website-Text-Scraper
Simple API with a single endpoint to scrape a website and other pages in the same domain. The response separates the text from the header, footer and body of each page.

# Running
Clone this repo and run:

```
go get
go run .
```

This will start the API at port 8080.

# Running with docker
Get and run the image from Dockerhub:
```
docker run -p 8080:8080 mtrentz/website_text_scraper:latest
```

# Usage
There is a single endpoint at /api/scrape/ that takes in a POST with url, depth, and max_requests. If not provided the default parameters is a depth of 2.

Max requests is supposed to be used to limit the amount of scrapes, but it doesn't guarantee that no more than the provided calls will be made.

# Example
```
curl -X POST http://localhost:8080/api/scrape -H "Content-Type: application/json" -d '{"url":"https://hltv.org", "max_requests":5}'
```

Result summary:
```
{
    "url": "https://www.hltv.org/",
    "page_amount": 5,
    "visited_at": "2022-07-27 18:07:59",
    "pages": [
        {
            "url": "https://www.hltv.org/",
            "header": "\nExpand\n\n\n\n\n\n\n\nAll\n(17)\nCasters\n(9)\nStreamers\n(7)\nOrganizers\n(1)\n\n\n\n\n",
            "text": "News\nMatches\nResults\nEvents\nStats\nGalleries\nRankings\nForums\nBetting\nLive\nFantasy\nForgot password ...",
            "footer": "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n",
            "visited_at": "2022-07-27 18:08:00"
        },
        ...
    ]
}
```
