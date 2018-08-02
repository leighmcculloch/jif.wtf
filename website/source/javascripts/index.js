(function () {
  var searchForm;
  var searchField;
  var resultContainer;
  var resultDisplay;
  var resultMessage;
  var resultUrl;
  var resultUrlDownload;
  var poweredByGiphy;
  var poweredByTenor;
  var lastSearchQuery;
  var lastSearchTimeoutId;

  var source = new URLSearchParams(window.location.search).get('source') || "tenor";
  var results;
  var resultIndex;

  function init() {
    searchForm = document.getElementById('search-form');
    searchField = document.getElementById('search-field');
    resultContainer = document.getElementById('result-container');
    resultDisplay = document.getElementById('result-display');
    resultMessage = document.getElementById('result-message');
    resultUrl = document.getElementById('result-url');
    resultUrlDownload = document.getElementById('result-url-download');
    navigation = document.getElementById('navigation-arrows');
    navigationLeft = document.getElementById('navigation-arrow-left');
    navigationRight = document.getElementById('navigation-arrow-right');
    poweredByGiphy = document.getElementById('powered-by-giphy');
    poweredByTenor = document.getElementById('powered-by-tenor');

    navigationLeft.addEventListener('click', function() { navigate(-1); });
    navigationRight.addEventListener('click', function() { navigate(1); });

    searchForm.addEventListener('submit', formSubmit);
    searchField.addEventListener('keyup', function() { searchWithDelay(); });
    searchField.addEventListener('change', function() { search(); });

    document.addEventListener('keydown', searchKeyDown);

    resultUrl.addEventListener('click', resultUrlSelect);

    searchField.focus();

    processPushState(window.location.hash);
  }

  function processPushState(pushState) {
    var goto = decodeURIComponent(pushState.replace(/^#/, ''));
    var gotoParts = goto.split('⇢');
    if (gotoParts.length >= 1) {
      searchField.value = gotoParts[0];
      var optionalIndex = undefined;
      if (gotoParts.length >= 2) {
        optionalIndex = gotoParts[1] - 1;
      }
      search(optionalIndex);
    }
  }

  function updatePushState(query, optionalIndex) {
    var state = '';
    if (lastSearchQuery !== undefined && lastSearchQuery.length > 0) {
      state += '#' + lastSearchQuery;
      if (optionalIndex !== undefined && optionalIndex > 0) {
        state += '⇢' + (optionalIndex + 1);
      }
    }
    window.history.pushState(null, null, state);
  }

  function formSubmit(e) {
    e.preventDefault();
    search();
    return false;
  }

  function searchKeyDown(e) {
    var TAB_KEY = 9;
    var ARROW_LEFT_KEY = 37;
    var ARROW_RIGHT_KEY = 39;

    switch (e.keyCode) {
      case TAB_KEY:
        navigate(e.shiftKey ? -1 : 1);
        e.preventDefault();
        return false;
      case ARROW_LEFT_KEY:
        navigate(-1);
        e.preventDefault();
        return false;
      case ARROW_RIGHT_KEY:
        navigate(1);
        e.preventDefault();
        return false;
    }
  }

  function navigate(vector) {
    resultIndex += vector;
    if (typeof results !== 'undefined') {
      if (resultIndex < 0) {
        resultIndex = results.length - 1;
      } else if (resultIndex >= results.length) {
        resultIndex = 0;
      }
    }
    displayResult(resultIndex);
  }

  function searchWithDelay(optionalIndex) {
    clearTimeout(lastSearchTimeoutId);
    lastSearchTimeoutId = setTimeout(function() { search(optionalIndex) }, 500);
  }

  function search(optionalIndex) {
    clearTimeout(lastSearchTimeoutId);
    if (lastSearchQuery == searchField.value) {
      return;
    }
    lastSearchQuery = searchField.value;
    updatePushState(lastSearchQuery, optionalIndex)

    var baseUrl;
    if (source == "tenor") {
      baseUrl = 'https://api.tenor.com/v1/search?q=';
    } else if (source == "tenor-function") {
      baseUrl = 'https://us-central1-jif-wtf-bf47b.cloudfunctions.net/search?q=';
    } else if (source == "tenor-appengine") {
      baseUrl = 'https://jif-wtf-bf47b.appspot.com/search?q=';
    } else if (source == "giphy") {
      baseUrl = 'https://api.giphy.com/v1/gifs/search?api_key=dc6zaTOxFJmzC&q=';
    } else {
      throw "Source unknown: " + source;
    }

    makeCORSRequest(
      baseUrl + encodeURIComponent(lastSearchQuery),
      function(response) {
        var data;
        if (source == "tenor") {
          data = response.results.map(function(r) {
            return {
              images: {
                original: {
                  url: r.media[0].gif.url,
                  mp4: r.media[0].mp4.url,
                  width: r.media[0].mp4.dims[0],
                  height: r.media[0].mp4.dims[1]
                }
              }
            };
          })
        } else {
          data = response.data;
        }
        onResults(lastSearchQuery, data, optionalIndex);
      }
    );
  }

  function onResults(query, newResults, optionalIndex) {
    if (query != lastSearchQuery) {
      return;
    }
    results = newResults;
    window.results = results;
    resultIndex = optionalIndex || 0;
    displayResult(resultIndex);
    if (query.length > 0) {
      if (source == "tenor" || source == "tenor-function" || source == "tenor-appengine") {
        poweredByTenor.style.display = 'block';
      } else if (source == "giphy") {
        poweredByGiphy.style.display = 'block';
      }
    }
  }

  function displayResult(index) {
    if (typeof results === 'undefined' || results == null || results.length == 0) {
      resultDisplay.src = '';
      if (searchField.value.length > 0) {
        resultMessage.innerHTML = 'No results. Search for something else.';
      } else {
        resultMessage.innerHTML = '';
      }
      resultUrl.innerHTML = '';
      resultUrlDownload.innerHTML = '';
      navigation.style.display = 'none';
      return;
    }

    navigation.style.display = 'block';

    if (index < 0 || index >= results.length) {
      index = 0;
    }

    updatePushState(lastSearchQuery, index);

    result = results[index].images.original;

    resultDisplay.type = 'video/mp4';
    resultDisplay.src = result.mp4;

    var height = Math.min(result.height, 200);
    var width = result.width * height/result.height;
    resultDisplay.width = width;
    resultDisplay.height = height;

    resultMessage.innerHTML = "Result " + (index + 1) + " of " + results.length;

    if (source == "tenor" || source == "tenor-function" || source == "tenor-appengine") {
      ids = getTenorIds(result.url);
    } else if (source == "giphy") {
      ids = getGiphyIds(result.url);
    } else {
      throw "Source unknown: " + source;
    }

    var url = getIntermediateGifUrl(ids.sourceId, ids.mediaId, ids.imageId);
    resultUrl.innerHTML = url;
    resultUrlDownload.innerHTML = '<a href="' + url + '" download="' + ids.imageId + '.gif">download</a>';
  }

  function getGiphyIds(giphyUrl) {
    var re = /http(s?):\/\/media(\d+).giphy.com\/media\/([0-9a-zA-Z]+)\/[a-z0-9_]+.gif/g;
    var matches = re.exec(giphyUrl);
    if (matches.length != 4) {
      return giphyUrl;
    }
    return { sourceId: "0", mediaId: matches[2], imageId: matches[3] }
  }

  function getTenorIds(tenorUrl) {
    var re = /http(s?):\/\/media.tenor.com\/images\/([0-9a-z]+)\/tenor.gif/g;
    var matches = re.exec(tenorUrl);
    if (matches.length != 3) {
      return tenorUrl;
    }
    return { sourceId: "1", mediaId: "0", imageId: matches[2] }
  }

  function getIntermediateGifUrl(sourceId, mediaId, imageId) {
    return "https://jif.wtf/" + sourceId + mediaId + "." + imageId + ".gif";
  }

  function resultUrlSelect(e) {
    selectTextOfElement(resultUrl);
  }

  function selectTextOfElement(element) {
    if (document.selection) {
      var range = document.body.createTextRange();
      range.moveToElementText(element);
      range.select();
    } else if (window.getSelection) {
      var range = document.createRange();
      range.selectNode(element);
      window.getSelection().addRange(range);
    }
  }

  function makeCORSRequest(url, complete) {
    var method = 'GET';

    var xhr = new XMLHttpRequest();
    if ('withCredentials' in xhr) {
      xhr.open(method, url, true);
    } else if (typeof XDomainRequest != 'undefined') {
      xhr = new XDomainRequest();
      xhr.open(method, url);
    } else {
      console.log('CORS is not supported.');
      return;
    }

    xhr.onload = function() {
      complete(JSON.parse(xhr.responseText));
    }
    xhr.onerror = function() {
      console.log('An XHR error occurred.');
    };

    xhr.send();

    return xhr;
  }

  window.addEventListener("load", init);
})();
