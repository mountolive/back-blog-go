import config from '~/lib/config.ts'
import fetchOne from '~/lib/fetchOne.ts'
import GoBack from '~/components/goBack.tsx'
import React from 'react'

export default function PostDetail(className: string, id: string) {
  const envs = config();
  const [post, isSyncing] = fetchOne(
    `https://api.leoponc.io/posts`,
    id,
  );

  return (
    <div className={className}>
      {isSyncing && (
        <p>...</p>
      )}
      {!isSyncing && post && (
        <div>
          <h1>{post.title}</h1>
          <p className="postdate">{post.createdAt}</p>
          <div dangerouslySetInnerHTML={{__html: post.content}}/>
        </div>
      )}
      {!isSyncing && !post && (
         <p>Not found :/</p>
      )}
    </div>
  )
}
