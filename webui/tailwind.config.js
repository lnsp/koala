module.exports = {
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
    purge: [
        './src/**/*.vue',
    ],
}