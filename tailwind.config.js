/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./components/*.{go,js,templ,html}"
    ],
    theme: {
      extend: {
        animation: {
            // Reduced from default 2s to 0.75s
            'pulse-fast': 'pulse 0.75s cubic-bezier(0.6, 1, 0.6, 1) infinite',
            },
        },
    },
    plugins: [],
  }
