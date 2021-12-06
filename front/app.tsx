import React, { FC } from 'react'

export default function App({ Page, pageProps }: { Page: FC, pageProps: Record<string, unknown> }) {
  return (
    <main>
      <head>
        <meta name="viewport" content="width=device-width" />
        <meta property="og:title" content="Leo Guercio's blog" />
        <meta property="og:image" content="https://www.aparences.net/wp-content/uploads/2017/11/de-chirico-portrait-apollinaire.jpg" />
      </head>
      <Page {...pageProps} />
    </main>
  )
}
