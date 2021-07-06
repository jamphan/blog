---
title: "Web-scraping Real-Estate data with Python"
date: 2020-10-05T11:51:18+11:00
draft: false
toc: true
tags: [python, webscrapping]
---

I was looking to collect some real-estate data in Australia to do some financial modelling.

I decided to quickly write up a scrapping tool that takes data that is publicly available on some real estate websites.

This is pretty easy with Python's [`Beautiful Soup`](https://pypi.org/project/beautifulsoup4/).

The basic idea:

1. Fetch the webpage.
2. Read-in the webpage into an object we can manipulate (the soup).
3. Parse the soup for the information we need.

## Fetching

I'm using Python's [``requests``](https://requests.readthedocs.io/en/master/) library.

The code:

``` py
class RequestManager(object):

    def __init__(self, host, outdir='data', max_cache_time=300):

        self._host = re.sub('/', '', host)
        self._outdir = os.path.join(outdir, self._host + os.sep)
        self._max_cache_time = max_cache_time

        print(self._outdir)
        if not(os.path.isdir(self._outdir)):
            os.mkdir(self._outdir)

    def endpoint_to_path(self, endpoint: str) -> str:
        """ Will make URL path ascii safe
        """

        return os.path.join(self._outdir, *[re.sub('\W', '-', x)  for x in endpoint.split('/')])

    def touch(self, endpoint: str) -> typing.BinaryIO:
        """ Creates a tree in the outdir that reflects the URL path.

        A .meta file can be used to determine the latest file added (for
        caching)
        """

        path_safe = self.endpoint_to_path(endpoint)
        if not(os.path.isdir(path_safe)):
            os.makedirs(path_safe)

        unix_time_now = int(time.time())
        file_meta = os.path.join(path_safe, '.meta')
        file_out = os.path.join(path_safe, str(unix_time_now))

        with open(file_meta, 'w', encoding='utf-8') as fdmeta:
            fdmeta.write(str(unix_time_now))

        return open(file_out, 'wb', encoding='utf-8')

    def get_cache(self, endpoint: str):
        _path = self.endpoint_to_path(endpoint)
        with open(os.path.join(_path, '.meta'), 'r', encoding='utf-8') as fdmeta:
            last_opened = int(fdmeta.read())

        unix_time_now = int(time.time())
        if (unix_time_now - last_opened) > self._max_cache_time:
            return None
        else:
            return get_cache_time(endpoint, last_opened)

    def get_cache_time(self, endpoint: str, time: int):
        _path = self.endpoint_to_path(endpoint)
        timestamps = [int(f) for f in listdir(_path) if isfile(os.path.join(_path, f)) and f.isnumeric()]
        timestamps.sort(reverse=True)
        for ts in timestamps:

            # Return first file that where the given time is greater
            if (time >= ts):
                with open(os.path.join(_path, ts), 'rb', encoding='utf-8') as fd:
                    resp = pickle.load(fd)

                return resp

    def get(self, endpoint: str):

        if (cache := self.get_cache(endpoint)) is not None:
            return cache

        resp = ''#get request here
        with self.touch(endpoint) as fd:
            pickle.dump(resp, fd)

        return resp
```

#   # Request Manager

## Parsing