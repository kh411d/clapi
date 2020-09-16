# clapi

Self-hosted serverless/faas blogging clap api 

### Deploying on Vercel/Zeit 

Within the /api directory of your projects, Vercel will automatically recognize the languages listed on this page, through their file extensions, and serve them as Serverless Function. Go files in the api directory that export a function matching the net/http Go API will be served as Serverless Functions.

```
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	var res map[string]interface{}

	switch r.Method {
	case "GET":
		res = getClaps()
	case "POST":
		res = addClap()
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
```

It seems that when you put more than one go file in the `api` folder the system will examine all the file to look for net/HTTP exported functions, if there is no exported function in one of that files, it will resulting error thus the deployment will fail, so it is why I put all the code in on file handler.go

After deployed, you'll have some domains to be used on function,
```
$ curl -X POST https://clapi.vercel.app/api/handler
{"Body":"Adding clap","StatusCode":200}
$ curl -X GET https://clapi.vercel.app/api/handler
{"Body":"Get clap","StatusCode":200}
```