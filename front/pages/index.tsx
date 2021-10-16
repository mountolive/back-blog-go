import React from 'react'
import Pic from '~/components/pic.tsx'
import fetchAll from '~/lib/entries.ts'

export default function Home() {
  const imgSize = 200;
  // const envs = Deno.env.toObject();
  const [postsByDate, isSyncing] = fetchAll('http://localhost:8003/posts-by-date');

  return (
    <div className="page">
      <head>
        <title>Leo Guercio's (mountolive) site</title>
        <link rel="stylesheet" href="../style/index.css" />
      </head>
      <div className="logo"><Pic className="main-img" size={imgSize}/></div>
      <h1>Leo Guercio's (mountolive), or the blog that never was</h1>
      <p className="links">
        <a href="https://www.linkedin.com/in/leonardo-guercio-a9b31b35/" target="_blank">LinkedIn</a>
        <span></span>
        <a href="https://github.com/mountolive" target="_blank">Github</a>
      </p>
      <div className="entries">
        {isSyncing && (
          <p className="entries-txt">...</p>
        )}
        {!isSyncing && postsByDate.length > 0 && (
          <p className="entries-txt">Found something for you</p>
        )}
        {!isSyncing && postsByDate.length === 0 && (
         <p className="entries-txt">Nothing to see here</p>
        )}
      </div>
      <p className="copyinfo">All (potentially stupid) stuff in here were said by Leonardo Guercio</p>
    </div>
  )
}
