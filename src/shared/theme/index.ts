import { createTheme } from '@mantine/core'

export const theme = createTheme({
    scale: 1,
    fontSmoothing: true,
    focusRing: 'auto',
    white: '#ffffff',
    black: '#24292f',
    colors: {
        dark: [
            '#c9d1d9',
            '#b1bac4',
            '#8b949e',
            '#6e7681',
            '#484f58',
            '#30363d',
            '#21262d',
            '#161b22',
            '#0d1117',
            '#010409'
        ],
        gray: [
            '#f6f8fa',
            '#eaeef2',
            '#d0d7de',
            '#afb8c1',
            '#8c959f',
            '#6e7781',
            '#57606a',
            '#424a53',
            '#32383f',
            '#24292f'
        ],
        blue: [
            '#ddf4ff',
            '#b6e3ff',
            '#80ccff',
            '#54aeff',
            '#218bff',
            '#0969da',
            '#0550ae',
            '#033d8b',
            '#0a3069',
            '#002155'
        ],
        green: [
            '#dafbe1',
            '#aceebb',
            '#6fdd8b',
            '#4ac26b',
            '#2da44e',
            '#1a7f37',
            '#116329',
            '#044f1e',
            '#003d16',
            '#002d11'
        ],
        yellow: [
            '#fff8c5',
            '#fae17d',
            '#eac54f',
            '#d4a72c',
            '#bf8700',
            '#9a6700',
            '#7d4e00',
            '#633c01',
            '#4d2d00',
            '#3b2300'
        ],
        orange: [
            '#fff1e5',
            '#ffd8b5',
            '#ffb77c',
            '#fb8f44',
            '#e16f24',
            '#bc4c00',
            '#953800',
            '#762c00',
            '#5c2200',
            '#471700'
        ]
    },
    primaryShade: {
        light: 9,
        dark: 9
    },
    primaryColor: 'teal',
    autoContrast: true,
    luminanceThreshold: 0.3,
    fontFamily:
        "'-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Noto Sans', 'Helvetica', 'Arial', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji'",
    fontFamilyMonospace: "'SFMono-Regular', 'Consolas', 'Liberation Mono', 'Menlo', monospace",
    headings: {
        fontFamily:
            "'-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Noto Sans', 'Helvetica', 'Arial', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji'",
        fontWeight: '600'
    },
    defaultRadius: 'md'
})
