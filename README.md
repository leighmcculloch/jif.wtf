# jif.wtf

[jif.wtf](http://jif.wtf) is an alternative to searching for gifs on Giphy. Gives you:

* Fast loading of results as MP4s.
* Quick tabbing through the results.
* Barebones interface.
* Quick click & copy paste of direct gif URL.

![](README-screenshot.png)

# Why

Because I needed a tool that could take me from a reaction in my mind to a killer gif in under 3 seconds. The ability to tab through fully loaded gifs quickly was key to this and a feature of Insert Pic which I loved. I just needed it in the form of a website for use on computers where installing arbritrary software is not ideal.

Yes, I know. There are already many wesites that let me search for gifs. They are amazing, this does nothing special. It just has the UI I wanted.

# Deployment

## Static Website

Built using `middleman` and deploys to AWS S3.

```bash
cd website/
./deploy
```

## Websocket Service

Built in `go` and deploys to [Pivotal Website Services](https://run.pivotal.io).

```bash
cd service/
cf login
cf push
```

# Inspiration
Inspired by [Insert Pic](http://www.getinsertpic.com/).

# Powered by
![](README-powered-by-giphy.gif)
