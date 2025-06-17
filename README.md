# kakakikikeke's Random Podcast Player
[Here](https://random.kakakikikeke.com/)

# Build on local
* go mod tidy
* go fmt ./...
* go test ./...
* go build
* ./random-podcast
* curl localhost:8080

# Deploy to GAE
* gcloud config configurations activate default
* gcloud app deploy
* gcloud app logs tail -s default

or

* gcloud beta app repaire && gcloud app deploy --no-cache
