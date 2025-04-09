/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{html,ts}"
  ],
  theme: {
    extend: {
      colors: {
        primary: '#5C00BF',
        secondary: '#FF8600',
        lightblue: '#AEB8FE',
        dark: '#2E2E2E',
        accent: 'rgba(255, 134, 0, 0.38)',
        white: '#FFFFFF',
        black: '#000000',
      },
      fontFamily: {
        ocr: ['"OCR A Std Regular"', 'monospace'],
        press: ['"Press Start 2P"', 'cursive'],
      },
      fontSize: {
        'title-lg': '80px',
        'title-md': '30px',
        'text-lg': '20px',
        'text-sm': '14px',
      },
    },
  },
  plugins: [],
}
