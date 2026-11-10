/**
 * Konfigurasi Tailwind CSS untuk build produksi.
 *
 * Bangun ulang CSS setelah mengubah template HTML:
 *   npx tailwindcss@3.4.17 -i static/src/input.css -o static/css/tailwind.css --minify
 *
 * (CSS hasil build ikut di-commit supaya deploy tidak butuh Node.js)
 */
module.exports = {
  darkMode: 'class',
  content: ['./templates/*.html'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
        mono: ['Fira Code', 'monospace'],
      },
      colors: {
        accent: {
          cyan: '#06b6d4',
          blue: '#3b82f6',
          indigo: '#6366f1',
          purple: '#a855f7',
        },
        slate: {
          850: '#151e2e',
          900: '#0f172a',
          950: '#020617',
        },
      },
      animation: {
        'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
    },
  },
  plugins: [],
};
