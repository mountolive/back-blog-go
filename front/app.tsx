import React, { FC } from 'react'

export default function App({ Page, pageProps }: { Page: FC, pageProps: Record<string, unknown> }) {
  return (
    <main>
      <head>
        <meta name="viewport" content="width=device-width" />
        <meta name="title" property="og:title" content="Leo Guercio's blog" />
        <meta name="type" property="og:type" content="website" />
        <meta name="image" property="og:image" content="https://www.aparences.net/wp-content/uploads/2017/11/de-chirico-portrait-apollinaire.jpg" />
        <meta name="url" property="og:url" content="https://www.leoponc.io" />
      </head>
      <Page {...pageProps} />
    </main>
  )
}
