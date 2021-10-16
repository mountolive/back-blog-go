import React from 'react';
import PostSummary from '~/lib/PostSummary';

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
			<h1>{post.title}</h1>
			<p>{post.createdAt}</p>
		</div>
	)
}
