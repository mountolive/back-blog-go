import React from 'react'
import fetchOne from '~/lib/fetchOne.ts'

export default function PostDetail(className: string, id: string) {
  // const envs = Deno.env.toObject();
  const [post, isSyncing] = fetchOne(`http://localhost:8003/posts`, id);

  return (
    <div className={className}>
      {isSyncing && (
        <p>...</p>
      )}
      {!isSyncing && post && (
        <div>
          <h1>{post.title}</h1>
          <div>{post.content}</div>
          <p>{post.createdAt}</p>
        </div>
      )}
      {!isSyncing && !post && (
         <p>Not found :/</p>
      )}
    </div>
  )
}
