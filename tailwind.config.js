/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.js", "./views/**/*.templ", "./views/**/*.go", "./node_modules/flowbite/**/*.js"],
  theme: {
    fontFamily: {
      sans: ["Space Mono", "system-ui"]
    },
    extend: {},
  },
  plugins: [
    require("flowbite/plugin")
  ],
}
