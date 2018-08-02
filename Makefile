deploy: clean build push

clean:
	cd website && rm -fR build

run:
	cd website && bundle exec middleman server

build:
	cd website && bundle exec middleman build
	cd service && go build

push: push-service push-functions push-hosting

push-functions:
	firebase deploy --only functions

push-service:
	cd service && gcloud app deploy --project jif-wtf-bf47b

push-hosting:
	firebase deploy --only hosting
