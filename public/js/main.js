// legacy plugins
htmx.defineExtension("push-url-with-query", {
  onEvent: function (name, e) {
    if (name === "htmx:confirm") {
      // get value from current field (such as: search field)
      const url = e.target.getAttribute("hx-push-url-with-query");
      const name = e.target.getAttribute("name");
      const value = e.target.value;

      const targetURL = new URL(url, window.location.origin);
      const checkIfExist = targetURL.searchParams.get(name);

      if (checkIfExist === "") {
        // make sure if contains empty ?name=
        targetURL.searchParams.delete(name);

        targetURL.searchParams.append(name, value);
      } else {
        targetURL.searchParams.set(name, value);
      }

      // push to browser history
      window.history.pushState({}, "", targetURL);
    }
  },
});

// refetch URL after execute xhr
// this extension useful if you have other URL dependency and want to mount in another document
// References:
// https://htmx.org/api/
// https://htmx.org/events/
htmx.defineExtension("refetch-url", {
  onEvent: function (name, e) {
    if (name === "htmx:xhr:loadend") {
      const targetPath = e.target.getAttribute("hx-refetch-url");
      const fetchMethod = e.target.getAttribute("hx-refetch-method") || "GET";
      const outputTarget = e.target.getAttribute("hx-refetch-target") || "body";
      // with query, will append current URL query to target url
      const withQuery =
        e.target.getAttribute("hx-refetch-with-query") || "false";

      const currentURL = new URL(window.location.href);
      const targetURL = new URL(targetPath, window.location.origin);

      if (withQuery === "true") {
        if (currentURL.searchParams.size) {
          currentURL.searchParams.forEach((val, key) => {
            targetURL.searchParams.append(key, val);
          });
        }
      }

      // refetch other URL to update UI
      htmx.ajax(fetchMethod, targetURL.href, outputTarget);
    }
  },
});

// replace current URL with extra query from form data
htmx.defineExtension("push-current-url", {
  onEvent: function (name, e) {
    // var formData = new FormData(e.target)
    if (name === "htmx:configRequest") {
      const currentURL = new URL(window.location.href);
      const formData = new FormData(e.target);

      formData.forEach((val, key) => {
        const checkIfExist = currentURL.searchParams.get(key);

        if (checkIfExist === "") {
          currentURL.searchParams.delete(key);
          currentURL.searchParams.append(key, val);
        } else {
          currentURL.searchParams.set(key, val);
        }
      });

      // push to browser history
      window.history.pushState({}, "", currentURL);
    }
  },
});

// debug event in HTMX
htmx.defineExtension("debug", {
  onEvent: function (name, e) {
    console.log("name: ", name);
    console.log("event: ", e);
  },
});

// flowbite init
htmx.defineExtension("flowbite-drawer", {
  onEvent: function (name, e) {
    if (name === "htmx:trigger") {
      const drawer = e.target.getAttribute("hx-flowbite-drawer");
      
      FlowbiteInstances._instances.Drawer[drawer].show();
    }
  },
});

// qr code
const qrTarget = document.getElementById("qrsrt");
if (qrTarget) {
  const qrContent = qrTarget.textContent;
  var qrcode = new QRCode(document.getElementById("qrcode"), {
    text: qrContent,
    width: 400,
    height: 400,
    correctLevel: QRCode.CorrectLevel.L,
  });
}
