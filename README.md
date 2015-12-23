spoofer.go
==========

#### what/why?

mitigating CORS would be the objective of this project. control over certain
headers being sent will expand the capabilities of how data is pushed around
on the internet.

#### endpoints

```/spoof``` would perhaps be the only endpoint ever necessary, but that might
be subject to change.

#### howto

access the ```/spoof``` endpoint with url parameters:

```url```: URL to be gotten

```referer```: in-case you want to replace ```referer``` header
