/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./templates/*.html"],
    theme: {
        fontFamily: { sans: ["Roboto", "sans-serif"] },
    },
    plugins: [
        require('@tailwindcss/forms'),
    ],
};
