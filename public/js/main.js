// temp solution from https://github.com/bigskysoftware/htmx/issues/2156
htmx.defineExtension("push-extra-query", {
  onEvent: function (name, e) {
    console.log(name);
    if (name === "htmx:confirm") {
      // get value from current field (such as: search field)
      const name = e.target.getAttribute("name");
      const value = e.target.value;

      const currentURL = new URL(window.location);
      const checkIfExist = currentURL.searchParams.get(name);

      if (checkIfExist === "") {
        currentURL.searchParams.append(name, value);
      } else {
        currentURL.searchParams.set(name, value);
      }

      // push to browser history
      window.history.pushState({}, "", currentURL);
    }
  },
});
