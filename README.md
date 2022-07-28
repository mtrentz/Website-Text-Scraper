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
There is a single endpoint at /api/scrape/ that takes in a POST with url, depth, and max_requests. If not provided, depth defaults to 2 and max_requests defaults to 200. For unlimited depth or max_requests, set them as -1.

Max requests is supposed to be used to limit the amount of scrapes, but it doesn't guarantee that no more than the exact amount of requests will be made.

# Example
```
curl -X POST http://localhost:8080/api/scrape/ -d '{"url":"https://www.bbc.com/", "max_requests":5}' -H "Content-Type: application/json"
```

Result summary:
```
{
    "url": "https://www.bbc.com/",
    "page_amount": 5,
    "visited_at": "2022-07-27 19:23:55",
    "pages": [
        {
            "url": "https://www.bbc.com/",
            "header": "Home\nNews\nSport\nWeather\niPlayer\nSounds\nBitesize\nCBeebies\nCBBC\nFood\nHome\nNews\nSport\nReel\nWorklife ...",
            "text": "BBC Homepage\nGas prices soar as Russia cuts German supply\nThe Nord Stream 1 pipeline is now operating at just ...",
            "footer": "\n\n\n\nHome\nNews\nSport\nWeather\niPlayer\nSounds\nBitesize\nCBeebies\nCBBC\nFood\nHome\nNews\nSport\nReel\n ...",
            "visited_at": "2022-07-27 19:23:56"
        },
        ...
    ]
}
```

Passsing a depth of 2 and unlimited requests
```
curl -X POST http://localhost:8080/api/scrape/ -d '{"url":"https://www.bbc.com/", "depth":2, "max_requests":-1}' -H "Content-Type: application/json"
```
