/** @type {import('tailwindcss').Config} */
export const content = [
  './*.html', // Include your HTML files here
  './*.go', // Include your Go files if you have HTML in them
  './**/*.html',
  './**/*.go',
];
export const theme = {
  extend: {},
};
export const plugins = [
  require('daisyui'),
];

export const daisyui = {
  themes: ["light", "dark"]
}
