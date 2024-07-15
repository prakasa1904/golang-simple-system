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
