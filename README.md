# clapi

Self-hosted serverless/faas blogging clap api 

### Key-Value Database

Currently, I only support redis at [lambda.store](https://lambda.store/), for free account you'll get Max 5000 Commands Daily and 256 MB Max data size per DB

### Environment Variables

Set these vars on any serverless provider you choose,

- `REDIS_HOST` (i.e. us1-xxxxx-xxxx-32223.lambda.store:32223)
- `REDIS_PASSWORD` (i.e. j20fj0293jf0293jf02090209fj02j0239)

Netlify only, edit `GO_IMPORT_PATH` value on netlify.toml file corresponds to your Github account

### Localhost

Build locally, don't forget to set your REDIS env vars,
```
$ go build -i
$ ./clapi

$ curl -X POST http://0.0.0.0:3000/clap?url=http://clapi/clap  
{"statusCode":200,"headers":null,"multiValueHeaders":null,"body":"ok"}  
$ curl -X GET http://0.0.0.0:3000/clap?url=http://clapi/clap
{"statusCode":200,"headers":null,"multiValueHeaders":null,"body":"7"}
```

### Deploying on Vercel/Zeit (Git Integration)

Go files in the `/api` directory that export a function matching the net/http Go API will be served as Serverless Functions.

```
//ServeHTTP net/http handler
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	...
}
```

Need to put all the code in one file `/api/handler.go`, otherwise deployment will fail

After deployed, as example you may access your function as this,  
url: [YOUR-DOMAIN]/api/handler

```
$ curl -X POST -d '{"url":"http://clapi/clap"}' https://clapi.vercel.app/api/handler
{"statusCode":200,"headers":null,"multiValueHeaders":null,"body":"ok"}

$ curl -X GET https://clapi.vercel.app/api/handler?url=http://clapi/clap
{"statusCode":200,"headers":null,"multiValueHeaders":null,"body":"4"}
```

### Deploying on Netlify (Git Integration)

Netlify can build your source Go functions into AWS Lambda compatible binaries.

Before building your Go source, Netlify needs to know the expected Go import path for your project. Use the `GO_IMPORT_PATH` environment variable to set the right import path. You can do this in your _netlify.toml_ file. The path value should point to your source repository on your Git provider, for example github.com/kh411d/clapi.

After successful deployment, you need to add env variables required for Redis in  
`Settings > Build & deploy > Environment > Environment variables`  
and then re-deploy the app from,  
`Deploys > Trigger deploy > Deploy site`

After deployed, as example you may access your function as this,  
URL: [YOUR-DOMAIN]/.netlify/functions/clapi

```
$ curl -X GET https://flamboyant-khorana-5b394b.netlify.app/.netlify/functions/clapi?url=http://clapi/clap
$ 6

$ curl -X  POST -d '{"url":"http://clapi/clap"}' https://flamboyant-khorana-5b394b.netlify.app/.netlify/functions/clapi
$ ok
```

### Deploying on Fly.io (flyctl)

I haven't tried this myself but it seems very possible, you need to have this `flyctl` tool installed on your system, upon installing, you have to register for fly.io account, simply follow the guidelines [Installing flyctl](https://fly.io/docs/getting-started/installing-flyctl/), and then continue with [Build, Deploy, and Run a Go Application](https://fly.io/docs/getting-started/golang/)

`flyctl init` will create a `fly.toml` file, Select Go (Go Builtin) for the builder, and then get ready for deployment `flyctl deploy`, keep in mind that flyctl will create docker images on your local computer, it will download all required images such as golang image approx 1gb, so make sure you have good storage space and internet connection.

`flyctl open` will open a browser on the HTTP version of the site.


