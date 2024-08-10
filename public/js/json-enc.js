htmx.defineExtension("json-enc", {
  onEvent: function (name, evt) {
    if (name === "htmx:configRequest") {
      evt.detail.headers["Content-Type"] = "application/json";
    }
  },

  encodeParameters: function (xhr, parameters, elt) {
    xhr.overrideMimeType("text/json");
    
    // convert to correct data type
    const validParameters = {};
    const regexNumbOnly = new RegExp('^[0-9]$');
    const keys = Object.keys(parameters);

    keys.forEach((key) => {
      const value = parameters.get(key);

      if (regexNumbOnly.test(value)) {
        validParameters[key] = Number(value);
      } else {
        validParameters[key] = value;
      }

    });

    console.log(validParameters);

    return JSON.stringify(validParameters);
  },
});
