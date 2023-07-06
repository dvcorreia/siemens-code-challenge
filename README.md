<p align="center">Code challenge for the position of Cloud Native Developer at Siemens.</p>

- [Problem statement](#problem-statement)
- [Interpretation](#interpretation)
- [Constrains](#constrains)
- [Notes](#notes)
- [Getting started](#getting-started)
- [Debugging](#debugging)


## Problem statement

* it takes some time until a unicorn is produced, the request is blocked on requesters site and he needs to wait 

* to improve the situation adjust the code, so that the requester is receiving a request-id, with this request-id he can poll and validate if unicorns are produced

* if the unicorn is produced it should be returned though using fifo principle

* adjust the code, so that every x seconds a new unicorn is produced at put to a store, which can be used to fulfill the requests (LIFO Store)

* make sure, duplicate capabilities are not added to the unicorn

* improve the overall code

* if any requirements are not clear, compile meaningful assumptions

The provided reference implementation can be found at `cmd/reference`.

## Interpretation

> to improve the situation adjust the code, so that the requester is receiving a request-id, with this request-id he can poll and validate if unicorns are produced

We could approach this in many ways.
I decided to go with a custom http header - `X-Unicorn-Order-Id`, which can be changed programmatically.
I did not go with the `X-Request-Id` header because it's not official in the HTTP spec and it could be confused with a request ID that identifies each unique http request, and not the request transaction _per se_.

> if the unicorn is produced it should be returned though using fifo principle

> adjust the code, so that every x seconds a new unicorn is produced at put to a store, which can be used to fulfill the requests (LIFO Store)

I am a bit confused in what it is expected in these two points, so I'am not certain how it would fit in the reference code.
Are the LIFO store and the random unicorn FIFO production two separate things and add up, or one replaces the other?

**My interpretation is as follows:**

Server logic:
- Every x seconds, a factory produces a unicorn.
- If there are no orders pending, that unicorn will be stored in a LIFO store. Upon new order, these stored unicorns can fulfill the request.
- While orders are pending, any generated unicorn will be used to fulfill these orders in a FIFO principle (first ones to order, are the firsts to receive).

User API: 
- A user can order for an amount of unicorn to be produced.
- If there are no unicorns available for deliver, the available ones will and an Order ID will be delivered for subsequent poll of unicorns. 
- If there are enough unicorns to fullfil the order, the order will be fulfilled.
- Upon fulfilling an order, that order will be forgotten, so it will be considered invalid upon request with the same ID.


## Constrains

It will only be used the standard library available for go 1.19.
The project structure does not follow exactly how I would structured it for a production service.

The name and adjectives are loaded from text files.
It was kept this way, since it could be a requirement for someone other than a programmer to change or add its values.
If that's not the case, I would recommend what was done with the capabilities data, which follows closely what Docker's [moby](https://github.com/moby/moby/blob/master/pkg/namesgenerator/names-generator.go). also does.

## Notes

This was made after dinner until late at night, so the code is definitely not my best.
Some early decision were left in place, despite not being the best abstractions, due to time constrains.
The same goes for tests.

If anything more is to be done, please reach out to me at [dv_correia@hotmail.com](mainto:dv_correia@hotmail.com).

## Getting started

Compile the application:

```console
go build ./cmd/unicorn
```

Check its help:

```console
./unicorn

Usage of ./unicorn:
  -addr string
        http server address (default ":8000")
  -rate duration
        period in which the production line will generate a new unicorn (default 5s)
```

Run the application:

```console
./unicorn
```

## Debugging

If you run the application, you will get something like this:

```logs
unicorn: main.go:38: setting up service ...
unicorn: main.go:39: config: prod-rate=5s
unicorn: main.go:91: listening http at :8000
unicorn: logger.go:24: storage: stored unicorn<cheerful-josephina>, now with 1
unicorn: logger.go:24: storage: stored unicorn<hurtful-karoline>, now with 2
unicorn: logger.go:24: storage: stored unicorn<edible-celinda>, now with 3
unicorn: logger.go:36: storage: collected 1 from the requested 1
unicorn: logger.go:23: Request GET /unicorns?amount=1 ~ status=200 id= took=178.759µs size=130
unicorn: logger.go:24: storage: stored unicorn<frivolous-loretta>, now with 3
^Cunicorn: main.go:108: by by, from unicorn application
```

It shows 4 unicorns being generated and stored in the LIFO store (3 before and 1 after the request).

It also shows a HTTP request to `localhost:8000/unicorns?amount=1`, where you see 1 unicorn being collected from the store and being returned:

```console
curl "localhost:8000/unicorns?amount=1 | jq"
```
```json
{
   "pending":0,
   "orderId":"QALjNQXJGGjX11Bt",
   "unicorns":[
      {
         "name":"edible-celinda",
         "capabilities":[
            "change color",
            "swim",
            "design"
         ]
      }
   ]
}
```

You can make a big request and pool for more unicorns. Starting the server from scratch to be able to clearly see what's happening - here are the logs and commands used to request the server:

```logs
unicorn: main.go:38: setting up service ...
unicorn: main.go:39: config: prod-rate=5s
unicorn: main.go:91: listening http at :8000
unicorn: logger.go:24: storage: stored unicorn<unrealistic-elton>, now with 1
unicorn: logger.go:36: storage: collected 1 from the requested 20
unicorn: logger.go:23: Request GET /unicorns?amount=20 ~ status=200 id= took=824.152µs size=126
unicorn: logger.go:23: Request GET /unicorns ~ status=200 id=847umsuGRb8MiKO6 took=99.892µs size=355
^Cunicorn: main.go:108: by by, from unicorn application
```

We can see that we saved a unicorn in the store before requesting.
That unicorn was retrieved from the store and returned.
Them we were able to retrieve 4 more unicorns, so 15 were left to be produced.

> Note: no logs show unicorns being put in the store because they are fulfilling the orders.

```console
curl "localhost:8000/unicorns?amount=20"
{"pending":19,"orderId":"847umsuGRb8MiKO6","unicorns":[{"name":"unrealistic-elton","capabilities":["design","walk","talk"]}]}
```

```console
curl "localhost:8000/unicorns" --header "X-Unicorn-Order-Id: 847umsuGRb8MiKO6" | jq
```

```json
{
  "pending": 15,
  "orderId": "847umsuGRb8MiKO6",
  "unicorns": [
    {
      "name": "shabby-jeane",
      "capabilities": [
        "fullfill wishes",
        "fighting capabilities",
        "fly"
      ]
    },
    {
      "name": "brilliant-sharika",
      "capabilities": [
        "lazy",
        "fullfill wishes",
        "design"
      ]
    },
    {
      "name": "frivolous-piedad",
      "capabilities": [
        "walk",
        "cry",
        "fly"
      ]
    },
    {
      "name": "superficial-lucio",
      "capabilities": [
        "swim",
        "cry",
        "run"
      ]
    }
  ]
}
```