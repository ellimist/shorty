Shorty Challenge
================

The trendy modern question for developer inteviews seems to be, "how to create an url shortner". Not wanting to fall too far from the cool kids, we have a challenge for you!

## The Challenge

The challenge, if you chose to accept it, is to create a micro service to shorten urls, in the style that TinyURL and bit.ly made popular.

## Rules

1. The service must expose HTTP endpoints according to the definition below.
2. The service must be self contained, you can use any language and technology you like, but it must be possible to set it up from a fresh install of Ubuntu Server 14.04, by following the steps you write in the README.
3. It must be well tested, it must also be possible to run the entire test suit with a single command from the directory of your repository.
4. The service must be versioned using git and submitted by making a Pull Request against this repository, git history **should** be meaningful.
5. You don't have to use a datastore, you can have all data in memory, but we'd be more impressed if you do use one.

**Good Luck!** — not that you need any ;)

## API Documentation

**All responses must be encoded in JSON and have the appropriate Content-Type header**


### POST /shorten

```
POST /shorten
Content-Type: "application/json"

{
  "url": "http://example.com",
  "shortcode": "example"
}
```

Attribute | Description
--------- | -----------
**url**   | url to shorten
shortcode | preferential shortcode

##### Returns:

```
201 Created
Content-Type: "application/json"

{
  "shortcode": :shortcode
}
```

A random shortcode is generated if none is requested, the generated short code has exactly 6 alpahnumeric characters and passes the following regexp: ```^[0-9a-zA-Z_]{6}$```.

##### Errors:

Error | Description
----- | ------------
400   | ```url``` is not present
409   | The the desired shortcode is already in use. **Shortcodes are case-sensitive**.
422   | The shortcode fails to meet the regex following regexp: ```^[0-9a-zA-Z_]{4,}$```.


### GET /:shortcode

```
GET /:shortcode
Content-Type: "application/json"
```

Attribute      | Description
-------------- | -----------
**shortcode**  | url encoded shortcode

##### Returns

**302** response with the location header pointing to the shortened URL

```
HTTP/1.1 302 Found
Location: http://www.example.com
```

##### Errors

Error | Description
----- | ------------
404   | The ```shortcode``` cannot be found in the system

### GET /:shortcode/stats

```
GET /:code
Content-Type: "application/json"
```

Attribute      | Description
-------------- | -----------
**shortcode**  | url encoded shortcode

##### Returns

```
200 OK
Content-Type: "application/json"

{
  "startDate": "2012-04-23T18:25:43.511Z",
  "lastSeenDate": "2012-04-23T18:25:43.511Z",
  "redirectCount": 0
}
```

Attribute         | Description
--------------    | -----------
**startDate**     | date when the url was encoded, conformant to [ISO8601](http://en.wikipedia.org/wiki/ISO_8601)
**redirectCount** | number of times the endpoint ```GET /shortcode``` was called
lastSeenDate      | date of the last time the a redirect was issued, not present if ```redirectCount == 0``

##### Errors

Error | Description
----- | ------------
404   | The ```shortcode``` cannot be found in the system


