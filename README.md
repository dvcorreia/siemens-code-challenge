# unicorn

Code challenge for the position of Cloud Native Developer at Siemens.

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
I decided to go with a custom http header - `Unicorn-Request-Id`.
I did not go with the `X-Request-Id` header because it's not official in the HTTP spec and it could be confused with a request ID that identifies each unique http request, and not the request transaction _per se_.




## Assumptions

It will only be used the standard library available for go 1.19.
The project structure does not follow exactly how I would structured it for a production service.

---

The name and adjectives are loaded from text files.
It was keept this way, since it could be a requirement for someone other than a programmer to change or add its values.
If that's not the case, I would recommend what was done with the capabilities data, which follows closely what Docker also does:
[https://github.com/moby/moby/blob/master/pkg/namesgenerator/names-generator.go](https://github.com/moby/moby/blob/master/pkg/namesgenerator/names-generator.go).