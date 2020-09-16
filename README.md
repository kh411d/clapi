# clapi

Self-hosted serverless/faas blogging clap api 

### Key-Value Database

Currently, I only support redis at [lambda.store](https://lambda.store/), for free account you'll get Max 5000 Commands Daily and 256 MB Max data size per DB

### Environment Variables

Set these vars on any serverless provider you choose,

- `REDIS_HOST` (i.e. us1-xxxxx-xxxx-32223.lambda.store:32223)
- `REDIS_PASSWORD` (i.e. j20fj0293jf0293jf02090209fj02j0239)

Netlify,

Edit `GO_IMPORT_PATH` value corresponds to your Github account

### Deploying on Vercel/Zeit 

Within the `/api` directory of your projects, Vercel will automatically recognize the languages listed on this page, through their file extensions, and serve them as Serverless Function. Go files in the api directory that export a function matching the net/http Go API will be served as Serverless Functions.

```
//ServeHTTP net/http handler
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	...
}
```

It seems that when you put more than one go file in the `api` folder the system will examine all the file to look for net/HTTP exported functions, if there is no exported function in one of that files, it will resulting error thus the deployment will fail, so it is why I put all the code in on file handler.go

After deployed, you'll have some domains to be used on function,
```
$ curl -X POST -d '{"url":"http://clapi/clap"}' https://clapi.vercel.app/api/handler
{"statusCode":200,"headers":null,"multiValueHeaders":null,"body":"ok"}

$ curl -X GET https://clapi.vercel.app/api/handler?url=http://clapi/clap
{"statusCode":200,"headers":null,"multiValueHeaders":null,"body":"4"}
```
