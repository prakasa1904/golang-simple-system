const plugin = require("tailwindcss/plugin");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/**/*.{html,js}", "./files/assets/**/*.{css}", "./node_modules/flowbite/**/*.js"],
  theme: {
    extend: {},
  },
  plugins: [
    require("flowbite/plugin"),
  ],
};
