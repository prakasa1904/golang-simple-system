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

// qr code
var qrcode = new QRCode(document.getElementById("qrcode"), {
  text: "2@VlomCzv/gCK36hNHq0cqEG2UAmTI8Az55EzP6pLlcVwtLefGHrAqjnJSURYzBvKv/GXlZXgk1swP4g==,C3iPkKk599olKykuR8zNhj0cbp2VwxVzc3tuNLupwCI=,KN0FHG9BIIEGAwwsooCHgJJrMy1mXRIn+o2F3gG5bFk=,dAC7C5RdO9p35CtL/fhTcNykdj3wz+4SFvK4If1hH40=",
  width: 300,
  height: 300,
  colorDark : "#000000",
  colorLight : "#ffffff",
  correctLevel : QRCode.CorrectLevel.M
});

