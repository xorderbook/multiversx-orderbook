/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        panel: '#27272B',
        main: '#202020',
      },
      fontFamily: {
        sans: ['"Inter"', 'sans-serif']

      }
    },
    backgroundImage: {
      // eslint-disable-next-line quotes
      'mvx-white': "url('../multiversx-white.svg')"
    }
  },
  plugins: []
};
