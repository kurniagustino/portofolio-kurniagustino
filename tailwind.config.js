/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./public/index.html"],
  darkMode: "class",
  theme: {
    container: {
      center: true,
      padding: "16px",
    },
    extend: {
      colors: {
        primary: "#22c55e",
        secondary: "#60a5fa",
        dark: "#0f172a",
      },
      screens: {
        "2xl": "1320px",
      },
    },
  },
  plugins: [],
};
