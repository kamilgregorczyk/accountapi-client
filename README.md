![Build](https://github.com/kamilgregorczyk/accountapi-client/workflows/Build/badge.svg) ![Docker-Compose](https://github.com/kamilgregorczyk/accountapi-client/workflows/Docker-Compose/badge.svg)

Author: Kamil Gregorczyk

Dear reviewers, thanks for looking into my home task. I'm new to Go (I come from java & python world) but I tried to not make a lot of rookie mistakes.

What I did:

* Implemented two clients. The http oriented one exists in case there's a need to reuse that http client for other domain oriented clients, the second one is Account orriented.

* For resilience I implemented retries with exponential backoff.

* I'm also requiring a timeout to be defined from the client as I noticed that builtin go client doesn't have any

* Every call in `account.Client` has the ability to pass context. I initially though about populating trace/span ids for distributed tracing but then I realised that Go doesn't have any generic interface for that and some monitoring tools could be not compatible with that one. So it's up to the caller to populate context with its own tracing tools or do different deadlines for calls.

* I added healthcheck & wait to accountapi in docker-compose to avoid restarts of accountapi in case postgres is not yet ready

* I tried to keep number of testing libs low therefore I (somewhat) did BDD style tests with go's built in logs, it's not the nicest way of doing it but it resembles BDD to some degree.

* I added basic validation of requests, like checking if ID is actually UUID V1 to avoid pointless requests


What I would have done more/differently:

* If infra doesn't support it, I would have also add circuit breaker, I didn't add it as it's more complicated than doing simple retries and on a prod ready solution I would have used a 3rd party lib.

* Use some simple BDD framework for tests

* Since the documentation of form3 domains is very accurate (at least with accountapi) I'd probably generate such clients
