const functions = require('firebase-functions');
const request = require('request');
const querystring = require('querystring');

exports.search = functions.https.onRequest((req, res) => {
	request(
		'https://api.tenor.com/v1/search?'
		+ querystring.stringify(
			{
				key: functions.config().tenor.key,
				safesearch: 'mild',
				limit: 25,
				media_filter: 'minimal',
				q: req.query.q
			}
		),
		function (error, tenorRes, body) {
			if (error) {
				res.end();
				return;
			}
			if (tenorRes.statusCode != 200) {
				res.end();
				return;
			}
			var parsedRes = JSON.parse(body);
			console.log(parsedRes);
			var results = [];
		  parsedRes.results.forEach(function(r, i) {
			  var media = r.media[0];
				results.push({
					images: {
						original: {
						  url: media.gif.url,
							mp4: media.mp4.url,
							width: media.mp4.dims[0],
							height: media.mp4.dims[1]
						}
					}
		    })
			})
			res.status(200);
			res.setHeader('content-type', 'application/json');
			res.set('Access-Control-Allow-Origin', "*")
			res.set('Access-Control-Allow-Methods', 'GET')
			res.send({ source: 'tenor', data: results });
		}
	);
});
