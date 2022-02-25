module.exports = {
    content: ["./src/**/*.{html,js,vue}"],
    theme: {
        extend: {
            fontFamily: {
                'sans': ['Inter', 'Helvetica', 'Arial', 'sans-serif'],
            }
        }
    },
    variants: {
        backgroundColor: ['responsive', 'odd', 'hover', 'focus'],
    },
    plugins: [require('@tailwindcss/forms')]
}