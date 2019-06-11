// Adapted from github.com/jcuga/golongpoll
// MIT

(function(jq) {

  'use strict';

  jq.longpoll = function(url, category, eventFnc, options) {

    if (typeof window.console == 'undefined') {
      window.console = {
        log: function(
          msg) {}
      };
    }

    var settings = jq.extend({
      timeout: 45, // seconds
      delaySuccess: 10, // milliseconds
      delayError: 3000 // milliseconds
    }, options);

    // var timeout = 45; // in seconds
    var baseUrl = url + '?timeout=' + encodeURIComponent(settings.timeout) +
      '&category=' + encodeURIComponent(category) + '&since_time=';
    // Start checking for any events that occurred after page load time (right now)
    // Notice how we use .getTime() to have num milliseconds since epoch in UTC
    // This is the time format the longpoll server uses.
    var sinceTime = (new Date(Date.now())).getTime();
    // var delaySuccess = 10; // 10 ms
    // var delayError = 3000; // 3 sec
    (function poll() {
      jq.ajax({
        url: baseUrl + sinceTime,
        success: function(data) {
          if (data && data.events && data.events.length > 0) {
            // NOTE: these events are in chronological order (oldest first)
            for (var i = 0; i < data.events.length; i++) {
              var event = data.events[i];
              eventFnc(event);
              sinceTime = event.timestamp;
            }
            // success!  start next longpoll
            setTimeout(poll, settings.delaySuccess);
            return;
          }
          if (data && data.timeout) {
            console.log(
              "No events, checking again.");
            // no events within timeout window, start another longpoll:
            setTimeout(poll, settings.delaySuccess);
            return;
          }
          if (data && data.error) {
            console.log("Error response: " + data.error);
            console.log("Trying again shortly...");
            setTimeout(poll, settings.delayError);
            return;
          }
          // We should have gotten one of the above 3 cases:
          // either nonempty event data, a timeout, or an error.
          console.log(
            "Didn't get expected event data, try again shortly..."
          );
          setTimeout(poll, settings.delayError);
        },
        dataType: "json",
        error: function(data) {
          console.log(
            "Error in ajax request--trying again shortly..."
          );
          setTimeout(poll, settings.delayError);
        }
      });
    })();

    return jq;

  };

}(jQuery));
