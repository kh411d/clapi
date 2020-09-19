# clapi

Self-hosted serverless/faas blogging clap api 

### Frontend

You can use this api for your own clap button or maybe you could give it a try for my ready made app [https://github.com/kh411d/clap-it](https://github.com/kh411d/clap-it)

### Database

Pick one of these databases provided, if you set both FaunaDB and Redis env vars, then FaunaDB will be the most likely to be chosen.

###### Fauna.com (FaunaDB)

- Go to [Fauna.com](https://fauna.com/) web console
- Create a collection, name it as `claps`, it will be used as key-value documents
- Create index for `claps` collection, 
  - Name it as `url_idx`
  - Set terms value as `data.url`
  - Tick mark the unique box.
- Create a database access key from the _Security_ tab in the left navigation, make sure the role is set to Admin, this access key needs to be set later on `FAUNADB_SECRET_KEY` env variable.

###### Lambda.store (Redis)

- Go to [Lambda.store](https://lambda.store/) web console
- Create a new database, name it as you like.
- Click on the database that you've just created, take a note for _Endpoint_, _Port_, and _Password_, these credentials need to be set later on `REDIS_HOST` and `REDIS_ PASSWORD` env variables.


### Environment Variables

Set these vars on any serverless provider you choose,

- `FAUNADB_SECRET_KEY` 
- `REDIS_HOST` (`[ENDPOINT]:[PORT]` i.e. `us1-xxxxx-xxxx-32223.lambda.store:32223`)
- `REDIS_PASSWORD` 
- `URL_HOST` (i.e. khal.web.id, _used to validate url input_)

Netlify only, edit `GO_IMPORT_PATH` value on `netlify.toml` file corresponds to your Github account

### Localhost

Build locally, don't forget to set your env vars,
```
$ go build -i
$ ./clapi

$ curl -X POST -d '2' http://0.0.0.0:3000/?url=http://clapi/clap  
true
$ curl -X GET http://0.0.0.0:3000/?url=http://clapi/clap
7
```

### Deploying on Vercel/Zeit (Git Integration)

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/git/external?repository-url=https%3A%2F%2Fgithub.com%2Fkh411d%2Fclapi)

Go files in the `/api` directory that export a function matching the net/http Go API will be served as Serverless Functions.

```
//ServeHTTP net/http handler
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	...
}
```

After deployed, as example you may access your function as this,  
__Base API URL__: `[YOUR-VERCEL-DOMAIN]/api/handler`

```
$ curl -X POST -d '2' https://clapi.vercel.app/api/handler?url=http://clapi/clap
true

$ curl -X GET https://clapi.vercel.app/api/handler?url=http://clapi/clap
7
```

### Deploying on Netlify (Git Integration)

[![Deploy to Netlify](https://www.netlify.com/img/deploy/button.svg)](https://app.netlify.com/start/deploy?repository=https://github.com/kh411d/clapi)

Netlify can build your source Go functions into AWS Lambda compatible binaries.

Before building your Go source, Netlify needs to know the expected Go import path for your project. Use the `GO_IMPORT_PATH` environment variable to set the right import path. You can do this in your _netlify.toml_ file. The path value should point to your source repository on your Git provider, for example github.com/kh411d/clapi.

After successful deployment, you need to add env variables required for Redis in  
`Settings > Build & deploy > Environment > Environment variables`  
and then re-deploy the app from,  
`Deploys > Trigger deploy > Deploy site`

After deployed, as example you may access your function as this,  
__Base API URL__: `[YOUR-NETLIFY-DOMAIN]/.netlify/functions/clapi`

```
$ curl -X POST -d '2' https://flamboyant-khorana-5b394b.netlify.app/.netlify/functions/clapi?url=http://clapi/clap
$ true

$ curl -X GET https://flamboyant-khorana-5b394b.netlify.app/.netlify/functions/clapi?url=http://clapi/clap
$ 7
```

### Deploying on Fly.io (flyctl)

I haven't tried this myself but it seems very possible, you need to have this `flyctl` tool installed on your system, upon installing, you have to register for fly.io account, simply follow the guidelines [Installing flyctl](https://fly.io/docs/getting-started/installing-flyctl/), and then continue with [Build, Deploy, and Run a Go Application](https://fly.io/docs/getting-started/golang/)

`flyctl init` will create a `fly.toml` file, Select Go (Go Builtin) for the builder, and then get ready for deployment `flyctl deploy`, keep in mind that flyctl will create docker images on your local computer, it will download all required images such as golang image approx 1gb, so make sure you have good storage space and internet connection.

`flyctl open` will open a browser on the HTTP version of the site.


