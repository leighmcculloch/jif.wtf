(function () {
  var SERVICE_ADDRESS = 'https://jifwtfservice.cfapps.io';

  var searchForm;
  var searchField;
  var resultContainer;
  var resultDisplay;
  var resultMessage;
  var resultUrl;
  var lastSearchQuery;
  var lastSearchTimeoutId;

  var results;
  var resultIndex;

  function init() {
    searchForm = document.getElementById('search-form');
    searchField = document.getElementById('search-field');
    resultContainer = document.getElementById('result-container');
    resultDisplay = document.getElementById('result-display');
    resultMessage = document.getElementById('result-message');
    resultUrl = document.getElementById('result-url');
    navigation = document.getElementById('navigation-arrows');
    navigationLeft = document.getElementById('navigation-arrow-left');
    navigationRight = document.getElementById('navigation-arrow-right');

    navigationLeft.addEventListener('click', function() { navigate(-1); });
    navigationRight.addEventListener('click', function() { navigate(1); });

    searchForm.addEventListener('submit', formSubmit);
    searchField.addEventListener('keyup', searchWithDelay);
    searchField.addEventListener('change', search);

    document.addEventListener('keydown', searchKeyDown);

    resultUrl.addEventListener('click', resultUrlSelect);

    searchField.focus();

    searchField.value = window.location.hash.replace(/^#/, '');
    if (searchField.value.length > 0) {
      searchWithDelay();
    }
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

  function searchWithDelay() {
    clearTimeout(lastSearchTimeoutId);
    lastSearchTimeoutId = setTimeout(search, 500);
  }

  function search() {
    clearTimeout(lastSearchTimeoutId);
    if (lastSearchQuery == searchField.value) {
      return;
    }
    lastSearchQuery = searchField.value;
    window.history.pushState(null, null, '#' + lastSearchQuery);
    makeCORSRequest(
      SERVICE_ADDRESS + '/search?query=' + encodeURIComponent(lastSearchQuery),
      function(response) {
        onResults(response.query, response.search_results);
      }
    );
  }

  function onResults(query, newResults) {
    if (query != lastSearchQuery) {
      return;
    }
    results = newResults;
    window.results = results;
    resultIndex = 0;
    displayResult(resultIndex);
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
      navigation.style.display = 'none';
      return;
    }

    navigation.style.display = 'block';

    var url = results[index].gif;

    resultDisplay.type = 'video/mp4';
    resultDisplay.src = results[index].mp4;
    resultDisplay.width = results[index].width;
    resultDisplay.height = results[index].height;

    resultMessage.innerHTML = "Result " + (index + 1) + " of " + results.length;
    resultUrl.innerHTML = url;
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
