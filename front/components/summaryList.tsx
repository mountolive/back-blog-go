import React from 'react';
import PostSummary from '~/lib/fetchAll.ts';

export default function SummaryList(className: string, posts: PostSummary[]) {
  return (
    <div className={className}>
      {
        posts.map((post: PostSummary) => {
          return Summary(post);
        })
      }
    </div>
  );
}

function Summary(post: PostSummary) {
  return (
    <div>
      <h1 onClick={_event => goToPost(post.id, post.title)}>{post.title}</h1>
      <p>{post.createdAt}</p>
    </div>
  )
}

function goToPost(id: string, title: string) {
  location.href = `/post?id=${id}&title=${urlFriendly(title)}`
}

function urlFriendly(txt: string): string {
  return txt.replace(/[^a-z0-9_]+/gi, '-').replace(/^-|-$/g, '').toLowerCase()
}
