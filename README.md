# kakakikikeke's Random Podcast Player
[Here](https://random.kakakikikeke.com/)

# Build on local
* go mod tidy
* go fmt ./...
* go build
* python /opt/homebrew/share/google-cloud-sdk/bin/dev_appserver.py .

# Deploy to GAE
* gcloud config configurations activate default
* gcloud app deploy
* gcloud app logs tail -s default

or

* gcloud beta app repaire && gcloud app deploy --no-cache
