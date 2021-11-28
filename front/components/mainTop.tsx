import React from 'react'
import Pic from '~/components/pic.tsx'
import Header from '~/components/header.tsx'

interface ComponentClasses {
  logoCls: string,
  linksCls: string,
}

export default function MainTop({ logoCls, linksCls }: ComponentClasses) {
  return (
    <>
      <Header/>
      <div className={logoCls}><Pic /></div>
      <h1>Leo Guercio's (mountolive), or the blog that never was</h1>
      <h2>Software Engineer, it seems</h2>
      <p className={linksCls}>
        <a href="https://www.linkedin.com/in/leonardo-guercio-a9b31b35/" target="_blank">LinkedIn</a>
        <span></span>
        <a href="https://github.com/mountolive" target="_blank">Github</a>
      </p>
    </>
  )
}
