export CLOUDFLARE_ZONE = f24e7be338605506f4cd84f8ff6dda42

deploy: clean build push cdn

clean:
	rm -fR build

run:
	bundle exec middleman server

build:
	bundle exec middleman build

push:
	git branch -D gh-pages 2>/dev/null | true
	git branch -D gh-pages-draft 2>/dev/null | true
	cd .. && \
		git checkout -b gh-pages-draft && \
		git add -f website/build && \
		git commit -m "Deploy to gh-pages" && \
		git subtree split --prefix website/build -b gh-pages && \
		git push --force origin gh-pages:gh-pages && \
		git checkout - && \
		cd -

cdn:
	curl -X DELETE "https://api.cloudflare.com/client/v4/zones/$(CLOUDFLARE_ZONE)/purge_cache" \
		-H "X-Auth-Email: $(CLOUDFLARE_EMAIL)" \
		-H "X-Auth-Key: $(CLOUDFLARE_CLIENT_API_KEY)" \
		-H "Content-Type: application/json" \
		--data '{"purge_everything":true}'

