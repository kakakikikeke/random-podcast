# local
/Users/kakakikikeke/go/src/bitbucket.org/team-kaka/golang-random-podcast に clone すること  
google-cloud-sdk を homebrew でインストールしてあること  
.python-version に記載の python がプロジェクトで使えるようになっていること

* ln -s /Users/kakakikikeke/go/src/bitbucket.org/team-kaka/golang-random-podcast /Users/kakakikikeke/data/repo
cd /Users/kakakikikeke/go/src/bitbucket.org/team-kaka/golang-random-podcast
* go fmt
* go build
* python /opt/homebrew/share/google-cloud-sdk/bin/dev_appserver.py .

# deploy
* gcloud config configurations activate default
* gcloud app deploy
* gcloud app logs tail -s default

or

* gcloud beta app repaire && gcloud app deploy --no-cache
